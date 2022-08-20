package client

import (
	"context"
	"crypto/rand"
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
	Call(service string, path string, req interface{}, resp interface{}, opts ...CallOption) error
	Subscribe(topic string, callback interface{})
}

type client struct {
	codec  codec.Codec
	trans  transport.Transport
	broker broker.Broker
	log    logger.Logger
	disc   discovery.Discovery
	port   int16
}

func NewClient(opts *options.Options) Client {
	return &client{
		codec:  opts.Codec,
		trans:  opts.Transport,
		broker: opts.Broker,
		log:    opts.Logger,
		disc:   opts.Discovery,
		port:   opts.RPCPort,
	}
}

func (c *client) Call(service string, path string, reqval interface{}, respval interface{}, opts ...CallOption) error {
	var callopts CallOptions

	for _, o := range opts {
		o(&callopts)
	}

	if callopts.Context == nil {
		callopts.Context = context.Background()
	}

	if c.disc == nil {
		panic("no discovery has been set up")
	}

	// Find service address
	host, err := c.disc.Find(service)
	if err != nil {
		return fmt.Errorf("discover service: %w", err)
	}

	// Connect to service
	s, err := c.trans.Dial(callopts.Context, fmt.Sprintf("%s:%d", host, c.port))
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer s.Close()

	req := transport.NewMessage()
	req.SetRandomRequestID()
	req.SetPath(path)

	parentReq, hasParent := transport.GetContextRequest(callopts.Context)
	if hasParent {
		req.MessageHeaders[transport.HeaderParentRequestID] = parentReq.MessageHeaders[transport.HeaderRequestID]
	}

	// Encode request data
	req.Data, err = c.codec.Marshal(reqval)
	if err != nil {
		return fmt.Errorf("encode request: %w", err)
	}

	// Send request
	if err := s.Send(callopts.Context, req); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	// Receive response
	var respmsg transport.Message
	if err := s.Receive(callopts.Context, &respmsg); err != nil {
		return fmt.Errorf("receive message: %w", err)
	}

	// Check for server handler error
	if err, ok := respmsg.GetError(); ok {
		return err
	}

	// Decode response data
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

// https://stackoverflow.com/a/25736155
func pseudo_uuid() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%x%x-%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}
