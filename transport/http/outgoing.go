package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MouseHatGames/mice/logger"
	"github.com/MouseHatGames/mice/transport"
)

type httpOutgoingSocket struct {
	address string
	resp    chan *http.Response
	log     logger.Logger
}

var _ transport.Socket = (*httpOutgoingSocket)(nil)

func (s *httpOutgoingSocket) Close() error {
	s.log.Debugf("closing outgoing socket")
	return nil
}

func (s *httpOutgoingSocket) Send(ctx context.Context, msg *transport.Message) error {
	s.log.Debugf("sending request with %d bytes", len(msg.Data))

	b, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("encode message: %w", err)
	}
	br := bytes.NewReader(b)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://%s/request", s.address), br)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	s.resp <- resp

	return nil
}

func (s *httpOutgoingSocket) Receive(ctx context.Context, msg *transport.Message) error {
	resp, ok := <-s.resp
	if !ok {
		return io.EOF
	}

	if err := unmarshalMessage(resp.Body, msg); err != nil {
		return fmt.Errorf("read message: %w", err)
	}

	s.log.Debugf("received response with %d bytes", len(msg.Data))

	return nil
}
