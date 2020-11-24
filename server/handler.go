package server

import "reflect"

type handler struct {
	Name      string
	Instance  interface{}
	Endpoints map[string]*endpoint
}

func newHandler(h interface{}) *handler {
	handtype := reflect.TypeOf(h)
	handval := reflect.ValueOf(h)
	name := reflect.Indirect(handval).Type().Name()

	var endpoints map[string]*endpoint

	for i := 0; i < handtype.NumMethod(); i++ {
		m := handtype.Method(i)

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

	// Functions must have inputs like (c context.Context, d *data.Data), plus one input for the receiver
	if m.Type.NumIn() != 3 {
		return nil
	}

	// Functions must return either 1 or 2 values
	if m.Type.NumOut() == 0 || m.Type.NumOut() > 2 {
		return nil
	}

	return &endpoint{
		Name:        m.Name,
		HandlerFunc: m.Func,
		In:          m.Type.In(2),
		Out:         m.Type.Out(0),
	}
}
