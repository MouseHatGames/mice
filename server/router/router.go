package router

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/MouseHatGames/mice/logger"
	"github.com/MouseHatGames/mice/options"
	"github.com/MouseHatGames/mice/tracing"
	"github.com/MouseHatGames/mice/transport"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var ErrMalformedPath = errors.New("malformed request path")
var ErrEndpointNotFound = errors.New("endpoint not found")

type Router interface {
	AddHandler(h interface{}, name string, methods []string)
	Handle(path string, req *transport.Message) ([]byte, error)
}

type router struct {
	handlers map[string]*handler
	log      logger.Logger
	opts     *options.Options
}

func NewRouter(opts *options.Options) Router {
	return &router{
		handlers: make(map[string]*handler),
		log:      opts.Logger.GetLogger("router"),
		opts:     opts,
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

func (s *router) Handle(path string, req *transport.Message) ([]byte, error) {
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

	in, err := s.decode(method.In, req.Data)
	if err != nil {
		return nil, fmt.Errorf("decode request: %w", err)
	}

	respValue := reflect.New(method.Out)

	ctx := transport.ContextWithRequest(context.Background(), req)
	ctx = tracing.ExtractFromMessage(ctx, req)

	ctx, span := s.opts.Tracer.Start(ctx, path, trace.WithAttributes(
		attribute.String("peer.service", s.opts.Name),
		attribute.Int("content_length", len(req.Data)),
	))
	defer span.End()

	ret := method.HandlerFunc.Call([]reflect.Value{
		reflect.ValueOf(handler.Instance),
		reflect.ValueOf(ctx),
		in,
		respValue,
	})

	if !ret[0].IsNil() {
		err := ret[0].Interface().(error)

		span.RecordError(err)
		span.SetStatus(codes.Error, "request handler failed")

		return nil, err
	}

	outdata, err := s.opts.Codec.Marshal(respValue.Interface())
	if err != nil {
		return nil, fmt.Errorf("encode response: %w", err)
	}

	return outdata, nil
}

func (s *router) decode(t reflect.Type, d []byte) (reflect.Value, error) {
	val := reflect.New(t)
	intf := val.Interface()

	if err := s.opts.Codec.Unmarshal(d, intf); err != nil {
		return reflect.Value{}, err
	}

	return val, nil
}
