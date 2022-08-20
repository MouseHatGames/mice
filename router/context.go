package router

import (
	"context"

	"github.com/MouseHatGames/mice/transport"
)

type contextKey int

const (
	keyRequest contextKey = iota
)

func ctxGetRequest(ctx context.Context) (*transport.Message, bool) {
	val := ctx.Value(keyRequest)
	if val == nil {
		return nil, false
	}

	return val.(*transport.Message), true
}

func ctxWithRequest(ctx context.Context, req *transport.Message) context.Context {
	return context.WithValue(ctx, keyRequest, req)
}
