package server

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/MouseHatGames/mice/logger"
	"github.com/MouseHatGames/mice/options"
	"github.com/MouseHatGames/mice/transport"
)

type Server interface {
	Start() error
	AddHandler(h interface{})
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
	l, err := s.opts.Transport.Listen(s.opts.ListenAddr)
	if err != nil {
		return err
	}

	if err := l.Accept(s.handle); err != nil {
		return fmt.Errorf("accept connections: %w", err)
	}

	return nil
}

func (s *server) AddHandler(h interface{}) {
	s.opts.Router.AddHandler(h)
}

func (s *server) handle(soc transport.Socket) {
	go func() {
		defer soc.Close()

		var req transport.Message

		for {
			err := soc.Receive(&req)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					s.log.Errorf("receive message: %s", err)
				} else {
					s.log.Debugf("socket eof")
				}
				break
			}

			//TODO: Fix potential race condition where the socket gets closed before this method sends the response
			go s.handleRequest(&req, soc)
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
