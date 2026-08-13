package log

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func newTestSlogLogger(buf *bytes.Buffer) *slogLogger {
	return &slogLogger{l: slog.New(slog.NewJSONHandler(buf, nil))}
}

func TestSlogLoggerLevelsProduceJSONWithFields(t *testing.T) {
	cases := []struct {
		name  string
		log   func(l *slogLogger)
		level string
	}{
		{"debug", func(l *slogLogger) { l.Debug("msg-d", "k", "v") }, "DEBUG"},
		{"info", func(l *slogLogger) { l.Info("msg-i", "k", "v") }, "INFO"},
		{"warn", func(l *slogLogger) { l.Warn("msg-w", "k", "v") }, "WARN"},
		{"error", func(l *slogLogger) { l.Error("msg-e", "k", "v") }, "ERROR"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			// DEBUG is below the default handler level, so use a handler that emits it.
			l := &slogLogger{l: slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))}
			tc.log(l)

			var got map[string]any
			if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
				t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
			}
			if got["level"] != tc.level {
				t.Errorf("level = %v, want %v", got["level"], tc.level)
			}
			if got["k"] != "v" {
				t.Errorf("field k = %v, want %v", got["k"], "v")
			}
		})
	}
}

func TestSlogLoggerWithAddsFieldsWithoutMutatingParent(t *testing.T) {
	var buf bytes.Buffer
	parent := newTestSlogLogger(&buf)

	child := parent.With("request_id", "abc-123")
	child.Info("child event")

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got["request_id"] != "abc-123" {
		t.Errorf("request_id = %v, want abc-123", got["request_id"])
	}

	buf.Reset()
	parent.Info("parent event")
	var gotParent map[string]any
	if err := json.Unmarshal(buf.Bytes(), &gotParent); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := gotParent["request_id"]; ok {
		t.Errorf("parent logger unexpectedly carries request_id from child: %v", gotParent)
	}
}

func TestSlogLoggerFatalExitsWithStatus1(t *testing.T) {
	if os.Getenv("LOG_TEST_FATAL_SUBPROCESS") == "1" {
		l := newSlogLogger()
		l.Fatal("boom", "reason", "test")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestSlogLoggerFatalExitsWithStatus1")
	cmd.Env = append(os.Environ(), "LOG_TEST_FATAL_SUBPROCESS=1")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected process to exit with an error, got %v (output: %s)", err, out.String())
	}
	if exitErr.ExitCode() != 1 {
		t.Fatalf("exit code = %d, want 1 (output: %s)", exitErr.ExitCode(), out.String())
	}
	if !strings.Contains(out.String(), "boom") {
		t.Fatalf("expected output to contain fatal message, got %q", out.String())
	}
}
