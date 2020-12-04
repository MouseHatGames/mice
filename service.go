package mice

import (
	"errors"
	"flag"
	"fmt"
	"os"

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

	Env() options.Environment
	Config() config.Config

	Server() server.Server
	Client() client.Client

	Start() error
}

type starter interface {
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
	svc.options.RPCPort = options.DefaultRPCPort
	svc.options.Environment = getEnvironment()

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

func getEnvironment() options.Environment {
	arg := flag.String("env", "", "")
	flag.Parse()

	if *arg != "" {
		return options.ParseEnvironment(*arg)
	}

	if env, ok := os.LookupEnv("MICE_ENV"); ok {
		return options.ParseEnvironment(env)
	}

	return options.EnvironmentDevelopment
}

func (s *service) Apply(opts ...options.Option) {
	for _, o := range opts {
		o(&s.options)
	}
}

func (s *service) Env() options.Environment {
	return s.options.Environment
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

	if err := tryStart(map[string]interface{}{
		"broker":    s.options.Broker,
		"codec":     s.options.Codec,
		"config":    s.options.Config,
		"discovery": s.options.Discovery,
		"logger":    s.options.Logger,
		"transport": s.options.Transport,
	}); err != nil {
		return err
	}

	s.options.Logger.Infof("starting on %s environment", s.options.Environment)

	return s.server.Start()
}

func tryStart(objs map[string]interface{}) error {
	for k, v := range objs {
		if s, ok := v.(starter); ok {
			if err := s.Start(); err != nil {
				return fmt.Errorf("start %s: %w", k, err)
			}
		}
	}
	return nil
}

func (s *service) Server() server.Server {
	return s.server
}

func (s *service) Client() client.Client {
	return s.client
}
