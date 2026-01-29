package cmd

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

type UpdateInput[T ID, P UpdateData[T]] struct {
	Tx    *sqlx.Tx
	Table string
	Data  P
	Func  *func(tx *sqlx.Tx, input P) error
}

type UpdateMultipleInput[T ID, P UpdateData[T]] struct {
	Tx    *sqlx.Tx
	Table string
	Data  []P
	Func  *func(tx *sqlx.Tx, input P) error
}

type UpdateData[T ID] interface {
	GetID() T
}

func UpdateMultiple[T ID, P UpdateData[T]](input UpdateMultipleInput[T, P]) error {
	for _, v := range input.Data {
		ds := goqu.Dialect("postgres").
			Update(input.Table).
			Set(v).
			Where(goqu.C("id").Eq(v.GetID()))
		query, params, err := ds.ToSQL()
		if err != nil {
			return err
		}
		_, err = input.Tx.Exec(query, params...)
		if err != nil {
			return err
		}
		if input.Func != nil {
			err = (*input.Func)(input.Tx, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Update[T ID, P UpdateData[T]](input UpdateInput[T, P]) error {
	ds := goqu.Dialect("postgres").
		Update(input.Table).
		Set(input.Data).
		Where(goqu.C("id").Eq(input.Data.GetID()))

	query, params, err := ds.ToSQL()
	if err != nil {
		return err
	}
	_, err = input.Tx.Exec(query, params...)
	if err != nil {
		return err
	}
	if input.Func != nil {
		err = (*input.Func)(input.Tx, input.Data)
		if err != nil {
			return err
		}
	}
	return nil
}
