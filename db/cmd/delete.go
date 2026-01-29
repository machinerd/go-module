package cmd

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
)

type DeleteAllExceptIDsInput[T any, P any] struct {
	Tx              *sqlx.Tx
	Table           string
	ConditionColumn string
	ConditionID     T
	ExceptIDs       []P
}

func DeleteAllExceptIDs[T any, P any](input DeleteAllExceptIDsInput[T, P]) error {
	where := []exp.Expression{
		goqu.C(input.ConditionColumn).Eq(input.ConditionID),
	}
	ds := goqu.Dialect("postgres").Delete(input.Table)
	if len(input.ExceptIDs) > 0 {
		where = append(where, goqu.C("id").NotIn(input.ExceptIDs))
	}
	ds = ds.Where(where...)
	query, params, err := ds.ToSQL()
	if err != nil {
		return err
	}
	_, err = input.Tx.Exec(query, params...)
	return err
}
