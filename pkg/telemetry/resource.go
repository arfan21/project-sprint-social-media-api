package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func newResource() *resource.Resource {
	var sdkResource *resource.Resource
	extraResources, _ := resource.New(
		context.Background(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithContainer(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ContainerName(serviceName),
			semconv.ServiceVersion(version),
		),
	)
	sdkResource, _ = resource.Merge(
		resource.Default(),
		extraResources,
	)
	return sdkResource
}
