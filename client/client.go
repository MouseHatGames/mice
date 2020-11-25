package client

import (
	"context"
	"fmt"

	"github.com/MouseHatGames/mice/options"
	"github.com/MouseHatGames/mice/transport"
)

type Client interface {
	Call(ctx context.Context, service string, path string, req interface{}, resp interface{}, opts ...CallOption) error
}

type client struct {
	opts *options.Options
}

type CallError struct {
	msg string
}

func (c *CallError) Error() string {
	return c.msg
}

func NewClient(opts *options.Options) Client {
	return &client{opts: opts}
}

func (c *client) Call(ctx context.Context, service string, path string, reqval interface{}, respval interface{}, opts ...CallOption) error {
	var callopts CallOptions

	for _, o := range opts {
		o(&callopts)
	}

	if callopts.Context == nil {
		callopts.Context = context.Background()
	}

	s, err := c.opts.Transport.Dial(ctx, service)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer s.Close()

	reqid := "TODO: Generate ID"
	req := transport.NewMessage()
	req.Headers[transport.HeaderRequestID] = reqid
	req.Headers[transport.HeaderPath] = path

	req.Data, err = c.opts.Codec.Marshal(reqval)
	if err != nil {
		return fmt.Errorf("encode request: %w", err)
	}

	if err := s.Send(callopts.Context, req); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	var respmsg transport.Message
	if err := s.Receive(ctx, &respmsg); err != nil {
		return fmt.Errorf("receive message: %w", err)
	}

	if err, ok := respmsg.Headers[transport.HeaderError]; ok {
		return &CallError{err}
	}

	fmt.Println(string(respmsg.Data))
	if err := c.opts.Codec.Unmarshal(respmsg.Data, respval); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}
