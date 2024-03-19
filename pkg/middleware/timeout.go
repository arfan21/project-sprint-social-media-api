package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type OptFunc func(*Config)

var defaultConfig = Config{
	duration: 30 * time.Second,
}

type Config struct {
	duration     time.Duration
	excludePaths []string
}

func WithExcludePaths(paths ...string) OptFunc {
	return func(c *Config) {
		c.excludePaths = append(c.excludePaths, paths...)
	}
}

func Timeout(duration time.Duration, opt ...OptFunc) fiber.Handler {
	config := defaultConfig

	if duration > 0 {
		config.duration = duration
	}

	for _, o := range opt {
		o(&config)
	}

	return func(c *fiber.Ctx) error {
		for _, path := range config.excludePaths {
			if c.Path() == path {
				return c.Next()
			}
		}
		ctx, cancel := context.WithTimeout(c.UserContext(), config.duration)
		defer cancel()
		c.SetUserContext(ctx)
		return c.Next()
	}
}
