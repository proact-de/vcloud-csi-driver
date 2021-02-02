package vcloud

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Href        string
	Insecure    bool
	Username    string
	Password    string
	Org         string
	Datacenters []string
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
func WithHref(v string) Option {
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

// WithOrg provides a function to set the org option.
func WithOrg(v string) Option {
	return func(o *Options) {
		o.Org = v
	}
}

// WithDatacenters provides a function to set the datacenters option.
func WithDatacenters(v []string) Option {
	return func(o *Options) {
		o.Datacenters = v
	}
}
