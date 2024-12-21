package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func EchoWithLogger(c echo.Context, log zerolog.Logger) echo.Context {
	c.SetRequest(c.Request().WithContext(WithLogger(c.Request().Context(), log)))
	return c
}

func FromEcho(c echo.Context) zerolog.Logger {
	return FromContext(c.Request().Context())
}
