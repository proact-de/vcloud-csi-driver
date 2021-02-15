package node

import (
	"github.com/proact-de/vcloud-csi-driver/pkg/service/mount"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/resize"
	"github.com/proact-de/vcloud-csi-driver/pkg/service/volume"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Server     string
	Datacenter string
	Volume     *volume.Service
	Mount      *mount.Service
	Resize     *resize.Service
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// WithServer provides a function to set the server option.
func WithServer(v string) Option {
	return func(o *Options) {
		o.Server = v
	}
}

// WithDatacenter provides a function to set the datacenter option.
func WithDatacenter(v string) Option {
	return func(o *Options) {
		o.Datacenter = v
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

// WithResize provides a function to set the resize option.
func WithResize(v *resize.Service) Option {
	return func(o *Options) {
		o.Resize = v
	}
}
