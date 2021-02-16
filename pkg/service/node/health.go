package node

import (
	"context"
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/kubernetes/pkg/volume/util/fs"
)

// NodeGetVolumeStats implements the CSI standard definition.
func (s *Service) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing volume id")
	}

	if req.VolumePath == "" {
		return nil, status.Error(codes.InvalidArgument, "missing volume path")
	}

	exists, err := s.mount.Exists(ctx, req.VolumePath)

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf(
			"failed to check for volume existence: %s",
			err,
		))
	}

	if !exists {
		return nil, status.Error(codes.Unavailable, fmt.Sprintf(
			"volume %s is not available on %v",
			req.VolumePath,
			s.server,
		))
	}

	availableBytes, totalBytes, usedBytes, totalInodes, freeInodes, usedInodes, err := fs.FsInfo(req.VolumePath)

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf(
			"failed to get volume stats: %s",
			err,
		))
	}

	return &csi.NodeGetVolumeStatsResponse{
		Usage: []*csi.VolumeUsage{
			{
				Unit:      csi.VolumeUsage_BYTES,
				Available: availableBytes,
				Total:     totalBytes,
				Used:      usedBytes,
			},
			{
				Unit:      csi.VolumeUsage_INODES,
				Available: freeInodes,
				Total:     totalInodes,
				Used:      usedInodes,
			},
		},
	}, nil
}
