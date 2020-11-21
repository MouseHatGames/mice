package transport

type Transport interface {
	Listen(addr string) error
	Send(msg *Message) error
}
