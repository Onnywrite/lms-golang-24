package model

import (
	"database/sql"

	"github.com/Onnywrite/lms-golang-24/pkg/calc"
	"github.com/google/uuid"
)

type Status string

const (
	StatusPending    Status = "Pending"
	StatusEvaluating Status = "Evaluating"
	StatusError      Status = "Error"
	StatusDone       Status = "Done"
)

type Expression struct {
	Id     uuid.UUID
	Status Status
	// TODO: ProgressPercent int
	Result sql.Null[float64]
	Error  sql.Null[string]
}

type Task struct {
	TaskId         uuid.UUID
	ExpressionId   uuid.UUID
	Arg1, Arg2     sql.Null[float64]
	Arg1Id, Arg2Id uuid.NullUUID // IDs of the tasks this one depends on
	Operator       calc.Operator
	Status         Status
	Result         sql.Null[float64]
	Error          sql.Null[string]
}
