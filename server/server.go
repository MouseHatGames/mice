package server

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/MouseHatGames/mice/broker"
	"github.com/MouseHatGames/mice/logger"
	"github.com/MouseHatGames/mice/options"
	"github.com/MouseHatGames/mice/transport"
)

type Server interface {
	Start() error
	AddHandler(h interface{}, methods ...string)
	Publish(ctx context.Context, topic string, data interface{}) error
}

type server struct {
	opts *options.Options
	log  logger.Logger
}

func NewServer(opts *options.Options) Server {
	return &server{
		opts: opts,
		log:  opts.Logger.GetLogger("server"),
	}
}

func (s *server) Start() error {
	ctx := context.Background()

	l, err := s.opts.Transport.Listen(ctx, fmt.Sprintf(":%d", s.opts.RPCPort))
	if err != nil {
		return err
	}

	if err := l.Accept(ctx, s.handle); err != nil {
		return fmt.Errorf("accept connections: %w", err)
	}

	return nil
}

func (s *server) AddHandler(h interface{}, methods ...string) {
	s.opts.Router.AddHandler(h, methods)
}

func (s *server) handle(soc transport.Socket) {
	go func() {
		defer soc.Close()

		var req transport.Message
		ctx := context.Background()

		for {
			err := soc.Receive(ctx, &req)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					s.log.Errorf("receive message: %s", err)
				} else {
					s.log.Debugf("socket eof")
				}
				break
			}

			s.handleRequest(&req, soc)
		}
	}()
}

func (s *server) handleRequest(req *transport.Message, soc transport.Socket) {
	path, ok := req.Headers[transport.HeaderPath]
	if !ok {
		s.log.Errorf("missing path header")
		return
	}

	var resp transport.Message
	resp.Headers = map[string]string{
		transport.HeaderRequestID: req.Headers[transport.HeaderRequestID],
	}

	ret, err := s.opts.Router.Handle(path, req.Data)

	if err != nil {
		resp.SetError(err)
	} else {
		resp.Data = ret
	}

	if err := soc.Send(context.Background(), &resp); err != nil {
		s.log.Errorf("send response: %s", err)
	}
}

func (s *server) Publish(ctx context.Context, topic string, data interface{}) error {
	if s.opts.Broker == nil {
		panic("no broker has been declared")
	}

	b, err := s.opts.Codec.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}

	if err := s.opts.Broker.Publish(ctx, topic, &broker.Message{Data: b}); err != nil {
		return fmt.Errorf("publish message: %w", err)
	}

	return nil
}
