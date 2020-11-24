package server

import (
	"context"
	"errors"

	"github.com/MouseHatGames/mice/options"
	"github.com/MouseHatGames/mice/router"
	"github.com/MouseHatGames/mice/transport"
)

type Server interface {
	Start() error
}

type server struct {
	opts *options.Options
}

func NewServer(opts *options.Options) Server {
	return &server{
		opts: opts,
	}
}

func (s *server) Start() error {
	l, err := s.opts.Transport.Listen(s.opts.ListenAddr)
	if err != nil {
		return err
	}

	go l.Accept(s.handle)

	return nil
}

func (s *server) handle(soc transport.Socket) {
	defer soc.Close()

	go func() {
		var req transport.Message

		for {
			err := soc.Receive(&req)
			if err != nil {
				s.opts.Logger.Errorf("receive message: %s", err)
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
		s.opts.Logger.Errorf("missing path header")
		return
	}

	ret, err := s.opts.Router.Handle(path, req.Data)

	var resp transport.Message
	resp.Headers = map[string]string{
		transport.HeaderRequestID: req.Headers[transport.HeaderRequestID],
	}

	var hdlrerr router.HandlerError

	if errors.As(err, &hdlrerr) {
		resp.SetError(err)
	} else {
		d, err := s.opts.Codec.Marshal(ret)
		if err != nil {
			resp.SetError(err)
		} else {
			resp.Data = d
		}
	}

	if err := soc.Send(context.Background(), &resp); err != nil {
		s.opts.Logger.Errorf("send response: %s", err)
	}
}
