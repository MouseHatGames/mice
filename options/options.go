package options

import (
	"github.com/MouseHatGames/mice/broker"
	"github.com/MouseHatGames/mice/codec"
	"github.com/MouseHatGames/mice/config"
	"github.com/MouseHatGames/mice/discovery"
	"github.com/MouseHatGames/mice/logger"
	"github.com/MouseHatGames/mice/transport"
	"go.opentelemetry.io/otel/trace"
)

// Options holds the configuration for a service instance
type Options struct {
	Name        string
	RPCPort     int16
	Environment Environment

	Logger    logger.Logger
	Codec     codec.Codec
	Transport transport.Transport
	Broker    broker.Broker
	Config    config.Config
	Discovery discovery.Discovery
	Tracer    trace.Tracer
}

// DefaultRPCPort is the port that will be used for RPC connections if no other is specified
const DefaultRPCPort = 7070

// Option represents a function that can be used to mutate an Options object
type Option func(*Options)

// Name sets the name of the service
func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// RPCPort sets the port in which this service's RPC will listen on, as well as the port in which other services' RPC servers are listening on.
// Defaults to DefaultRPCPort
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

// WithEnvironment sets the environment that the service is running under
func WithEnvironment(e Environment) Option {
	return func(o *Options) {
		o.Environment = e
	}
}

// IfEnvironment conditionally applies options if running under a specific environment.
//
// Make sure this option comes after any WithEnvironment options, if any.
func IfEnvironment(env Environment, opts ...Option) Option {
	return func(o *Options) {
		if o.Environment == env {
			for _, opt := range opts {
				opt(o)
			}
		}
	}
}

// IfNotEnvironment conditionally applies options if running under any environment except for a specific one.
//
// Make sure this option comes after any WithEnvironment options, if any.
func IfNotEnvironment(env Environment, opts ...Option) Option {
	return func(o *Options) {
		if o.Environment != env {
			for _, opt := range opts {
				opt(o)
			}
		}
	}
}

// Tracer sets the OpenTelemetry tracer to use.
func Tracer(tracer trace.Tracer) Option {
	return func(o *Options) {
		o.Tracer = tracer
	}
}
