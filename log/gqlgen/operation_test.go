package gqlgen

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	stdlog "github.com/machinerd/go-module/log"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type mockLogger struct {
	calls []string
}

func (m *mockLogger) Debug(msg string, args ...any)  { m.calls = append(m.calls, "debug:"+msg) }
func (m *mockLogger) Info(msg string, args ...any)   { m.calls = append(m.calls, "info:"+msg) }
func (m *mockLogger) Warn(msg string, args ...any)   { m.calls = append(m.calls, "warn:"+msg) }
func (m *mockLogger) Error(msg string, args ...any)  { m.calls = append(m.calls, "error:"+msg) }
func (m *mockLogger) Fatal(msg string, args ...any)  { m.calls = append(m.calls, "fatal:"+msg) }
func (m *mockLogger) With(args ...any) stdlog.Logger { return m }

// withMockLogger swaps in a mock logger for the duration of the test. There
// is no exported getter on the log package, so the previous logger can't be
// restored afterwards; that's fine since each go test invocation is its own
// process and every test that logs sets its own mock first.
func withMockLogger(t *testing.T) *mockLogger {
	t.Helper()
	m := &mockLogger{}
	stdlog.SetLogger(m)
	return m
}

func TestOperationLabel(t *testing.T) {
	tests := []struct {
		name string
		oc   *graphql.OperationContext
		want string
	}{
		{
			name: "uses operation name when present",
			oc:   &graphql.OperationContext{OperationName: "GetCompany"},
			want: "GetCompany",
		},
		{
			name: "unknown when no operation and no name",
			oc:   &graphql.OperationContext{},
			want: "unknown",
		},
		{
			name: "falls back to selected field names",
			oc: &graphql.OperationContext{
				Operation: &ast.OperationDefinition{
					Operation: ast.Query,
					SelectionSet: ast.SelectionSet{
						&ast.Field{Name: "companyList"},
						&ast.Field{Name: "productList"},
					},
				},
			},
			want: "query{companyList,productList}",
		},
		{
			name: "falls back to bare operation type when selection set has no fields",
			oc: &graphql.OperationContext{
				Operation: &ast.OperationDefinition{Operation: ast.Mutation},
			},
			want: "mutation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OperationLabel(tt.oc); got != tt.want {
				t.Errorf("OperationLabel() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLogOperationsLogsStartAndDone(t *testing.T) {
	m := withMockLogger(t)

	oc := &graphql.OperationContext{OperationName: "GetCompany"}
	ctx := graphql.WithOperationContext(context.Background(), oc)

	mw := LogOperations(func(context.Context) string { return "req-1" })
	next := func(context.Context) graphql.ResponseHandler {
		return func(context.Context) *graphql.Response {
			return &graphql.Response{}
		}
	}

	resp := mw(ctx, next)(ctx)
	if resp == nil {
		t.Fatalf("expected a response")
	}

	want := []string{"debug:graphql operation start", "info:graphql operation done"}
	if len(m.calls) != len(want) {
		t.Fatalf("got %d calls, want %d: %v", len(m.calls), len(want), m.calls)
	}
	for i, w := range want {
		if m.calls[i] != w {
			t.Errorf("call %d = %q, want %q", i, m.calls[i], w)
		}
	}
}

func TestLogOperationsLogsFailureWhenResponseHasErrors(t *testing.T) {
	m := withMockLogger(t)

	oc := &graphql.OperationContext{OperationName: "GetCompany"}
	ctx := graphql.WithOperationContext(context.Background(), oc)

	mw := LogOperations(func(context.Context) string { return "req-1" })
	next := func(context.Context) graphql.ResponseHandler {
		return func(context.Context) *graphql.Response {
			return &graphql.Response{Errors: gqlerror.List{gqlerror.Errorf("boom")}}
		}
	}

	mw(ctx, next)(ctx)

	want := []string{"debug:graphql operation start", "error:graphql operation failed"}
	if len(m.calls) != len(want) {
		t.Fatalf("got %d calls, want %d: %v", len(m.calls), len(want), m.calls)
	}
	for i, w := range want {
		if m.calls[i] != w {
			t.Errorf("call %d = %q, want %q", i, m.calls[i], w)
		}
	}
}
