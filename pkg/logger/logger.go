package logger

import (
	"context"
	"os"
	"sync"

	"github.com/agoda-com/opentelemetry-go/otelzerolog"
	otel "github.com/agoda-com/opentelemetry-logs-go"
	"github.com/arfan21/project-sprint-social-media-api/config"
	"github.com/rs/zerolog"
)

var loggerInstance zerolog.Logger
var once sync.Once

func Log(ctx context.Context) *zerolog.Logger {
	once.Do(func() {
		multi := zerolog.MultiLevelWriter(os.Stdout)
		loggerInstance = zerolog.New(multi).With().Timestamp().Logger()

		if config.Get().Env == "dev" {
			loggerInstance = loggerInstance.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		}

		if config.Get().Otel.EnableLogging {
			loggerInstance = loggerInstance.Hook(NewHook())
		}

	})

	newlogger := loggerInstance.With().Ctx(ctx).Logger()
	return &newlogger
}

var otelZerologHook otelzerolog.Hook
var onceHook sync.Once

func NewHook() otelzerolog.Hook {
	onceHook.Do(func() {
		logger := otel.GetLoggerProvider().Logger(
			config.Get().Service.Name,
		)

		otelZerologHook = otelzerolog.Hook{Logger: logger}

	})

	return otelZerologHook
}
