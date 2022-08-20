package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/MouseHatGames/mice/logger"
	"github.com/MouseHatGames/mice/transport"
)

type httpIncomingSocket struct {
	rw              http.ResponseWriter
	r               *http.Request
	log             logger.Logger
	closer          chan<- struct{}
	sentResponse    bool
	receivedRequest bool
}

var _ transport.Socket = (*httpIncomingSocket)(nil)

func (s *httpIncomingSocket) Close() error {
	s.log.Debugf("closing incoming socket")
	s.closer <- struct{}{}
	return nil
}

func (s *httpIncomingSocket) Send(ctx context.Context, msg *transport.Message) error {
	if s.sentResponse {
		return errors.New("response already sent")
	}
	s.sentResponse = true

	s.log.Debugf("sending response with %d bytes", len(msg.Data))

	if err := marshalMessage(s.rw, msg); err != nil {
		return fmt.Errorf("encode message: %w", err)
	}

	return nil
}

func (s *httpIncomingSocket) Receive(ctx context.Context, msg *transport.Message) error {
	if s.receivedRequest {
		return io.EOF
	}
	s.receivedRequest = true

	if err := unmarshalMessage(s.r.Body, msg); err != nil {
		return fmt.Errorf("read message: %w", err)
	}

	s.log.Debugf("received request with %d bytes", len(msg.Data))

	return nil
}
