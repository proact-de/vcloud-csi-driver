package node

import (
	"github.com/proact-de/vcloud-csi-driver/pkg/service/mount"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/resize"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/stats"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/volume"
	"github.com/proact-de/vcloud-csi-driver/pkg/vcloud"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Client *vcloud.Client
	Volume *volume.Service
	Mount  *mount.Service
	Resize *resize.Service
	Stats  *stats.Service
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// WithClient provides a function to set the client option.
func WithClient(v *vcloud.Client) Option {
	return func(o *Options) {
		o.Client = v
	}
}

// WithVolume provides a function to set the volume option.
func WithVolume(v *volume.Service) Option {
	return func(o *Options) {
		o.Volume = v
	}
}

// WithMount provides a function to set the mount option.
func WithMount(v *mount.Service) Option {
	return func(o *Options) {
		o.Mount = v
	}
}

// WithResize provides a function to set the mount option.
func WithResize(v *resize.Service) Option {
	return func(o *Options) {
		o.Resize = v
	}
}

// WithStats provides a function to set the mount option.
func WithStats(v *stats.Service) Option {
	return func(o *Options) {
		o.Stats = v
	}
}
