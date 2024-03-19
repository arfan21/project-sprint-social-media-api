package telemetry

import "github.com/arfan21/project-sprint-social-media-api/config"

var (
	serviceName  = config.Get().Service.Name
	collectorURL = config.Get().Otel.ExporterOTLPEndpoint
	// insecure     = os.Getenv("INSECURE_MODE")
	version = config.Get().Service.Version
)
