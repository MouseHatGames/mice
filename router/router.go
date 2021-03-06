package router

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/MouseHatGames/mice/codec"
	"github.com/MouseHatGames/mice/logger"
)

var ErrMalformedPath = errors.New("malformed request path")
var ErrEndpointNotFound = errors.New("endpoint not found")

type Router interface {
	AddHandler(h interface{}, name string, methods []string)
	Handle(path string, data []byte) ([]byte, error)
}

type router struct {
	handlers map[string]*handler
	codec    codec.Codec
	log      logger.Logger
}

func NewRouter(cod codec.Codec, log logger.Logger) Router {
	return &router{
		handlers: make(map[string]*handler),
		codec:    cod,
		log:      log.GetLogger("router"),
	}
}

func (s *router) AddHandler(h interface{}, name string, methods []string) {
	metmap := make(map[string]bool, len(methods))
	for _, m := range methods {
		metmap[m] = true
	}

	hdl := newHandler(h, name, metmap)
	s.handlers[hdl.Name] = hdl

	for k := range hdl.Endpoints {
		s.log.Debugf("registered endpoint %s.%s", hdl.Name, k)
	}
}

func (s *router) Handle(path string, data []byte) ([]byte, error) {
	s.log.Debugf("request to %s", path)

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

	respValue := reflect.New(method.Out)

	ret := method.HandlerFunc.Call([]reflect.Value{
		reflect.ValueOf(handler.Instance),
		reflect.ValueOf(context.Background()),
		in,
		respValue,
	})

	if !ret[0].IsNil() {
		return nil, ret[0].Interface().(error)
	}

	outdata, err := s.codec.Marshal(respValue.Interface())
	if err != nil {
		return nil, fmt.Errorf("encode response: %w", err)
	}

	return outdata, nil
}

func (s *router) decode(t reflect.Type, d []byte) (reflect.Value, error) {
	val := reflect.New(t)
	intf := val.Interface()

	if err := s.codec.Unmarshal(d, intf); err != nil {
		return reflect.Value{}, err
	}

	return val, nil
}
