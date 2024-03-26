package telemetry

import (
	"context"

	"github.com/arfan21/project-sprint-social-media-api/config"
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
			semconv.ServiceName(config.Get().Service.Name),
			semconv.ContainerName(config.Get().Service.Name),
			semconv.ServiceVersion(config.Get().Service.Version),
		),
	)
	sdkResource, _ = resource.Merge(
		resource.Default(),
		extraResources,
	)
	return sdkResource
}
