package telemetry

import (
	"context"
	"time"

	"github.com/arfan21/project-sprint-social-media-api/config"
	"github.com/arfan21/project-sprint-social-media-api/pkg/logger"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"google.golang.org/grpc/credentials"
)

func InitMetric() (func(ctx context.Context) error, error) {
	ctx := context.Background()

	options := []sdkmetric.Option{
		sdkmetric.WithResource(newResource()),
	}

	if !config.Get().Otel.OnlyPrometheusExporter {
		secureOption := otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
		if config.Get().Otel.Insecure {
			secureOption = otlpmetricgrpc.WithInsecure()
		}

		exporter, err := otlpmetricgrpc.New(
			ctx,
			secureOption,
			otlpmetricgrpc.WithEndpoint(collectorURL),
		)

		if err != nil {
			logger.Log(ctx).Error().Err(err).Msg("could not initialize exporter")
			return nil, err
		}
		reader := sdkmetric.NewPeriodicReader(exporter)
		options = append(options, sdkmetric.WithReader(reader))

		prometheusExporter, err := prometheus.New()
		if err != nil {
			logger.Log(ctx).Error().Err(err).Msg("could not initialize prometheus exporter")
			return nil, err
		}

		options = append(options, sdkmetric.WithReader(prometheusExporter))
	} else {
		prometheusExporter, err := prometheus.New()
		if err != nil {
			logger.Log(ctx).Error().Err(err).Msg("could not initialize prometheus exporter")
			return nil, err
		}

		options = append(options, sdkmetric.WithReader(prometheusExporter))
	}

	metricProvider := sdkmetric.NewMeterProvider(options...)
	otel.SetMeterProvider(metricProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	logger.Log(context.Background()).Info().Msg("Starting runtime instrumentation:")
	err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		logger.Log(context.Background()).Error().Err(err).Msg("Failed to start runtime instrumentation")
		return nil, err
	}

	return metricProvider.Shutdown, nil
}
