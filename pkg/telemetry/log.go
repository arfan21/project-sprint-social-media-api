package telemetry

import (
	"context"

	otel "github.com/agoda-com/opentelemetry-logs-go"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogsgrpc"
	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"github.com/arfan21/project-sprint-social-media-api/config"
	"github.com/arfan21/project-sprint-social-media-api/pkg/logger"
)

func InitLogs() (func(context.Context) error, error) {
	ctx := context.Background()

	// secureOption := otlplogsgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	// if configuration.Env().Otel.Insecure {
	// 	secureOption = otlplogsgrpc.WithInsecure()
	// }
	secureOption := otlplogsgrpc.WithInsecure()

	grpcClient := otlplogsgrpc.NewClient(
		otlplogsgrpc.WithEndpoint(config.Get().Otel.ExporterOTLPEndpoint),
		secureOption,
	)
	exporter, err := otlplogs.NewExporter(ctx,
		otlplogs.WithClient(grpcClient),
	)
	if err != nil {
		logger.Log(ctx).Error().Err(err).Msg("logs: could not initialize exporter")
		return nil, err
	}
	loggerProvider := sdk.NewLoggerProvider(
		sdk.WithBatcher(exporter),
		sdk.WithResource(newResource()),
	)

	otel.SetLoggerProvider(loggerProvider)

	return loggerProvider.Shutdown, nil
}
