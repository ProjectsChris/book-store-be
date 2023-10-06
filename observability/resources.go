package observability

import (
	"context"

	"go.opentelemetry.io/otel/sdk/resource"
)

func NewResource(ctx context.Context) (*resource.Resource, error) {
	return resource.New(ctx)
}
