package app

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/Onnywrite/lms-golang-24/pkg/grace"
	"github.com/Onnywrite/lms-golang-24/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

type App struct {
	log    zerolog.Logger
	server *echo.Echo
	apps   grace.ShutdownGroup
	c      Config
}

type Config struct {
	Port string
}

func New() *App {
	return NewWithConfig(Config{
		Port: os.Getenv("SERVER_PORT"),
	})
}

func NewWithConfig(c Config) *App {
	log := zerolog.New(os.Stdout).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	server := echo.New()

	server.HideBanner = true
	server.HTTPErrorHandler = echoErrorHandler()
	server.Use(middleware.Recover(), middleware.CORS(), contextWithLogger(log))

	return &App{
		log:    log,
		server: server,
		apps:   grace.NewShutdownGroup(),
		c:      c,
	}
}

func (a *App) Run(ctx context.Context) error {
	a.apps.Add(a.server)

	go func() {
		err := a.server.Start(":" + a.c.Port)
		if err != nil {
			a.log.Error().Err(err).Msg("echo server stopped")
		}
	}()

	return a.apps.WaitAndClose(ctx) //nolint:wrapcheck
}

func echoErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		log := logger.FromEcho(c)

		var he *echo.HTTPError

		ok := errors.As(err, &he)
		if !ok {
			he = &echo.HTTPError{
				Internal: nil,
				Code:     http.StatusInternalServerError,
				Message:  "Internal server error",
			}
		}

		if he.Code >= http.StatusInternalServerError {
			log.Error().Err(err).Msg("http error")
		}

		c.Response().Status = he.Code
		c.Response().WriteHeader(he.Code)

		if c.Request().Method == http.MethodHead {
			return
		}

		if he.Message != "" {
			_, _ = c.Response().Write([]byte(he.Message.(string)))
		}
	}
}
