package resize

import (
	"context"

	"github.com/proact-de/vcloud-csi-driver/pkg/model"
	"github.com/rs/zerolog/log"
	"k8s.io/kubernetes/pkg/util/resizefs"
	"k8s.io/mount-utils"
	"k8s.io/utils/exec"
)

// Service defines the service for the resize component.
type Service struct {
	options Options
	resizer *resizefs.ResizeFs
}

// NewService simply initializes a new resize service.
func NewService(opts ...Option) *Service {
	options := newOptions(opts...)

	return &Service{
		options: options,
		resizer: resizefs.NewResizeFs(&mount.SafeFormatAndMount{
			Interface: mount.New(""),
			Exec:      exec.New(),
		}),
	}
}

// Opts defines the available options for the resize handler.
type Opts struct {
	Volume *model.Volume
	Target string
}

// Resize is resizing the requested volume.
func (s *Service) Resize(ctx context.Context, opts Opts) error {
	log.Debug().
		Str("name", opts.Volume.Name).
		Str("target", opts.Target).
		Str("device", opts.Volume.Device).
		Msg("Resizing volume")

	if _, err := s.resizer.Resize(opts.Volume.Device, opts.Target); err != nil {
		log.Error().
			Err(err).
			Str("name", opts.Volume.Name).
			Str("target", opts.Target).
			Str("device", opts.Volume.Device).
			Msg("Failed to resize volume")

		return err
	}

	return nil
}
