package volume

import (
	"context"
	"errors"
	"fmt"

	"github.com/proact-de/vcloud-csi-driver/pkg/model"
	"github.com/proact-de/vcloud-csi-driver/pkg/vcloud"
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

var (
	// ErrVolumeNotFound TODO
	ErrVolumeNotFound = errors.New("volume not found")

	// ErrVolumeAlreadyExists TODO
	ErrVolumeAlreadyExists = errors.New("volume does already exist")

	// ErrServerNotFound TODO
	ErrServerNotFound = errors.New("server not found")

	// ErrStillAttached TODO
	ErrStillAttached = errors.New("volume is attached")

	// ErrNotAttached TODO
	ErrNotAttached = errors.New("volume is not attached")

	// ErrAttachLimitReached TODO
	ErrAttachLimitReached = errors.New("max number of attachments per server reached")

	// ErrLockedServer TODO
	ErrLockedServer = errors.New("server is locked")
)

// Service defines the service for the resize component.
type Service struct {
	client *vcloud.Client
}

// NewService simply initializes a new resize service.
func NewService(opts ...Option) *Service {
	options := newOptions(opts...)

	return &Service{
		client: options.Client,
	}
}

// FindOpts defines the available options for the find handler.
type FindOpts struct {
	ID string
}

// Find simply tries to find the requested volume.
func (s *Service) Find(ctx context.Context, opts FindOpts) (*model.Volume, error) {
	disk, err := s.client.Find(opts.ID)

	if err != nil {
		if errors.Is(err, govcd.ErrorEntityNotFound) {
			return nil, ErrVolumeNotFound
		}
	}

	return disk, nil
}

// CreateOpts defines the available options for the create handler.
type CreateOpts struct {
	Name    string
	Size    int64
	Profile string
}

// Create is handling the creation of a volume for a server.
func (s *Service) Create(ctx context.Context, opts CreateOpts) error {
	return fmt.Errorf("not implemented")
}

// DeleteOpts defines the available options for the delete handler.
type DeleteOpts struct {
	ID string
}

// Delete is handling the deletion of a volume from a server.
func (s *Service) Delete(ctx context.Context, opts DeleteOpts) error {
	return fmt.Errorf("not implemented")
}

// AttachOpts defines the available options for the attach handler.
type AttachOpts struct {
	Volume *model.Volume
	Server *model.Server
}

// Attach is handling the attaching of a volume to a server.
func (s *Service) Attach(ctx context.Context, opts AttachOpts) error {
	return nil
}

// DetachOpts defines the available options for the detach handler.
type DetachOpts struct {
	Volume *model.Volume
	Server *model.Server
}

// Detach is handling the detaching of a volume from a server.
func (s *Service) Detach(ctx context.Context, opts DetachOpts) error {
	return nil
}

// ResizeOpts defines the available options for the resize handler.
type ResizeOpts struct {
	Volume *model.Volume
	Size   int64
}

// Resize is handling the resizing of a volume from a server.
func (s *Service) Resize(ctx context.Context, opts ResizeOpts) error {
	return nil
}
