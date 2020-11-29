package client

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/MouseHatGames/mice/broker"
	"github.com/MouseHatGames/mice/codec"
	"github.com/MouseHatGames/mice/discovery"
	"github.com/MouseHatGames/mice/logger"
	"github.com/MouseHatGames/mice/options"
	"github.com/MouseHatGames/mice/transport"
)

var ErrMustBeFunc = errors.New("value must be a function")
var ErrInvalidInput = errors.New("func must have 1 input")
var ErrInputPointer = errors.New("the func must take a pointer as an input")

type Client interface {
	Call(ctx context.Context, service string, path string, req interface{}, resp interface{}, opts ...CallOption) error
	Subscribe(topic string, callback interface{})
}

type client struct {
	codec  codec.Codec
	trans  transport.Transport
	broker broker.Broker
	log    logger.Logger
	disc   discovery.Discovery
}

type CallError struct {
	msg string
}

func (c *CallError) Error() string {
	return c.msg
}

func NewClient(opts *options.Options) Client {
	return &client{
		codec:  opts.Codec,
		trans:  opts.Transport,
		broker: opts.Broker,
		log:    opts.Logger,
	}
}

func (c *client) Call(ctx context.Context, service string, path string, reqval interface{}, respval interface{}, opts ...CallOption) error {
	var callopts CallOptions

	for _, o := range opts {
		o(&callopts)
	}

	if callopts.Context == nil {
		callopts.Context = context.Background()
	}

	s, err := c.trans.Dial(ctx, service)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer s.Close()

	reqid := "TODO: Generate ID"
	req := transport.NewMessage()
	req.Headers[transport.HeaderRequestID] = reqid
	req.Headers[transport.HeaderPath] = path

	req.Data, err = c.codec.Marshal(reqval)
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

	if err := c.codec.Unmarshal(respmsg.Data, respval); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func (c *client) Subscribe(topic string, callback interface{}) {
	if c.broker == nil {
		panic("no broker has been declared")
	}

	cb, err := c.createCallback(callback)
	if err != nil {
		panic(err.Error())
	}

	if err := c.broker.Subscribe(context.Background(), topic, cb); err != nil {
		c.log.Errorf("failed to subscribe to %s: %s", topic, err)
	}
}

func (c *client) createCallback(intf interface{}) (func(*broker.Message), error) {
	if fn, ok := intf.(func(*broker.Message)); ok {
		return fn, nil
	}

	val := reflect.ValueOf(intf)
	if val.Kind() != reflect.Func {
		return nil, ErrMustBeFunc
	}

	typ := val.Type()

	if typ.NumIn() != 1 {
		return nil, ErrInvalidInput
	}

	// hasctx := typ.NumIn() == 2

	// if hasctx && typ.In(0) != reflect.TypeOf(context.Background()) {
	// 	return nil, errors.New("the first input must be of type context.Context")
	// }

	if typ.In(0).Kind() != reflect.Ptr {
		return nil, ErrInputPointer
	}
	datatyp := typ.In(0).Elem()

	return func(msg *broker.Message) {
		data := reflect.New(datatyp)

		err := c.codec.Unmarshal(msg.Data, data.Interface())
		if err != nil {
			c.log.Errorf("failed to unmarshal event data: %s", err)
			return
		}

		val.Call([]reflect.Value{data})
	}, nil
}
