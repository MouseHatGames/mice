package broker

import "context"

type Broker interface {
	Close() error
	Publish(ctx context.Context, topic string, data *Message) error
	Subscribe(ctx context.Context, topic string, callback func(*Message)) error
}

type Message struct {
	Data []byte
}
