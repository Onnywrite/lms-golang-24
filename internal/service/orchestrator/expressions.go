package orchestrator

import (
	"context"

	"github.com/google/uuid"
)

type ScheduleOutput struct {
	ExpressionId uuid.UUID
}

func (s *Service) Schedule(context.Context, string) (ScheduleOutput, error) {
	return ScheduleOutput{}, nil
}

type Expression struct {
	Id     uuid.UUID
	Status string
	Result float64
	Error  *string
}

type ExpressionsOutput struct {
	Expressions []Expression
}

func (s *Service) Expressions(context.Context) (ExpressionsOutput, error) {
	return ExpressionsOutput{}, nil
}

type ExpressionByIdInput struct {
	ExpressionId uuid.UUID
}

type ExpressionByIdOutput struct {
	Expression Expression
}

func (s *Service) ExpressionById(context.Context, ExpressionByIdInput) (ExpressionByIdOutput, error) {
	return ExpressionByIdOutput{}, nil
}
