package node

import (
	"context"
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/mount"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/volume"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NodePublishVolume implements the CSI standard definition.
func (s *Service) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing volume id")
	}

	if req.TargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "missing target path")
	}

	if req.StagingTargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "missing staging target path")
	}

	if req.VolumeCapability == nil {
		return nil, status.Error(codes.InvalidArgument, "missing volume capability")
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

	switch {
	case req.VolumeCapability.GetBlock() != nil:
		if err := s.mount.Publish(ctx, mount.PublishOpts{
			Volume:  v,
			Target:  req.TargetPath,
			Staging: v.Device,
			IsBlock: true,
		}); err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf(
				"failed to publish block volume: %s",
				err,
			))
		}

		return &csi.NodePublishVolumeResponse{}, nil
	case req.VolumeCapability.GetMount() != nil:
		capability := req.VolumeCapability.GetMount()

		if err := s.mount.Publish(ctx, mount.PublishOpts{
			Volume:     v,
			Target:     req.TargetPath,
			Staging:    req.StagingTargetPath,
			Readonly:   req.Readonly,
			FSType:     capability.FsType,
			MountFlags: capability.MountFlags,
		}); err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf(
				"failed to publish mount volume: %s",
				err,
			))
		}

		return &csi.NodePublishVolumeResponse{}, nil
	}

	return nil, status.Error(codes.InvalidArgument, "unsupported volume capability to publish")
}

// NodeUnpublishVolume implements the CSI standard definition.
func (s *Service) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing volume id")
	}

	if req.TargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "missing target path")
	}

	v, err := s.volume.Find(ctx, volume.FindOpts{
		ID: req.VolumeId,
	})

	if err != nil {
		switch err {
		case volume.ErrVolumeNotFound:
			return &csi.NodeUnpublishVolumeResponse{}, nil
		default:
			return nil, status.Error(codes.Internal, fmt.Sprintf(
				"failed to get volume: %s",
				err,
			))
		}
	}

	if err := s.mount.Unpublish(ctx, mount.UnpublishOpts{
		Volume: v,
		Target: req.TargetPath,
	}); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf(
			"failed to unpublish volume: %s",
			err,
		))
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

// NodeStageVolume implements the CSI standard definition.
func (s *Service) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing volume id")
	}

	if req.StagingTargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "missing staging target path")
	}

	if req.VolumeCapability == nil {
		return nil, status.Error(codes.InvalidArgument, "missing volume capability")
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

	switch {
	case req.VolumeCapability.GetBlock() != nil:
		return &csi.NodeStageVolumeResponse{}, nil
	case req.VolumeCapability.GetMount() != nil:
		capability := req.VolumeCapability.GetMount()

		if err := s.mount.Stage(ctx, mount.StageOpts{
			Volume:     v,
			Target:     req.StagingTargetPath,
			FSType:     capability.FsType,
			MountFlags: capability.MountFlags,
		}); err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf(
				"failed to stage volume: %s",
				err,
			))
		}

		return &csi.NodeStageVolumeResponse{}, nil
	}

	return nil, status.Error(codes.InvalidArgument, "stage volume: unsupported volume capability")
}

// NodeUnstageVolume implements the CSI standard definition.
func (s *Service) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing volume id")
	}

	if req.StagingTargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "missing staging target path")
	}

	v, err := s.volume.Find(ctx, volume.FindOpts{
		ID: req.VolumeId,
	})

	if err != nil {
		switch err {
		case volume.ErrVolumeNotFound:
			return &csi.NodeUnstageVolumeResponse{}, nil
		default:
			return nil, status.Error(codes.Internal, fmt.Sprintf(
				"failed to get volume: %s",
				err,
			))
		}
	}

	if err := s.mount.Unstage(ctx, mount.UnstageOpts{
		Volume: v,
		Target: req.StagingTargetPath,
	}); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf(
			"failed to unstage volume: %s",
			err,
		))
	}

	return &csi.NodeUnstageVolumeResponse{}, nil
}
