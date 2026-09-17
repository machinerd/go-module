package audit

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/machinerd/go-module/db/cmd"
)

type CreateInput[T cmd.ID, P cmd.CreateData[T]] struct {
	Tx     *sqlx.Tx
	Table  string
	Change Change
	Map    func(Change) P
}

func Create[T cmd.ID, P cmd.CreateData[T]](input CreateInput[T, P]) (T, error) {
	var zero T
	if input.Tx == nil || input.Table == "" || input.Map == nil {
		return zero, fmt.Errorf("audit.Create requires Tx, Table and Map")
	}
	return cmd.Create(cmd.CreateInput[T, P]{Tx: input.Tx, Table: input.Table, Data: input.Map(input.Change)})
}

type CreateMultipleInput[T cmd.ID, P cmd.CreateData[T]] struct {
	Tx    *sqlx.Tx
	Table string
	Diff  DiffInput
	Map   func(Change) P
}

func CreateMultiple[T cmd.ID, P cmd.CreateData[T]](input CreateMultipleInput[T, P]) ([]T, error) {
	if input.Tx == nil || input.Table == "" || input.Map == nil {
		return nil, fmt.Errorf("audit.CreateMultiple requires Tx, Table and Map")
	}
	changes, err := Diff(input.Diff)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	rows := make([]P, 0, len(changes))
	for _, change := range changes {
		rows = append(rows, input.Map(change))
	}
	return cmd.CreateMultiple(cmd.CreateMultipleInput[T, P]{Tx: input.Tx, Table: input.Table, Data: rows})
}
