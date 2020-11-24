package options

import "github.com/MouseHatGames/mice/codec"

// Options holds the configuration for a service instance
type Options struct {
	Name       string
	ListenAddr string
	Codec      codec.Codec
}

// Option represents a function that can be used to mutate an Options object
type Option func(*Options)

// Name sets the name of the service
func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// ListenAddr sets the address in which the gRPC server will listen on
func ListenAddr(addr string) Option {
	return func(o *Options) {
		o.ListenAddr = addr
	}
}
