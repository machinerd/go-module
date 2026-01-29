package clause

import "github.com/doug-martin/goqu/v9/exp"

type Join struct {
	As        string
	Table     exp.Expression
	Condition exp.JoinCondition
}

func ContainsAs(as string) func(v Join) bool {
	return func(v Join) bool {
		return v.As == as
	}
}
