package options

import (
	"github.com/MouseHatGames/mice/broker"
	"github.com/MouseHatGames/mice/codec"
	"github.com/MouseHatGames/mice/config"
	"github.com/MouseHatGames/mice/logger"
	"github.com/MouseHatGames/mice/router"
	"github.com/MouseHatGames/mice/transport"
)

// Options holds the configuration for a service instance
type Options struct {
	Name    string
	RPCPort int16

	Logger    logger.Logger
	Codec     codec.Codec
	Transport transport.Transport
	Router    router.Router
	Broker    broker.Broker
	Config    config.Config
}

// Option represents a function that can be used to mutate an Options object
type Option func(*Options)

// Name sets the name of the service
func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// RPCPort sets the port in which this service's RPC will listen on, as well as the port in which other services' RPC servers are listening on
func RPCPort(port int16) Option {
	return func(o *Options) {
		o.RPCPort = port
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
