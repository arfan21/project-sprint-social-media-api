package api

import (
	"context"

	"github.com/arfan21/project-sprint-social-media-api/config"
	"github.com/arfan21/project-sprint-social-media-api/internal/server"
	dbpostgres "github.com/arfan21/project-sprint-social-media-api/pkg/db/postgres"
	"github.com/arfan21/project-sprint-social-media-api/pkg/logger"
	"github.com/arfan21/project-sprint-social-media-api/pkg/telemetry"
)

func Serve() error {
	_, err := config.LoadConfig()
	if err != nil {
		return err
	}

	_, err = config.ParseConfig(config.GetViper())
	if err != nil {
		return err
	}

	if config.Get().Otel.EnableLogging {
		logShutdown, err := telemetry.InitLogs()
		if err != nil {
			return err
		}

		defer logShutdown(context.Background())

		// logger.Log(context.Background()).Info().Msg("otel logging enabled")
	} else {
		// logger.Log(context.Background()).Warn().Msg("otel logging disabled")
	}

	if config.Get().Otel.EnableTracing {
		tracerShutdown, err := telemetry.InitTracer()
		if err != nil {
			return err
		}
		defer tracerShutdown(context.Background())

		// this log called for initialize hook
		// logger.Log(context.Background()).Info().Msgf("otel tracing enabled with service name: %s", config.Get().Service.Name)
	} else {
		// logger.Log(context.Background()).Warn().Msg("otel tracing disabled")
	}

	if config.Get().Otel.EnableMetrics {
		metricShutdown, err := telemetry.InitMetric()
		if err != nil {
			return err
		}

		defer metricShutdown(context.Background())

		if config.Get().Otel.OnlyPrometheusExporter {
			logger.Log(context.Background()).Info().Msg("otel metric enabled only with prometheus")
		} else {
			logger.Log(context.Background()).Info().Msg("otel metric enabled with otlp & prometheus")
		}
	} else {
		logger.Log(context.Background()).Warn().Msg("otel metric disabled")
	}

	db, err := dbpostgres.NewPgx()
	if err != nil {
		return err
	}

	server := server.New(
		db,
	)
	return server.Run()
}
