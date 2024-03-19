package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

const (
	TraceIDHeaderKey = "X-Trace-ID"
)

func TraceID() fiber.Handler {
	return func(c *fiber.Ctx) error {

		userCtx := c.UserContext()

		// get tracer
		span := trace.SpanFromContext(userCtx)

		// set trace id
		if span.SpanContext().HasTraceID() {
			c.Set(TraceIDHeaderKey, span.SpanContext().TraceID().String())
		}

		return c.Next()
	}
}
