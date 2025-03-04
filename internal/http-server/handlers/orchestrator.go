package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/Onnywrite/lms-golang-24/internal/model"
	"github.com/Onnywrite/lms-golang-24/internal/service/orchestrator"
	"github.com/Onnywrite/lms-golang-24/pkg/erix"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Orchestrator struct {
	Expressioner Expressioner
	Tasker       Tasker
}

type Expressioner interface {
	Schedule(context.Context, string) (orchestrator.ScheduleOutput, error)
	Expressions(context.Context) (orchestrator.ExpressionsOutput, error)
	ExpressionById(context.Context,
		orchestrator.ExpressionByIdInput,
	) (orchestrator.ExpressionByIdOutput, error)
}

func (h Orchestrator) Calculate() echo.HandlerFunc {
	type Request struct {
		Expression string `json:"expression"`
	}

	type Response struct {
		Id uuid.UUID `json:"id"`
	}

	return func(c echo.Context) error {
		var req Request
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		result, err := h.Expressioner.Schedule(c.Request().Context(), req.Expression)
		if err != nil {
			return echo.NewHTTPError(erix.HttpCode(err), erix.LastReason(err))
		}

		return c.JSON(http.StatusOK, Response{Id: result.Expression.Id})
	}
}

func (h Orchestrator) Expressions() echo.HandlerFunc {
	type Expression struct {
		Id     uuid.UUID         `json:"id"`
		Status model.Status      `json:"status"`
		Result sql.Null[float64] `json:"result,omitempty"`
		Error  sql.Null[string]  `json:"error,omitempty"`
	}

	type Response struct {
		Expressions []Expression `json:"expressions"`
	}

	return func(c echo.Context) error {
		result, err := h.Expressioner.Expressions(c.Request().Context())
		if err != nil {
			return echo.NewHTTPError(erix.HttpCode(err), erix.LastReason(err))
		}

		expressions := make([]Expression, len(result.Expressions))
		for i, expr := range result.Expressions {
			expressions[i] = Expression(expr)
		}

		return c.JSON(http.StatusOK, Response{Expressions: expressions})
	}
}

func (h Orchestrator) ExpressionById() echo.HandlerFunc {
	type Response struct {
		Id     uuid.UUID         `json:"id"`
		Status string            `json:"status"`
		Result sql.Null[float64] `json:"result"`
		Error  sql.Null[string]  `json:"error,omitempty"`
	}

	return func(c echo.Context) error {
		expressionIdStr := c.Param("id")

		expressionId, err := uuid.Parse(expressionIdStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		result, err := h.Expressioner.ExpressionById(c.Request().Context(),
			orchestrator.ExpressionByIdInput{
				ExpressionId: expressionId,
			})
		if err != nil {
			return echo.NewHTTPError(erix.HttpCode(err), erix.LastReason(err))
		}

		return c.JSON(http.StatusOK, Response(result.Expression))
	}
}

type Tasker interface {
	Task(context.Context) (orchestrator.TaskOutput, error)
	AcceptTaskResult(context.Context, orchestrator.AcceptTaskResultInput) error
}

func (h Orchestrator) GetNewTask() echo.HandlerFunc {
	type Response struct {
		Id        uuid.UUID `json:"id"`
		Args      []float64 `json:"args"`
		Operation string    `json:"operation"`
	}

	return func(c echo.Context) error {
		result, err := h.Tasker.Task(c.Request().Context())
		if err != nil {
			return echo.NewHTTPError(erix.HttpCode(err), erix.LastReason(err))
		}

		return c.JSON(http.StatusOK, Response(result.Task))
	}
}

func (h Orchestrator) AcceptTaskResult() echo.HandlerFunc {
	type Request struct {
		Result float64 `json:"result"`
		Error  *string `json:"error"`
	}

	return func(c echo.Context) error {
		taskIdStr := c.Param("id")

		taskId, err := uuid.Parse(taskIdStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		var req Request
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		err = h.Tasker.AcceptTaskResult(c.Request().Context(),
			orchestrator.AcceptTaskResultInput{
				CompletedTask: orchestrator.CompletedTask{
					Id:     taskId,
					Result: req.Result,
				},
			})
		if err != nil {
			return echo.NewHTTPError(erix.HttpCode(err), erix.LastReason(err))
		}

		return c.NoContent(http.StatusOK)
	}
}
