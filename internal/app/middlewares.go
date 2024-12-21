package app

import (
	"github.com/Onnywrite/lms-golang-24/pkg/logger"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func contextWithLogger(log zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			traceId := c.Request().Header.Get("X-Trace-Id")
			if traceId == "" {
				traceId = uuid.NewString()
			}

			c = logger.EchoWithLogger(c, log.With().Str("trace_id", traceId).Logger())

			return next(c)
		}
	}
}
