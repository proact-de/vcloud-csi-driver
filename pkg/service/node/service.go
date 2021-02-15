package node

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/mount"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/resize"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/volume"
)

const (
	// MaxVolumesPerNode defines the maximum volumes per node.
	MaxVolumesPerNode int64 = 50
)

// Service defines the service for the node component.
type Service struct {
	server     string
	datacenter string
	volume     *volume.Service
	mount      *mount.Service
	resize     *resize.Service
}

// NewService simply initializes a new node service.
func NewService(opts ...Option) *Service {
	options := newOptions(opts...)

	return &Service{
		server:     options.Server,
		datacenter: options.Datacenter,
		volume:     options.Volume,
		mount:      options.Mount,
		resize:     options.Resize,
	}
}

// NodeGetCapabilities implements the CSI standard definition.
func (s *Service) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
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
						Type: csi.NodeServiceCapability_RPC_GET_VOLUME_STATS,
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
		},
	}, nil
}
