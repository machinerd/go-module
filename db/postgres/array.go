package postgres

import (
	"fmt"
	"strings"
)

// MakeArrayString creates a PostgreSQL array string representation
func MakeArrayString[T any](values []T) string {
	var builder strings.Builder
	builder.WriteString("{")
	for i, v := range values {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprint(v))
	}
	builder.WriteString("}")
	return builder.String()
}
