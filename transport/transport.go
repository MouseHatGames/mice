package transport

import (
	"context"
	"errors"
)

var ErrTooManyHeaders = errors.New("too many headers")
var ErrHeaderTooLong = errors.New("header is longer than 255 characters")

type Transport interface {
	Listen(ctx context.Context, addr string) (Listener, error)
	Dial(ctx context.Context, addr string) (Socket, error)
}

type Socket interface {
	Close() error
	Send(ctx context.Context, msg *Message) error
	Receive(ctx context.Context, msg *Message) error
}

type Listener interface {
	Close() error
	Accept(ctx context.Context, fn func(Socket)) error
}
