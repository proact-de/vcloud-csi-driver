package node

import (
	"context"
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service defines the service for the node component.
type Service struct {
	options Options
}

// NewService simply initializes a new node service.
func NewService(opts ...Option) *Service {
	return &Service{
		options: newOptions(opts...),
	}
}

// NodePublishVolume implements the CSI standard definition.
func (s *Service) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "NodePublishVolume").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// NodeUnpublishVolume implements the CSI standard definition.
func (s *Service) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "NodeUnpublishVolume").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// NodeStageVolume implements the CSI standard definition.
func (s *Service) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "NodeStageVolume").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// NodeUnstageVolume implements the CSI standard definition.
func (s *Service) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "NodeUnstageVolume").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// NodeExpandVolume implements the CSI standard definition.
func (s *Service) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "NodeExpandVolume").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// NodeGetVolumeStats implements the CSI standard definition.
func (s *Service) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "NodeGetVolumeStats").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// NodeGetInfo implements the CSI standard definition.
func (s *Service) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "NodeGetInfo").
		Msg("")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// NodeGetCapabilities implements the CSI standard definition.
func (s *Service) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	log.Info().
		Str("req", fmt.Sprintf("%+v", req)).
		Str("method", "NodeGetCapabilities").
		Msg("")

	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: []*csi.NodeServiceCapability{
			{
				Type: &csi.NodeServiceCapability_Rpc{
					Rpc: &csi.NodeServiceCapability_RPC{
						Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
					},
				},
			},
			{
				Type: &csi.NodeServiceCapability_Rpc{
					Rpc: &csi.NodeServiceCapability_RPC{
						Type: csi.NodeServiceCapability_RPC_EXPAND_VOLUME,
					},
				},
			},
			{
				Type: &csi.NodeServiceCapability_Rpc{
					Rpc: &csi.NodeServiceCapability_RPC{
						Type: csi.NodeServiceCapability_RPC_GET_VOLUME_STATS,
					},
				},
			},
		},
	}, nil
}
