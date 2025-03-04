package orchestrator

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Onnywrite/lms-golang-24/internal/model"
	"github.com/Onnywrite/lms-golang-24/pkg/calc"
	"github.com/Onnywrite/lms-golang-24/pkg/erix"

	"github.com/google/uuid"
)

var (
	ErrInvalidExpression = erix.NewStatus("invalid expression", erix.CodeBadRequest)
)

type ScheduleOutput struct {
	model.Expression
}

func (s *Service) Schedule(ctx context.Context, exprStr string) (ScheduleOutput, error) {
	var null ScheduleOutput

	const op = "orchestrator.Service.Schedule"

	expr, err := calc.NewCalculator(exprStr)
	if err != nil {
		return null, ErrInvalidExpression.Wrap(err)
	}

	tree, err := expr.BuildAST()
	if err != nil {
		return null, ErrInvalidExpression.Wrap(err)
	}

	expression := model.Expression{
		Id:     uuid.Must(uuid.NewV7()),
		Status: model.StatusPending,
	}

	evaluatedNodes := 0
	tasks := make([]model.Task, 0, len(tree))

	for _, node := range tree {
		if len(node.SubnodesIdxs) == 0 {
			evaluatedNodes++

			continue
		}

		if len(node.SubnodesIdxs) != 2 {
			return null, ErrInvalidExpression.WithField("subnodesCount", len(node.SubnodesIdxs))
		}

		task := model.Task{
			TaskId:       uuid.Must(uuid.NewV7()),
			ExpressionId: expression.Id,
			Operator:     node.Operator,
			Status:       model.StatusPending,
		}

		subnodeIdx1 := node.SubnodesIdxs[0]
		if tree[subnodeIdx1].Value.Valid {
			task.Arg1 = sql.Null[float64]{
				V:     tree[subnodeIdx1].Value.V,
				Valid: true,
			}
		} else {
			task.Arg1Id = uuid.NullUUID{
				UUID:  tasks[subnodeIdx1-evaluatedNodes].TaskId,
				Valid: true,
			}
		}

		subnodeIdx2 := node.SubnodesIdxs[0]
		if tree[subnodeIdx2].Value.Valid {
			task.Arg2 = sql.Null[float64]{
				V:     tree[subnodeIdx2].Value.V,
				Valid: true,
			}
		} else {
			task.Arg2Id = uuid.NullUUID{
				UUID:  tasks[subnodeIdx2-evaluatedNodes].TaskId,
				Valid: true,
			}
		}

		tasks = append(tasks, task)
	}

	s.tx.Atomically(ctx, func(ctx context.Context) error {
		expression, err = s.expessions.SaveExpression(ctx, expression)
		if err != nil {
			return fmt.Errorf("%s - save expression: %w", op, err)
		}

		err = s.tasks.SaveTasks(ctx, tasks)
		if err != nil {
			return fmt.Errorf("%s - save tasks: %w", op, err)
		}

		return nil
	})

	return ScheduleOutput{
		Expression: expression,
	}, nil
}

type Expression struct {
	Id     uuid.UUID
	Status model.Status
	Result sql.Null[float64]
	Error  sql.Null[string]
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
	Expression
}

func (s *Service) ExpressionById(context.Context, ExpressionByIdInput) (ExpressionByIdOutput, error) {
	return ExpressionByIdOutput{}, nil
}
