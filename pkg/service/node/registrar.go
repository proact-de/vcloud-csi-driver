package node

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

// NodeGetInfo implements the CSI standard definition.
func (s *Service) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	return &csi.NodeGetInfoResponse{
		NodeId:            s.server,
		MaxVolumesPerNode: MaxVolumesPerNode,
		AccessibleTopology: &csi.Topology{
			Segments: map[string]string{
				"failure-domain.beta.kubernetes.io/zone": s.datacenter,
			},
		},
	}, nil
}
