package clause

import "github.com/doug-martin/goqu/v9/exp"

type With struct {
	Name       string
	Expression exp.Expression
}
