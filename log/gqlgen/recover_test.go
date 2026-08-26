package gqlgen

import (
	"context"
	"testing"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

func TestRecoverFuncLogsAndReturnsGenericError(t *testing.T) {
	m := withMockLogger(t)

	recoverFn := RecoverFunc(func(context.Context) string { return "req-1" })
	err := recoverFn(context.Background(), "boom")

	gqlErr, ok := err.(*gqlerror.Error)
	if !ok || gqlErr.Message != "internal system error" {
		t.Fatalf("RecoverFunc() error = %v, want a generic internal error", err)
	}

	if len(m.calls) != 1 || m.calls[0] != "error:graphql operation panicked" {
		t.Fatalf("unexpected log calls: %v", m.calls)
	}
}

func TestRecoverFuncRecoversFromARealPanic(t *testing.T) {
	withMockLogger(t)

	recoverFn := RecoverFunc(func(context.Context) string { return "req-1" })

	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = recoverFn(context.Background(), r)
			}
		}()
		panic("boom")
	}()

	gqlErr, ok := err.(*gqlerror.Error)
	if !ok || gqlErr.Message != "internal system error" {
		t.Fatalf("RecoverFunc() error = %v, want a generic internal error", err)
	}
}
