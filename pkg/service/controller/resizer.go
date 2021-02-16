package controller

import (
	"context"
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/proact-de/vcloud-csi-driver/pkg/model"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/volume"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ControllerExpandVolume implements the CSI standard definition.
func (s *Service) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing volume id")
	}

	size, err := sizeFromCapacityRange(req.GetCapacityRange())

	if err != nil {
		return nil, status.Errorf(codes.OutOfRange, "invalid capacity range: %v", err)
	}

	if err := s.volume.Resize(ctx, volume.ResizeOpts{
		Volume: &model.Volume{
			ID: req.VolumeId,
		},
		Size: size,
	}); err != nil {
		code := codes.Internal

		switch err {
		case volume.ErrVolumeNotFound:
			code = codes.NotFound
		}

		return nil, status.Error(code, fmt.Sprintf(
			"failed to expand volume: %s",
			err,
		))
	}

	v, err := s.volume.Find(ctx, volume.FindOpts{
		ID: req.VolumeId,
	})

	if err != nil {
		code := codes.Internal

		switch err {
		case volume.ErrVolumeNotFound:
			code = codes.NotFound
		}

		return nil, status.Error(code, fmt.Sprintf(
			"failed to expand volume: %s",
			err,
		))
	}

	nodeExpansionRequired := true

	if req.GetVolumeCapability() != nil {
		switch req.GetVolumeCapability().GetAccessType().(type) {
		case *csi.VolumeCapability_Block:
			nodeExpansionRequired = false
		}
	}

	return &csi.ControllerExpandVolumeResponse{
		CapacityBytes:         v.Size,
		NodeExpansionRequired: nodeExpansionRequired,
	}, nil
}
