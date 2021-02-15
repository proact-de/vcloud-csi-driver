package mount

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/proact-de/vcloud-csi-driver/pkg/model"
	"github.com/rs/zerolog/log"
	"k8s.io/mount-utils"
	"k8s.io/utils/exec"
)

const (
	// DefaultFilesystem defines the standard filesystem if nothing else is provided.
	DefaultFilesystem = "ext4"
)

// Service defines the service for the mount component.
type Service struct {
	options Options
	mounter *mount.SafeFormatAndMount
}

// NewService simply initializes a new mount service.
func NewService(opts ...Option) *Service {
	options := newOptions(opts...)

	return &Service{
		options: options,
		mounter: &mount.SafeFormatAndMount{
			Interface: mount.New(""),
			Exec:      exec.New(),
		},
	}
}

// PublishOpts defines the available options for the publish handler.
type PublishOpts struct {
	Volume     *model.Volume
	Target     string
	Staging    string
	IsBlock    bool
	Readonly   bool
	FSType     string
	MountFlags []string
}

// Publish is publishing the request volume.
func (s *Service) Publish(ctx context.Context, opts PublishOpts) error {
	if opts.IsBlock {
		if err := os.MkdirAll(
			filepath.Dir(opts.Target),
			os.FileMode(0750),
		); err != nil {
			return err
		}

		mountFile, err := os.OpenFile(
			opts.Target,
			os.O_CREATE,
			os.FileMode(0660),
		)

		if err != nil {
			return err
		}

		mountFile.Close()
	} else {
		if opts.FSType == "" {
			opts.FSType = DefaultFilesystem
		}

		if err := os.MkdirAll(
			opts.Target,
			os.FileMode(0750),
		); err != nil {
			return err
		}
	}

	options := []string{
		"bind",
	}

	if opts.Readonly {
		options = append(options, "ro")
	}

	options = append(options, opts.MountFlags...)

	log.Info().
		Str("name", opts.Volume.Name).
		Str("target", opts.Target).
		Str("staging", opts.Staging).
		Str("fs", opts.FSType).
		Bool("block", opts.IsBlock).
		Str("options", strings.Join(options, ", ")).
		Msg("Publishing volume")

	if err := s.mounter.Interface.Mount(
		opts.Staging,
		opts.Target,
		opts.FSType,
		options,
	); err != nil {
		log.Error().
			Err(err).
			Str("name", opts.Volume.Name).
			Str("target", opts.Target).
			Str("staging", opts.Staging).
			Str("fs", opts.FSType).
			Bool("block", opts.IsBlock).
			Str("options", strings.Join(options, ", ")).
			Msg("Failed to publish volume")

		return err
	}

	return nil
}

// UnpublishOpts defines the available options for the unpublish handler.
type UnpublishOpts struct {
	Volume *model.Volume
	Target string
}

// Unpublish is unpublishing the request volume.
func (s *Service) Unpublish(ctx context.Context, opts UnpublishOpts) error {
	log.Info().
		Str("name", opts.Volume.Name).
		Str("target", opts.Target).
		Msg("Unpublishing volume")

	if err := s.mounter.Interface.Unmount(opts.Target); err != nil {
		log.Error().
			Err(err).
			Str("name", opts.Volume.Name).
			Str("target", opts.Target).
			Msg("Failed to unpublish volume")

		return err
	}

	return nil
}

// StageOpts defines the available options for the stage handler.
type StageOpts struct {
	Volume     *model.Volume
	Target     string
	FSType     string
	MountFlags []string
}

// Stage is staging the request volume.
func (s *Service) Stage(ctx context.Context, opts StageOpts) error {
	if opts.FSType == "" {
		opts.FSType = DefaultFilesystem
	}

	log.Info().
		Str("name", opts.Volume.Name).
		Str("target", opts.Target).
		Str("fs", opts.FSType).
		Msg("Staging volume")

	isNotMountPoint, err := s.mounter.Interface.IsLikelyNotMountPoint(opts.Target)

	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(
				opts.Target,
				os.FileMode(0750),
			); err != nil {
				return err
			}

			isNotMountPoint = true
		} else {
			return err
		}
	}

	if !isNotMountPoint {
		return fmt.Errorf("%q is not a valid mount point", opts.Target)
	}

	if err := s.mounter.FormatAndMount(
		opts.Volume.Device,
		opts.Target,
		opts.FSType,
		nil,
	); err != nil {
		log.Error().
			Err(err).
			Str("name", opts.Volume.Name).
			Str("target", opts.Target).
			Str("fs", opts.FSType).
			Msg("Failed to stage volume")

		return err
	}

	return nil
}

// UnstageOpts defines the available options for the unstage handler.
type UnstageOpts struct {
	Volume *model.Volume
	Target string
}

// Unstage is unstaging the request volume.
func (s *Service) Unstage(ctx context.Context, opts UnstageOpts) error {
	log.Info().
		Str("name", opts.Volume.Name).
		Str("target", opts.Target).
		Msg("Unstaging volume")

	if err := s.mounter.Interface.Unmount(opts.Target); err != nil {
		log.Error().
			Err(err).
			Str("name", opts.Volume.Name).
			Str("target", opts.Target).
			Msg("Failed to unstage volume")

		return err
	}

	return nil
}

// Exists is checking if volume path exists.
func (s *Service) Exists(ctx context.Context, target string) (bool, error) {
	_, err := os.Stat(target)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
