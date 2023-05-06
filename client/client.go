package client

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"reflect"

	"github.com/MouseHatGames/mice/auth"
	"github.com/MouseHatGames/mice/broker"
	"github.com/MouseHatGames/mice/options"
	"github.com/MouseHatGames/mice/tracing"
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
	opts *options.Options
	port int16
}

func NewClient(opts *options.Options) Client {
	return &client{
		opts: opts,
		port: opts.RPCPort,
	}
}

func (c *client) Call(service string, path string, reqval interface{}, respval interface{}, opts ...CallOption) error {
	var callopts CallOptions

	for _, o := range opts {
		o(&callopts)
	}

	parentReq, hasParent := transport.GetContextRequest(callopts.Context)

	ctx := callopts.Context
	if ctx == nil {
		ctx = context.Background()
	}

	ctx = tracing.ExtractFromMessage(ctx, parentReq)

	if c.opts.Discovery == nil {
		panic("no discovery has been set up")
	}

	// Find service address
	host, err := c.opts.Discovery.Find(service)
	if err != nil {
		return fmt.Errorf("discover service: %w", err)
	}

	// Connect to service
	s, err := c.opts.Transport.Dial(ctx, fmt.Sprintf("%s:%d", host, c.port))
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer s.Close()

	req := transport.NewMessage()
	req.SetRandomRequestID()
	req.SetPath(path)

	if id, ok := auth.GetUserID(ctx); ok {
		req.SetUserID(id)
	}

	tracing.InjectToMessage(ctx, req)

	if hasParent {
		req.MessageHeaders[transport.HeaderParentRequestID] = parentReq.MessageHeaders[transport.HeaderRequestID]
	}

	// Encode request data
	req.Data, err = c.opts.Codec.Marshal(reqval)
	if err != nil {
		return fmt.Errorf("encode request: %w", err)
	}

	// Send request
	if err := s.Send(ctx, req); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	// Receive response
	var respmsg transport.Message
	if err := s.Receive(ctx, &respmsg); err != nil {
		return fmt.Errorf("receive message: %w", err)
	}

	// Check for server handler error
	if err, ok := respmsg.GetError(); ok {
		return err
	}

	// Decode response data
	if err := c.opts.Codec.Unmarshal(respmsg.Data, respval); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func (c *client) Subscribe(topic string, callback interface{}) {
	if c.opts.Broker == nil {
		panic("no broker has been declared")
	}

	cb, err := c.createCallback(callback)
	if err != nil {
		panic(err.Error())
	}

	if err := c.opts.Broker.Subscribe(context.Background(), topic, cb); err != nil {
		c.opts.Logger.Errorf("failed to subscribe to %s: %s", topic, err)
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

		err := c.opts.Codec.Unmarshal(msg.Data, data.Interface())
		if err != nil {
			c.opts.Logger.Errorf("failed to unmarshal event data: %s", err)
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
