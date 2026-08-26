package gqlgen

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/99designs/gqlgen/graphql"
	"github.com/machinerd/go-module/log"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// RecoverFunc returns a gqlgen graphql.RecoverFunc that logs a resolver
// panic (with a stack trace and the request id getRequestID extracts from
// ctx) through go-module/log. Without this, gqlgen's DefaultRecover writes
// plain text straight to stderr, with no request id and no structured
// fields, bypassing the rest of the logging setup entirely. Wire it up
// with:
//
//	svr.SetRecoverFunc(gqlgen.RecoverFunc(myRequestIDFunc))
func RecoverFunc(getRequestID RequestIDFunc) graphql.RecoverFunc {
	return func(ctx context.Context, recovered any) error {
		requestID := getRequestID(ctx)
		log.Error("graphql operation panicked", "err", fmt.Sprint(recovered), "request_id", requestID, "stack", string(debug.Stack()))
		return gqlerror.Errorf("internal system error")
	}
}
