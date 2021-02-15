package node

import (
	"context"
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/resize"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/volume"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NodeExpandVolume implements the CSI standard definition.
func (s *Service) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing volume id")
	}

	if req.VolumePath == "" {
		return nil, status.Error(codes.InvalidArgument, "missing volume path")
	}

	v, err := s.volume.Find(ctx, volume.FindOpts{
		ID: req.VolumeId,
	})

	if err != nil {
		switch err {
		case volume.ErrVolumeNotFound:
			return nil, status.Error(codes.NotFound, "volume not found")
		default:
			return nil, status.Error(codes.Internal, fmt.Sprintf(
				"failed to get volume: %s",
				err,
			))
		}
	}

	exists, err := s.mount.Exists(ctx, v.Device)

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf(
			"failed to check for volume existence: %s",
			err,
		))
	}

	if !exists {
		return nil, status.Error(codes.Unavailable, fmt.Sprintf(
			"volume %s is not available on %v",
			v.Device,
			s.server,
		))
	}

	if err := s.resize.Resize(ctx, resize.Opts{
		Volume: v,
		Target: req.VolumePath,
	}); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf(
			"failed to resize volume: %s",
			err,
		))
	}

	return &csi.NodeExpandVolumeResponse{
		CapacityBytes: v.Size,
	}, nil
}
