package log

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLevelHandlerMissingLevel(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/admin/log-level", nil)
	w := httptest.NewRecorder()

	LevelHandler()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestLevelHandlerInvalidLevel(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/admin/log-level?level=verbose", nil)
	w := httptest.NewRecorder()

	LevelHandler()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestLevelHandlerAllValidLevels(t *testing.T) {
	original := std
	t.Cleanup(func() { std = original })

	cases := []struct {
		query string
		want  Level
	}{
		{"debug", LevelDebug},
		{"info", LevelInfo},
		{"warn", LevelWarn},
		{"warning", LevelWarn},
		{"error", LevelError},
		{"DEBUG", LevelDebug}, // case-insensitive
	}

	for _, tc := range cases {
		t.Run(tc.query, func(t *testing.T) {
			levelVar := &slog.LevelVar{}
			SetLogger(&slogLogger{
				l:        slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: levelVar})),
				levelVar: levelVar,
			})

			req := httptest.NewRequest(http.MethodPost, "/admin/log-level?level="+tc.query, nil)
			w := httptest.NewRecorder()

			LevelHandler()(w, req)

			if w.Code != http.StatusNoContent {
				t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
			}
			if levelVar.Level() != tc.want.slogLevel() {
				t.Errorf("level = %v, want %v", levelVar.Level(), tc.want.slogLevel())
			}
		})
	}
}

func TestLevelHandlerEndToEndTogglesDebugOutput(t *testing.T) {
	original := std
	t.Cleanup(func() { std = original })

	var buf bytes.Buffer
	l := newTestSlogLogger(&buf)
	SetLogger(l)

	l.Debug("hidden before level change")
	if buf.Len() != 0 {
		t.Fatalf("expected no output before level change, got %q", buf.String())
	}

	req := httptest.NewRequest(http.MethodPost, "/admin/log-level?level=debug", nil)
	w := httptest.NewRecorder()
	LevelHandler()(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}

	buf.Reset()
	l.Debug("visible after level change")
	if buf.Len() == 0 {
		t.Fatalf("expected Debug output to be visible after handler set level to debug")
	}
}
