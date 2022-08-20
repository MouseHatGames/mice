package tracing

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type noopTracer struct{}

func NoopTracer() trace.Tracer {
	return &noopTracer{}
}

func (noopTracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return ctx, &noopSpan{}
}

type noopSpan struct{}

func (s noopSpan) End(options ...trace.SpanEndOption) {
}

func (s noopSpan) AddEvent(name string, options ...trace.EventOption) {
}

func (s noopSpan) IsRecording() bool {
	return true
}

func (s noopSpan) RecordError(err error, options ...trace.EventOption) {
}

func (s noopSpan) SpanContext() trace.SpanContext {
	panic("not implemented")
}

func (s noopSpan) SetStatus(code codes.Code, description string) {
}

func (s noopSpan) SetName(name string) {
}

func (s noopSpan) SetAttributes(kv ...attribute.KeyValue) {
}

func (s noopSpan) TracerProvider() trace.TracerProvider {
	panic("not implemented")
}
