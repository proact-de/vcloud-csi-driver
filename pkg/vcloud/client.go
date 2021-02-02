package vcloud

// Client is a wrapper around the upstream vCloud Director client.
type Client struct {
	options Options
}

// NewClient simply initializes a new client wrapper.
func NewClient(opts ...Option) *Client {
	return &Client{
		options: newOptions(opts...),
	}
}
