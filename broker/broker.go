package broker

import "context"

type Broker interface {
	Connect(ctx context.Context) error
	Close() error
	Publish(ctx context.Context, topic string, data interface{}) error
	Subscribe(ctx context.Context, topic string, callback func(interface{})) error
}
