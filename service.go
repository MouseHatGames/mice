package mice

import (
	"context"
	"errors"
	"fmt"

	"github.com/MouseHatGames/mice/client"
	"github.com/MouseHatGames/mice/config"
	"github.com/MouseHatGames/mice/logger"
	"github.com/MouseHatGames/mice/options"
	"github.com/MouseHatGames/mice/router"
	"github.com/MouseHatGames/mice/server"
)

// Service represents a service that can receive and send requests
type Service interface {
	// Apply applies one or more options to the service's configuration
	Apply(opts ...options.Option)

	Config() config.Config

	Server() server.Server
	Client() client.Client

	Start() error
}

type service struct {
	options options.Options
	server  server.Server
	client  client.Client
}

// NewService instantiates a new service and initializes it with options
func NewService(opts ...options.Option) Service {
	svc := &service{}
	svc.options.Logger = logger.NewStdoutLogger()
	svc.options.RPCPort = 7070

	svc.Apply(opts...)

	if svc.options.Name == "" {
		panic("no name defined")
	}
	if svc.options.Codec == nil {
		panic("no codec defined")
	}
	if svc.options.Transport == nil {
		panic("no transport defined")
	}

	svc.options.Router = router.NewRouter(svc.options.Codec, svc.options.Logger)
	svc.server = server.NewServer(&svc.options)
	svc.client = client.NewClient(&svc.options)

	return svc
}

func (s *service) Apply(opts ...options.Option) {
	for _, o := range opts {
		o(&s.options)
	}
}

func (s *service) Config() config.Config {
	if s.options.Config == nil {
		panic("no config provider has been set up")
	}
	return s.options.Config
}

func (s *service) Start() error {
	if s.options.Name == "" {
		return errors.New("missing service name")
	}

	if s.options.Broker != nil {
		if err := s.options.Broker.Connect(context.Background()); err != nil {
			return fmt.Errorf("broker connect: %w", err)
		}
	}

	return s.server.Start()
}

func (s *service) Server() server.Server {
	return s.server
}

func (s *service) Client() client.Client {
	return s.client
}
