package identity

import (
	"context"
	"sync"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/proact-de/vcloud-csi-driver/pkg/version"
)

// Service defines the service for the identity component.
type Service struct {
	options Options
	mutex   sync.RWMutex
	ready   bool
}

// NewService simply initializes a new identity service.
func NewService(opts ...Option) *Service {
	return &Service{
		options: newOptions(opts...),
	}
}

// SetReady simply enables or disables the ready flag of the identity service.
func (s *Service) SetReady(ready bool) {
	s.mutex.Lock()
	s.ready = ready
	s.mutex.Unlock()
}

// GetPluginInfo implements the CSI standard definition.
func (s *Service) GetPluginInfo(context.Context, *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	return &csi.GetPluginInfoResponse{
		Name:          "vcloud.csi.proact.de",
		VendorVersion: version.String,
	}, nil
}

// GetPluginCapabilities implements the CSI standard definition.
func (s *Service) GetPluginCapabilities(context.Context, *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
					},
				},
			},
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_VOLUME_ACCESSIBILITY_CONSTRAINTS,
					},
				},
			},
			{
				Type: &csi.PluginCapability_VolumeExpansion_{
					VolumeExpansion: &csi.PluginCapability_VolumeExpansion{
						Type: csi.PluginCapability_VolumeExpansion_ONLINE,
					},
				},
			},
			{
				Type: &csi.PluginCapability_VolumeExpansion_{
					VolumeExpansion: &csi.PluginCapability_VolumeExpansion{
						Type: csi.PluginCapability_VolumeExpansion_OFFLINE,
					},
				},
			},
		},
	}, nil
}

// Probe implements the CSI standard definition.
func (s *Service) Probe(context.Context, *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	return &csi.ProbeResponse{
		Ready: &wrappers.BoolValue{
			Value: s.isReady(),
		},
	}, nil
}

func (s *Service) isReady() bool {
	s.mutex.RLock()
	ready := s.ready
	s.mutex.RUnlock()

	return ready
}
