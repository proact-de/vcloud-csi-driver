package controller

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/volume"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// MinVolumeSize defines the minimum volume size.
	MinVolumeSize int64 = 1073741824

	// MaxVolumeSize defines the maximum volume size.
	MaxVolumeSize int64 = 10995116277760

	// DefaultVolumeSize defines the default volume size.
	DefaultVolumeSize int64 = 5368709120

	// TopologyKey defines the topology key for kubernetes.
	TopologyKey = "failure-domain.beta.kubernetes.io/zone"
)

// Service defines the service for the controller component.
type Service struct {
	server     string
	datacenter string
	volume     *volume.Service
}

// NewService simply initializes a new controller service.
func NewService(opts ...Option) *Service {
	options := newOptions(opts...)

	return &Service{
		server:     options.Server,
		datacenter: options.Datacenter,
		volume:     options.Volume,
	}
}

// ValidateVolumeCapabilities implements the CSI standard definition.
func (s *Service) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing volume id")
	}

	if len(req.VolumeCapabilities) == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing volume capabilities")
	}

	return &csi.ValidateVolumeCapabilitiesResponse{
		Confirmed: &csi.ValidateVolumeCapabilitiesResponse_Confirmed{
			VolumeCapabilities: []*csi.VolumeCapability{
				{
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
					},
				},
			},
		},
	}, nil
}

// ControllerGetCapabilities implements the CSI standard definition.
func (s *Service) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
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
