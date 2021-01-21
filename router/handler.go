package router

import (
	"fmt"
	"reflect"
)

type handler struct {
	Name      string
	Instance  interface{}
	Endpoints map[string]*endpoint
}

type HandlerError struct {
	endpoint *endpoint
	handler  *handler
	err      error
}

func (e *HandlerError) Error() string {
	return fmt.Sprintf("handler %s.%s: %s", e.handler.Name, e.endpoint.Name, e.err)
}

func (e *HandlerError) Unwrap() error {
	return e.err
}

func newHandler(h interface{}, methods map[string]bool) *handler {
	handtype := reflect.TypeOf(h)
	handval := reflect.ValueOf(h)
	name := reflect.Indirect(handval).Type().Name()

	endpoints := make(map[string]*endpoint)

	for i := 0; i < handtype.NumMethod(); i++ {
		m := handtype.Method(i)

		if methods != nil && !methods[m.Name] {
			continue
		}

		e := getEndpoint(m)
		if e == nil {
			continue
		}

		endpoints[e.Name] = e
	}

	return &handler{
		Name:      name,
		Instance:  h,
		Endpoints: endpoints,
	}
}

func getEndpoint(m reflect.Method) *endpoint {
	// Exported methods have an empty PkgPath
	if m.PkgPath != "" {
		return nil
	}

	// Functions must have inputs like (c context.Context, req *data.Request, resp *data.Response), plus one input for the receiver
	if m.Type.NumIn() != 4 {
		return nil
	}

	// Functions must return an error value
	if m.Type.NumOut() != 1 {
		return nil
	}

	return &endpoint{
		Name:        m.Name,
		HandlerFunc: m.Func,
		In:          m.Type.In(2),
		Out:         m.Type.In(3),
	}
}
