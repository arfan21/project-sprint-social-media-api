package logger

import (
	"context"
	"os"
	"sync"

	"github.com/agoda-com/opentelemetry-go/otelzerolog"
	otel "github.com/agoda-com/opentelemetry-logs-go"
	"github.com/arfan21/project-sprint-social-media-api/config"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2/middleware/requestid"
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
	if req_id, ok := ctx.Value(requestid.ConfigDefault.ContextKey).(string); ok {
		newlogger = newlogger.With().Str(fiberzerolog.FieldRequestID, req_id).Logger()
	}

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
