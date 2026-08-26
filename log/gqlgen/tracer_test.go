package gqlgen

import (
	"context"
	"testing"

	"github.com/graph-gophers/dataloader/v7"
)

func TestLoggingTracerLogsBatchDone(t *testing.T) {
	m := withMockLogger(t)

	tracer := LoggingTracer[string, int]{
		Name:         "TestLoader",
		GetRequestID: func(context.Context) string { return "req-1" },
	}

	ctx, finish := tracer.TraceBatch(context.Background(), []string{"a", "b"})
	finish([]*dataloader.Result[int]{
		{Data: 1},
		{Data: 2},
	})
	_ = ctx

	want := []string{"debug:dataloader batch start", "info:dataloader batch done"}
	if len(m.calls) != len(want) {
		t.Fatalf("got %d calls, want %d: %v", len(m.calls), len(want), m.calls)
	}
	for i, w := range want {
		if m.calls[i] != w {
			t.Errorf("call %d = %q, want %q", i, m.calls[i], w)
		}
	}
}

func TestLoggingTracerLogsBatchFailedWhenAnyResultErrors(t *testing.T) {
	m := withMockLogger(t)

	tracer := LoggingTracer[string, int]{
		Name:         "TestLoader",
		GetRequestID: func(context.Context) string { return "req-1" },
	}

	_, finish := tracer.TraceBatch(context.Background(), []string{"a", "b"})
	finish([]*dataloader.Result[int]{
		{Data: 1},
		{Error: context.DeadlineExceeded},
	})

	want := []string{"debug:dataloader batch start", "error:dataloader batch failed"}
	if len(m.calls) != len(want) {
		t.Fatalf("got %d calls, want %d: %v", len(m.calls), len(want), m.calls)
	}
	for i, w := range want {
		if m.calls[i] != w {
			t.Errorf("call %d = %q, want %q", i, m.calls[i], w)
		}
	}
}

func TestLoggingTracerLoadAndLoadManyAreNoops(t *testing.T) {
	withMockLogger(t)

	tracer := LoggingTracer[string, int]{
		Name:         "TestLoader",
		GetRequestID: func(context.Context) string { return "req-1" },
	}

	if _, finish := tracer.TraceLoad(context.Background(), "a"); finish == nil {
		t.Fatalf("TraceLoad() returned a nil finish func")
	} else {
		finish(nil)
	}

	if _, finish := tracer.TraceLoadMany(context.Background(), []string{"a", "b"}); finish == nil {
		t.Fatalf("TraceLoadMany() returned a nil finish func")
	} else {
		finish(nil)
	}
}
