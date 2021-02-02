package controller

import (
	"github.com/proact-de/vcloud-csi-driver/pkg/service/volume"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Volume *volume.Service
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// WithVolume provides a function to set the volume option.
func WithVolume(v *volume.Service) Option {
	return func(o *Options) {
		o.Volume = v
	}
}
