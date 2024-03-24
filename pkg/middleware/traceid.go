package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
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

func RequestIdUser() fiber.Handler {
	return func(c *fiber.Ctx) error {

		val, ok := c.Locals(requestid.ConfigDefault.ContextKey).(string)
		if ok {
			userCtx := c.UserContext()
			userCtx = context.WithValue(userCtx, requestid.ConfigDefault.ContextKey, val)
			c.SetUserContext(userCtx)
		}

		return c.Next()
	}
}
