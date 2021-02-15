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

// ControllerPublishVolume implements the CSI standard definition.
func (s *Service) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing volume id")
	}

	if req.NodeId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing node id")
	}

	if req.VolumeCapability == nil {
		return nil, status.Error(codes.InvalidArgument, "missing volume capabilities")
	}

	if !isCapabilitySupported(req.VolumeCapability) {
		return nil, status.Error(codes.InvalidArgument, "capability is not supported")
	}

	if req.Readonly {
		return nil, status.Error(codes.InvalidArgument, "readonly volumes are not supported")
	}

	if err := s.volume.Attach(ctx, volume.AttachOpts{
		Volume: &model.Volume{
			Name: req.VolumeId,
		},
		Server: &model.Server{
			Name: req.NodeId,
		},
	}); err != nil {
		code := codes.Internal

		switch err {
		case volume.ErrVolumeNotFound:
			code = codes.NotFound
		case volume.ErrServerNotFound:
			code = codes.NotFound
		case volume.ErrStillAttached:
			code = codes.FailedPrecondition
		case volume.ErrAttachLimitReached:
			code = codes.ResourceExhausted
		case volume.ErrLockedServer:
			code = codes.Unavailable
		}

		return nil, status.Error(code, fmt.Sprintf("failed to publish volume: %s", err))
	}

	return &csi.ControllerPublishVolumeResponse{}, nil
}

// ControllerUnpublishVolume implements the CSI standard definition.
func (s *Service) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid volume id")
	}

	if err := s.volume.Detach(ctx, volume.DetachOpts{
		Volume: &model.Volume{
			Name: req.VolumeId,
		},
		Server: &model.Server{
			Name: req.NodeId,
		},
	}); err != nil {
		code := codes.Internal

		switch err {
		case volume.ErrVolumeNotFound:
			return &csi.ControllerUnpublishVolumeResponse{}, nil
		case volume.ErrServerNotFound:
			code = codes.NotFound
		case volume.ErrLockedServer:
			code = codes.Unavailable
		}

		return nil, status.Error(code, fmt.Sprintf("failed to unpublish volume: %s", err))
	}

	return &csi.ControllerUnpublishVolumeResponse{}, nil
}
