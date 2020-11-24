package mice

import (
	"errors"

	"github.com/MouseHatGames/mice/client"
	"github.com/MouseHatGames/mice/logger"
	"github.com/MouseHatGames/mice/options"
	"github.com/MouseHatGames/mice/router"
	"github.com/MouseHatGames/mice/server"
)

// Service represents a service that can receive and send requests
type Service interface {
	// Apply applies one or more options to the service's configuration
	Apply(opts ...options.Option)

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

	svc.Apply(opts...)

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

func (s *service) Start() error {
	if s.options.Name == "" {
		return errors.New("missing service name")
	}
	if s.options.ListenAddr == "" {
		return errors.New("missing listen address")
	}

	return s.server.Start()
}

func (s *service) Server() server.Server {
	return s.server
}

func (s *service) Client() client.Client {
	return s.client
}
