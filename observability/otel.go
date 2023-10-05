package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var Tracer *sdktrace.TracerProvider

// Init a tracer providers
func newTraceProvider(exp sdktrace.SpanExporter, ctx context.Context) (*sdktrace.TracerProvider, error) {
	r, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("Pippo", "Pluto"),
		),
	)

	if err != nil {
		return nil, err
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	), nil

}

// Init a tracer
func InitTracer() (func(context.Context) error, error) {
	ctx := context.Background()

	exporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint("172.20.10.2:4317"),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create a new tracer provider with a batch span processor and the given exporter.
	tp, err := newTraceProvider(exporter, ctx)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)
	Tracer = tp

	// Finally, set the tracer that can be used for this package.
	return exporter.Shutdown, nil
}
