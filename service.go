package mice

import (
	"errors"

	"github.com/MouseHatGames/mice/options"
)

// Service represents a service that can receive and send requests
type Service interface {
	Apply(opts ...options.Option)
	Start() error
}

type service struct {
	options options.Options
}

// NewService instantiates a new service and initializes it with options
func NewService(opts ...options.Option) Service {
	svc := &service{}
	svc.Apply(opts...)

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

	return nil
}
