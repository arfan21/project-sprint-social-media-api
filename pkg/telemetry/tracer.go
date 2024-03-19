package telemetry

import (
	"context"

	"github.com/arfan21/project-sprint-social-media-api/pkg/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func InitTracer() (func(context.Context) error, error) {
	// secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	// if len(insecure) > 0 {
	// 	secureOption = otlptracegrpc.WithInsecure()
	// }
	secureOption := otlptracegrpc.WithInsecure()

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)

	if err != nil {
		logger.Log(context.Background()).Error().Err(err).Msg("could not initialize exporter")
		return nil, err
	}

	// tracerProvider := sdktrace.NewTracerProvider(
	// 	sdktrace.WithSampler(sdktrace.AlwaysSample()),
	// 	sdktrace.WithBatcher(exporter),
	// 	sdktrace.WithResource(resources),
	// )
	// ginMode := os.Getenv("GIN_MODE")
	// if ginMode == "" || ginMode == gin.DebugMode {
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource()),
	)
	// }

	otel.SetTracerProvider(
		tracerProvider,
	)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return exporter.Shutdown, nil
}

func InitSpanWithRequestId(ctx context.Context, span trace.Span) (context.Context, trace.Span) {
	val, ok := ctx.Value(requestid.ConfigDefault.ContextKey).(string)
	if ok {
		span.SetAttributes(
			attribute.String("requestId", val),
		)

	}

	return ctx, span
}
