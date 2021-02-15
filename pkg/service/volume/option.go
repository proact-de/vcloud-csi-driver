package volume

import (
	"github.com/proact-de/vcloud-csi-driver/pkg/vcloud"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Client *vcloud.Client
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
