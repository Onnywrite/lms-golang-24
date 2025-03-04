package orchestrator

import (
	"context"

	"github.com/Onnywrite/lms-golang-24/internal/model"
	"github.com/Onnywrite/lms-golang-24/internal/storage"
)

type Service struct {
	tx         storage.Transactor
	expessions ExpressionsRepo
	tasks      TasksRepo
}

type ExpressionsRepo interface {
	SaveExpression(context.Context, model.Expression) (model.Expression, error)
}
type TasksRepo interface {
	SaveTasks(context.Context, []model.Task) error
}

type Dependencies struct {
	Transactor      storage.Transactor
	ExpressionsRepo ExpressionsRepo
	TasksRepo       TasksRepo
}

type Config struct {
	Dependencies
}

func New(deps Dependencies) *Service {
	return NewWithConfig(Config{
		Dependencies: deps,
	})
}

func NewWithConfig(c Config) *Service {
	return &Service{
		tx:         c.Transactor,
		expessions: c.ExpressionsRepo,
		tasks:      c.TasksRepo,
	}
}
