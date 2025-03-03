package httpserver

import (
	"github.com/Onnywrite/lms-golang-24/internal/http-server/handlers"
	"github.com/Onnywrite/lms-golang-24/internal/service/orchestrator"
	"github.com/labstack/echo/v4"
)

func RegisterApiV1(r *echo.Group, s *orchestrator.Service) {
	orch := handlers.Orchestrator{
		Expressioner: s,
		Tasker:       s,
	}

	r.POST("/calculate", orch.Calculate())
	r.GET("/expressions", orch.Expressions())
	r.GET("/expressions/:id", orch.ExpressionById())

	r.POST("/internal/tasks/get", orch.GetNewTask())
	r.POST("/internal/task/:id", orch.AcceptTaskResult())
}
