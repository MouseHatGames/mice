package transport

import (
	"context"
)

type contextKey int

const (
	keyRequest contextKey = iota
)

func GetContextRequest(ctx context.Context) (*Message, bool) {
	val := ctx.Value(keyRequest)
	if val == nil {
		return nil, false
	}

	return val.(*Message), true
}

func ContextWithRequest(ctx context.Context, req *Message) context.Context {
	return context.WithValue(ctx, keyRequest, req)
}
