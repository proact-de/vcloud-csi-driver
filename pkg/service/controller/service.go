package controller

import (
	"context"
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service defines the service for the controller component.
type Service struct {
	options Options
}

// NewService simply initializes a new controller service.
func NewService(opts ...Option) *Service {
	return &Service{
		options: newOptions(opts...),
	}
}

// GetCapacity implements the CSI standard definition.
func (s *Service) GetCapacity(ctx context.Context, req *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "GetCapacity").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// ValidateVolumeCapabilities implements the CSI standard definition.
func (s *Service) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "ValidateVolumeCapabilities").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// ListVolumes implements the CSI standard definition.
func (s *Service) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "ListVolumes").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// CreateVolume implements the CSI standard definition.
func (s *Service) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "CreateVolume").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// DeleteVolume implements the CSI standard definition.
func (s *Service) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "DeleteVolume").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// ListSnapshots implements the CSI standard definition.
func (s *Service) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "ListSnapshots").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// CreateSnapshot implements the CSI standard definition.
func (s *Service) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "CreateSnapshot").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// DeleteSnapshot implements the CSI standard definition.
func (s *Service) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "DeleteSnapshot").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// ControllerGetVolume implements the CSI standard definition.
func (s *Service) ControllerGetVolume(ctx context.Context, req *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "ControllerGetVolume").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// ControllerExpandVolume implements the CSI standard definition.
func (s *Service) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "ControllerExpandVolume").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// ControllerPublishVolume implements the CSI standard definition.
func (s *Service) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "ControllerPublishVolume").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// ControllerUnpublishVolume implements the CSI standard definition.
func (s *Service) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "ControllerUnpublishVolume").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// ControllerGetCapabilities implements the CSI standard definition.
func (s *Service) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "ControllerGetCapabilities").
		Msg("")

	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: []*csi.ControllerServiceCapability{
			{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
					},
				},
			},
			{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
					},
				},
			},
			{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_EXPAND_VOLUME,
					},
				},
			},
		},
	}, nil
}
