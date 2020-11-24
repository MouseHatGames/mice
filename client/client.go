package client

import "github.com/MouseHatGames/mice/options"

type Client interface {
	Call(service string, path string, val interface{}, opts ...CallOption) (resp interface{}, err error)
}

type client struct {
	opts *options.Options
}

func newClient(opts *options.Options) Client {
	return &client{opts: opts}
}

func (c *client) Call(service string, path string, val interface{}, opts ...CallOption) (resp interface{}, err error) {
	var callopts CallOptions

	for _, o := range opts {
		o(&callopts)
	}

	return nil, nil
}
