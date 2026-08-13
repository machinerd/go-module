package log

import "testing"

type mockLogger struct {
	calls   []string
	withLog *mockLogger
}

func (m *mockLogger) Debug(msg string, args ...any) { m.calls = append(m.calls, "debug:"+msg) }
func (m *mockLogger) Info(msg string, args ...any)  { m.calls = append(m.calls, "info:"+msg) }
func (m *mockLogger) Warn(msg string, args ...any)  { m.calls = append(m.calls, "warn:"+msg) }
func (m *mockLogger) Error(msg string, args ...any) { m.calls = append(m.calls, "error:"+msg) }
func (m *mockLogger) Fatal(msg string, args ...any) { m.calls = append(m.calls, "fatal:"+msg) }
func (m *mockLogger) With(args ...any) Logger {
	m.withLog = &mockLogger{}
	return m.withLog
}

func withMockLogger(t *testing.T) *mockLogger {
	t.Helper()
	original := std
	m := &mockLogger{}
	SetLogger(m)
	t.Cleanup(func() { std = original })
	return m
}

func TestPackageFunctionsDelegateToSetLogger(t *testing.T) {
	m := withMockLogger(t)

	Debug("d")
	Info("i")
	Warn("w")
	Error("e")
	Fatal("f")

	want := []string{"debug:d", "info:i", "warn:w", "error:e", "fatal:f"}
	if len(m.calls) != len(want) {
		t.Fatalf("got %d calls, want %d: %v", len(m.calls), len(want), m.calls)
	}
	for i, w := range want {
		if m.calls[i] != w {
			t.Errorf("call %d = %q, want %q", i, m.calls[i], w)
		}
	}
}

func TestWithDelegatesAndReturnsChildLogger(t *testing.T) {
	m := withMockLogger(t)

	child := With("request_id", "abc")
	if child != m.withLog {
		t.Fatalf("With() did not return the logger produced by std.With")
	}

	child.Info("hello")
	if len(m.withLog.calls) != 1 || m.withLog.calls[0] != "info:hello" {
		t.Fatalf("child logger did not receive call: %v", m.withLog.calls)
	}
}
