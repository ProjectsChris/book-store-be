package observability

import (
	"context"

	"go.opentelemetry.io/otel/sdk/resource"
)

// NewResource return a new resource
func NewResource(ctx context.Context) (*resource.Resource, error) {
	return resource.New(ctx)
}
