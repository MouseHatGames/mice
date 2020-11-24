package transport

import (
	"errors"
	"net"
)

var ErrTooManyHeaders = errors.New("too many headers")
var ErrHeaderTooLong = errors.New("header is longer than 255 characters")

type Transport interface {
	Listen(addr string) (Listener, error)
	Dial(addr string) (Socket, error)
}

type Socket interface {
	Close() error
	Send(msg *Message) error
	Receive(msg *Message) error
}

type Listener interface {
	Close() error
	Addr() net.Addr
	Accept(fn func(Socket)) error
}
