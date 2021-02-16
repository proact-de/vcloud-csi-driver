package controller

import (
	"context"
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/volume"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateVolume implements the CSI standard definition.
func (s *Service) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "missing name")
	}

	if len(req.VolumeCapabilities) == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing volume capabilities")
	}

	for i, cap := range req.VolumeCapabilities {
		if !isCapabilitySupported(cap) {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("capability %d is not supported", i))
		}
	}

	size, err := sizeFromCapacityRange(req.GetCapacityRange())

	if err != nil {
		return nil, status.Errorf(codes.OutOfRange, "invalid capacity range: %v", err)
	}

	profile := ""
	if value, ok := req.Parameters["storageProfile"]; ok {
		profile = value
	}

	resp := &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      req.Name,
			CapacityBytes: size,
			AccessibleTopology: []*csi.Topology{
				{
					Segments: map[string]string{
						"failure-domain.beta.kubernetes.io/zone": s.datacenter,
					},
				},
			},
		},
	}

	if err := s.volume.Create(ctx, volume.CreateOpts{
		Name:    req.Name,
		Size:    size,
		Profile: profile,
	}); err != nil {
		code := codes.Internal

		switch err {
		case volume.ErrVolumeAlreadyExists:
			code = codes.AlreadyExists
		}

		return nil, status.Error(code, fmt.Sprintf("failed to create volume: %s", err))
	}

	return resp, nil
}

// DeleteVolume implements the CSI standard definition.
func (s *Service) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing volume id")
	}

	if err := s.volume.Delete(ctx, volume.DeleteOpts{
		ID: req.VolumeId,
	}); err != nil {
		code := codes.Internal

		switch err {
		case volume.ErrVolumeNotFound:
			return &csi.DeleteVolumeResponse{}, nil
		case volume.ErrStillAttached:
			code = codes.FailedPrecondition
		}

		return nil, status.Error(code, fmt.Sprintf(
			"failed to delete volume: %s",
			err,
		))
	}

	return &csi.DeleteVolumeResponse{}, nil
}

// GetCapacity implements the CSI standard definition.
func (s *Service) GetCapacity(ctx context.Context, req *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
