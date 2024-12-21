package httpserver

import (
	"net/http"

	"github.com/Onnywrite/lms-golang-24/pkg/calc"

	"github.com/labstack/echo/v4"
)

func RegisterApiV1(r *echo.Group) {
	r.POST("/calculate", CalculateHand())
	r.POST("/panic", func(_ echo.Context) error {
		panic("test panic")
	})
}

func CalculateHand() echo.HandlerFunc {
	type Request struct {
		Expression string `json:"expression"`
	}

	type Response struct {
		Result float64 `json:"result"`
	}

	return func(c echo.Context) error {
		var req Request
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		result, err := calc.Calc(req.Expression)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		return c.JSON(http.StatusOK, Response{Result: result})
	}
}
