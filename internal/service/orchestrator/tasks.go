package orchestrator

import (
	"context"

	"github.com/google/uuid"
)

type Task struct {
	Id        uuid.UUID
	Args      []float64
	Operation string
}

type TaskOutput struct {
	Task Task
}

func (s *Service) Task(context.Context) (TaskOutput, error) {
	return TaskOutput{}, nil
}

type CompletedTask struct {
	Id     uuid.UUID
	Result float64
}

type AcceptTaskResultInput struct {
	CompletedTask CompletedTask
}

func (s *Service) AcceptTaskResult(context.Context, AcceptTaskResultInput) error {
	return nil
}
