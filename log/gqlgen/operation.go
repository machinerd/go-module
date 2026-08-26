// Package gqlgen provides reusable structured-logging integration points for
// gqlgen-based GraphQL subgraphs: operation start/duration/outcome logging,
// panic recovery, and dataloader batch tracing. Every log line goes through
// github.com/machinerd/go-module/log.
//
// None of this package depends on any particular web framework. Each
// subgraph supplies its own RequestIDFunc to pull a correlation id out of
// context.Context however it stores one (Gin, net/http, etc).
package gqlgen

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/machinerd/go-module/log"
	"github.com/vektah/gqlparser/v2/ast"
)

// RequestIDFunc extracts a caller-defined request/correlation id from ctx.
type RequestIDFunc func(ctx context.Context) string

// OperationLabel returns a name to log for a GraphQL operation. Most clients
// don't send an operationName for one-off/anonymous queries, so fall back to
// the top-level selected field names (e.g. "query{companyList}").
func OperationLabel(oc *graphql.OperationContext) string {
	if oc.OperationName != "" {
		return oc.OperationName
	}
	if oc.Operation == nil {
		return "unknown"
	}

	var fields []string
	for _, sel := range oc.Operation.SelectionSet {
		if f, ok := sel.(*ast.Field); ok {
			fields = append(fields, f.Name)
		}
	}
	if len(fields) == 0 {
		return string(oc.Operation.Operation)
	}
	return fmt.Sprintf("%s{%s}", oc.Operation.Operation, strings.Join(fields, ","))
}

// LogOperations returns an AroundOperations hook that logs the start, and
// then the duration and outcome, of every GraphQL operation, tagged with the
// request id getRequestID extracts from ctx. Wire it up with:
//
//	svr.AroundOperations(gqlgen.LogOperations(myRequestIDFunc))
func LogOperations(getRequestID RequestIDFunc) graphql.OperationMiddleware {
	return func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		oc := graphql.GetOperationContext(ctx)
		requestID := getRequestID(ctx)
		start := time.Now()
		operation := OperationLabel(oc)
		log.Debug("graphql operation start", "operation", operation, "request_id", requestID)

		responseHandler := next(ctx)
		return func(ctx context.Context) *graphql.Response {
			resp := responseHandler(ctx)
			duration := time.Since(start)
			if resp != nil && len(resp.Errors) > 0 {
				log.Error("graphql operation failed", "operation", operation, "request_id", requestID, "duration_ms", duration.Milliseconds(), "errors", resp.Errors)
			} else {
				log.Info("graphql operation done", "operation", operation, "request_id", requestID, "duration_ms", duration.Milliseconds())
			}
			return resp
		}
	}
}
