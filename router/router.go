package router

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/MouseHatGames/mice/codec"
)

var ErrMalformedPath = errors.New("malformed request path")
var ErrEndpointNotFound = errors.New("endpoint not found")

type Router interface {
	AddHandler(h interface{})
	Handle(path string, data []byte) ([]byte, error)
}

type router struct {
	handlers map[string]*handler
	codec    codec.Codec
}

func newRouter(cod codec.Codec) Router {
	return &router{
		handlers: make(map[string]*handler),
		codec:    cod,
	}
}

func (s *router) AddHandler(h interface{}) {
	hdl := newHandler(h)
	s.handlers[hdl.Name] = hdl
}

func (s *router) Handle(path string, data []byte) ([]byte, error) {
	dotidx := strings.IndexRune(path, '.')
	if dotidx == -1 {
		return nil, ErrMalformedPath
	}

	hndname := path[:dotidx]
	metname := path[dotidx+1:]

	handler, ok := s.handlers[hndname]
	if !ok {
		return nil, ErrEndpointNotFound
	}

	method, ok := handler.Endpoints[metname]
	if !ok {
		return nil, ErrEndpointNotFound
	}

	in, err := s.decode(method.In, data)
	if err != nil {
		return nil, fmt.Errorf("decode request: %w", err)
	}

	ret := method.HandlerFunc.Call([]reflect.Value{
		reflect.ValueOf(handler.Instance),
		reflect.ValueOf(context.Background()),
		*in,
	})

	if len(ret) == 2 && !ret[1].IsNil() {
		return nil, &HandlerError{
			endpoint: method,
			handler:  handler,
			err:      ret[1].Interface().(error),
		}
	}

	outdata, err := s.codec.Marshal(ret[1].Interface())
	if err != nil {
		return nil, fmt.Errorf("encode response: %w", err)
	}

	return outdata, nil
}

func (s *router) decode(t reflect.Type, d []byte) (*reflect.Value, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	val := reflect.New(t)
	intf := val.Interface()

	if err := s.codec.Unmarshal(d, intf); err != nil {
		return nil, err
	}

	return &val, nil
}
