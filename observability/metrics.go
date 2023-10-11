package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

// InitMetric init new metric
func InitMetric(ctx context.Context, endPoint string) (func(context.Context) error, error) {
	res, err := NewResource(ctx)
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

// newMeterProvider create a new meter provider
func newMeterProvider(ctx context.Context, endPoint string, res *resource.Resource) (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(ctx,
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
