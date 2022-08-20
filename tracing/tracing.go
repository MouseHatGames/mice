package tracing

import (
	"context"
	"strings"

	"github.com/MouseHatGames/mice/transport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var Tracer = otel.GetTracerProvider().Tracer("mice")
var Propagator = &propagation.TraceContext{}

const headerPrefix = "tracing-"

func ExtractFromMessage(ctx context.Context, msg *transport.Message) context.Context {
	if msg == nil {
		return ctx
	}

	return Propagator.Extract(ctx, &carrierMessage{msg})
}

func InjectToMessage(ctx context.Context, msg *transport.Message) {
	if msg != nil {
		Propagator.Inject(ctx, &carrierMessage{msg})
	}
}

type carrierMessage struct {
	msg *transport.Message
}

func (c *carrierMessage) Get(key string) string {
	return c.msg.MessageHeaders[headerPrefix+key]
}

func (c *carrierMessage) Set(key string, value string) {
	c.msg.MessageHeaders[headerPrefix+key] = value
}

func (c *carrierMessage) Keys() []string {
	keys := make([]string, 0, len(c.msg.MessageHeaders))

	for h := range c.msg.MessageHeaders {
		if strings.HasPrefix(h, headerPrefix) {
			keys = append(keys, h)
		}
	}

	return keys
}
