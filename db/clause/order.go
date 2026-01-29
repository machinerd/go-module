package clause

import (
	"slices"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

func SetOrder(sortList []string, sortMap map[string]string, defaultSorts []string) []exp.OrderedExpression {
	var orders []exp.OrderedExpression
	if len(sortList) == 0 && len(defaultSorts) > 0 {
		sortList = defaultSorts
	}
	for _, defaultSort := range defaultSorts {
		if !slices.Contains(sortList, defaultSort) {
			sortList = append(sortList, defaultSort)
		}
	}
	for _, sort := range sortList {
		key := strings.Replace(sort, "-", "", 1)
		var o exp.OrderedExpression
		field := sortMap[key]
		if field == "" {
			continue
		}
		ie := goqu.I(field)
		if strings.HasPrefix(sort, "-") {
			o = ie.Desc().NullsLast()
		} else {
			o = ie.Asc().NullsLast()
		}
		orders = append(orders, o)
	}

	return orders
}
