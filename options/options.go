package options

import (
	"github.com/MouseHatGames/mice/codec"
	"github.com/MouseHatGames/mice/logger"
	"github.com/MouseHatGames/mice/router"
	"github.com/MouseHatGames/mice/transport"
)

// Options holds the configuration for a service instance
type Options struct {
	Name       string
	ListenAddr string
	Logger     logger.Logger
	Codec      codec.Codec
	Transport  transport.Transport
	Router     router.Router
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

// Logger sets the logger that will receive the log messages sent by the library
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// Codec sets the codec that will transform the messages
func Codec(c codec.Codec) Option {
	return func(o *Options) {
		o.Codec = c
	}
}

// Transport sets the transport that will deliver and receive messages to and from other services
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}

// Router sets the router
func Router(r router.Router) Option {
	return func(o *Options) {
		o.Router = r
	}
}
