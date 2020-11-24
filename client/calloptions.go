package client

import "context"

// CallOptions represents configuration that apply to a single call
type CallOptions struct {
	Context context.Context
}

type CallOption func(*CallOptions)

// Context sets the context to be used for a call. Defaults to context.Background()
func Context(c context.Context) CallOption {
	return func(o *CallOptions) {
		o.Context = c
	}
}
