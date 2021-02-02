package stats

import (
	"golang.org/x/sys/unix"
)

// Service defines the service for the stats component.
type Service struct {
	options Options
}

// NewService simply initializes a new stats service.
func NewService(opts ...Option) *Service {
	return &Service{
		options: newOptions(opts...),
	}
}

// ByteFilesystemStats checks the available bytes of the filesystem.
func (s *Service) ByteFilesystemStats(path string) (int64, int64, int64, error) {
	statfs := &unix.Statfs_t{}

	if err := unix.Statfs(path, statfs); err != nil {
		return 0, 0, 0, err
	}

	totalBytes := int64(statfs.Blocks) * int64(statfs.Bsize)
	availableBytes := int64(statfs.Bavail) * int64(statfs.Bsize)
	usedBytes := (int64(statfs.Blocks) - int64(statfs.Bfree)) * int64(statfs.Bsize)

	return totalBytes, availableBytes, usedBytes, nil
}

// InodeFilesystemStats checks the available inodes of the filesystem.
func (s *Service) InodeFilesystemStats(path string) (int64, int64, int64, error) {
	statfs := &unix.Statfs_t{}

	if err := unix.Statfs(path, statfs); err != nil {
		return 0, 0, 0, nil
	}

	total := int64(statfs.Files)
	used := int64(statfs.Files) - int64(statfs.Ffree)
	free := int64(statfs.Ffree)

	return total, used, free, nil
}
