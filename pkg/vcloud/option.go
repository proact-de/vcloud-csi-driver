package vcloud

import (
	"net/url"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Href         *url.URL
	Insecure     bool
	Username     string
	Password     string
	Organization string
	Datacenter   string
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// WithHref provides a function to set the href option.
func WithHref(v *url.URL) Option {
	return func(o *Options) {
		o.Href = v
	}
}

// WithInsecure provides a function to set the insecure option.
func WithInsecure(v bool) Option {
	return func(o *Options) {
		o.Insecure = v
	}
}

// WithUsername provides a function to set the username option.
func WithUsername(v string) Option {
	return func(o *Options) {
		o.Username = v
	}
}

// WithPassword provides a function to set the password option.
func WithPassword(v string) Option {
	return func(o *Options) {
		o.Password = v
	}
}

// WithOrganization provides a function to set the organization option.
func WithOrganization(v string) Option {
	return func(o *Options) {
		o.Organization = v
	}
}

// WithDatacenter provides a function to set the datacenter option.
func WithDatacenter(v string) Option {
	return func(o *Options) {
		o.Datacenter = v
	}
}
