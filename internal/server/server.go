package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arfan21/project-sprint-social-media-api/config"
	_ "github.com/arfan21/project-sprint-social-media-api/docs"
	"github.com/arfan21/project-sprint-social-media-api/pkg/exception"
	"github.com/arfan21/project-sprint-social-media-api/pkg/logger"
	"github.com/arfan21/project-sprint-social-media-api/pkg/middleware"
	"github.com/arfan21/project-sprint-social-media-api/pkg/pkgutil"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

const (
	ctxTimeout = 5
)

type Server struct {
	app *fiber.App
	db  *pgxpool.Pool
}

func New(
	db *pgxpool.Pool,
) *Server {
	app := fiber.New(fiber.Config{
		Prefork:      true,
		ErrorHandler: exception.FiberErrorHandler,
	})

	if config.Get().Otel.EnableMetrics {
		app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	}

	timeout := time.Duration(config.Get().Service.Timeout) * time.Second
	app.Use(middleware.Timeout(timeout))

	app.Use(cors.New())
	if config.Get().Otel.EnableMetrics || config.Get().Otel.EnableTracing {
		app.Use(otelfiber.Middleware())
		app.Use(middleware.TraceID())
	}
	if (!config.Get().Otel.EnableTracing && config.Get().Otel.EnableMetrics) || config.Get().Otel.EnableMetrics {
		app.Use(requestid.New())
		app.Use(middleware.RequestIdUser())
	}

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		// Logger: logger.Log(context.Background()),
		GetLogger: func(c *fiber.Ctx) zerolog.Logger {
			return *logger.Log(c.UserContext())
		},
	}))

	app.Use(recover.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	return &Server{
		app: app,
		db:  db,
	}
}

func (s *Server) Run() error {
	s.Routes()
	s.app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(pkgutil.HTTPResponse{
			Message: "Not Found",
		})
	})
	ctx := context.Background()
	go func() {
		if err := s.app.Listen(pkgutil.GetPort()); err != nil {
			logger.Log(ctx).Fatal().Err(err).Msg("failed to start server")
		}
	}()

	// go func() {
	// 	logger.Log(ctx).Info().Msgf("Starting prometheus exporter on port %s", config.Get().Otel.ExporterPrometheusPort)
	// 	http.Handle(config.Get().Otel.ExporterPrometheusPath)
	// 	if err := http.ListenAndServe(pkgutil.GetPort(config.Get().Otel.ExporterPrometheusPort), nil); err != nil {
	// 		logger.Log(ctx).Fatal().Err(err).Msg("failed to start prometheus exporter")
	// 	}
	// }()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	_, shutdown := context.WithTimeout(ctx, ctxTimeout*time.Second)
	defer shutdown()

	logger.Log(ctx).Info().Msg("shutting down server")
	return s.app.Shutdown()
}
