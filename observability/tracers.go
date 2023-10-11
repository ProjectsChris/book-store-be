package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// InitTracer init a tracer
func InitTracer(ctx context.Context, endPoint string) (func(context.Context) error, error) {
	exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endPoint),
	))
	if err != nil {
		return nil, err
	}

	tp, err := newTraceProvider(exporter, ctx)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)
	return exporter.Shutdown, nil
}

// newTraceProvider init a new tracer providers
func newTraceProvider(exp sdktrace.SpanExporter, ctx context.Context) (*sdktrace.TracerProvider, error) {
	r, err := NewResource(ctx)
	if err != nil {
		return nil, err
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	), nil
}
