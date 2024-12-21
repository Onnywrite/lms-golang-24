package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *App) registerHealthcheck() {
	a.server.GET("/api/v1/healthz", func(c echo.Context) error {
		return c.JSONBlob(http.StatusOK, []byte(`{"message": "service healthy"}`))
	})
}
