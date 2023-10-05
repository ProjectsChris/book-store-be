package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
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

func newMeterProvider(ctx context.Context, endPoint string, res *resource.Resource) (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(endPoint),
	)
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(
			metricExporter,
			metric.WithInterval(3*time.Second)),
		),
	)

	return meterProvider, nil
}

func InitMetric(ctx context.Context, endPoint string) (func(context.Context) error, error) {

	res, err := resource.New(ctx)
	if err != nil {
		return nil, err
	}

	meterProvider, err := newMeterProvider(ctx, endPoint, res)
	if err != nil {
		return nil, err
	}

	otel.SetMeterProvider(meterProvider)
	return meterProvider.Shutdown, nil
}

// Init a tracer
func InitTracer(ctx context.Context, endPoint string) (func(context.Context) error, error) {
	exporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(endPoint),
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
