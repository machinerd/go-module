package audit

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DataLogRowsInput[K any] struct {
	Tx         *sqlx.Tx
	Query      func(K) (string, []interface{}, error)
	QueryInput K
}

func DataLogRows[T any, K any](ctx context.Context, input DataLogRowsInput[K]) ([]*T, error) {
	if input.Tx == nil || input.Query == nil {
		return nil, fmt.Errorf("audit.DataLogRows requires Tx and Query")
	}
	query, args, err := input.Query(input.QueryInput)
	if err != nil {
		return nil, fmt.Errorf("build audit snapshot query: %w", err)
	}
	var rows []*T
	if err := input.Tx.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("read audit snapshot: %w", err)
	}
	return rows, nil
}
