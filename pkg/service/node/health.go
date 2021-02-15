package node

import (
	"context"
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	totalBytes, availableBytes, usedBytes, err := s.byteStats(req.VolumePath)

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf(
			"failed to get volume byte stats: %s",
			err,
		))
	}

	totalINodes, usedINodes, freeINodes, err := s.inodeStats(req.VolumePath)

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf(
			"failed to get volume inode stats: %s",
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
				Available: freeINodes,
				Total:     totalINodes,
				Used:      usedINodes,
			},
		},
	}, nil
}

// byteStats checks the available bytes of the filesystem.
func (s *Service) byteStats(path string) (int64, int64, int64, error) {
	statfs := &unix.Statfs_t{}

	if err := unix.Statfs(path, statfs); err != nil {
		return 0, 0, 0, err
	}

	total := int64(statfs.Blocks) * int64(statfs.Bsize)
	available := int64(statfs.Bavail) * int64(statfs.Bsize)
	used := (int64(statfs.Blocks) - int64(statfs.Bfree)) * int64(statfs.Bsize)

	return total, available, used, nil
}

// inodeStats checks the available inodes of the filesystem.
func (s *Service) inodeStats(path string) (int64, int64, int64, error) {
	statfs := &unix.Statfs_t{}

	if err := unix.Statfs(path, statfs); err != nil {
		return 0, 0, 0, nil
	}

	total := int64(statfs.Files)
	used := int64(statfs.Files) - int64(statfs.Ffree)
	free := int64(statfs.Ffree)

	return total, used, free, nil
}
