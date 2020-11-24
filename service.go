package mice

import (
	"errors"
	"fmt"
	"net"

	"github.com/MouseHatGames/mice/logger"
	"github.com/MouseHatGames/mice/options"
	"github.com/MouseHatGames/mice/router"
	"github.com/MouseHatGames/mice/server"
	"google.golang.org/grpc"
)

// Service represents a service that can receive and send requests
type Service interface {
	// Apply applies one or more options to the service's configuration
	Apply(opts ...options.Option)

	Server() server.Server

	Start() error
}

type service struct {
	options options.Options
	server  server.Server
}

// NewService instantiates a new service and initializes it with options
func NewService(opts ...options.Option) Service {
	svc := &service{}
	svc.options.Logger = logger.NewStdoutLogger()

	svc.Apply(opts...)

	svc.options.Router = router.NewRouter(svc.options.Codec)

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

	lis, err := net.Listen("tcp", s.options.ListenAddr)
	if err != nil {
		return fmt.Errorf("failed to open tcp listener: %w", err)
	}

	srv := grpc.NewServer()

	if err := srv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve grpc: %w", err)
	}

	return nil
}

func (s *service) Server() server.Server {
	return nil
}
