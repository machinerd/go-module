package cmd

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

// DB type 관련 선언
type ID interface {
	string | int
}

type CreateInput[T ID, P CreateData[T]] struct {
	Tx    *sqlx.Tx
	Table string
	Data  P
	Func  *func(tx *sqlx.Tx, input P) error
}

type CreateMultipleInput[T ID, P CreateData[T]] struct {
	Tx    *sqlx.Tx
	Table string
	Data  []P
	Func  *func(tx *sqlx.Tx, input P) error
}

type CreateData[T ID] interface {
	SetID(id T)
}

func CreateMultiple[T ID, P CreateData[T]](input CreateMultipleInput[T, P]) ([]T, error) {
	var ids []T
	for _, v := range input.Data {
		ds := goqu.Dialect("postgres").Insert(input.Table).Rows(v).Returning("id")
		query, params, err := ds.ToSQL()
		if err != nil {
			return ids, err
		}
		r, err := input.Tx.Queryx(query, params...)
		if err != nil {
			return ids, err
		}
		var id T
		for r.Next() {
			if err := r.Scan(&id); err != nil {
				return ids, err
			} else {
				v.SetID(id)
				ids = append(ids, id)
			}
		}
		if input.Func != nil {
			err = (*input.Func)(input.Tx, v)
			if err != nil {
				return nil, err
			}
		}
	}
	return ids, nil
}

func Create[T ID, P CreateData[T]](input CreateInput[T, P]) (T, error) {
	var id T
	ds := goqu.Dialect("postgres").Insert(input.Table).
		Rows(input.Data).
		Returning("id")
	query, params, err := ds.ToSQL()
	if err != nil {
		return id, err
	}
	r, err := input.Tx.Queryx(query, params...)
	if err != nil {
		return id, err
	}
	for r.Next() {
		if err := r.Scan(&id); err != nil {
			return id, err
		} else {
			input.Data.SetID(id)
		}
	}
	if input.Func != nil {
		err = (*input.Func)(input.Tx, input.Data)
		if err != nil {
			return id, err
		}
	}
	return id, nil
}
