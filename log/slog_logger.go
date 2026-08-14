package log

import (
	"log/slog"
	"os"
)

type slogLogger struct {
	l        *slog.Logger
	levelVar *slog.LevelVar
}

func newSlogLogger() *slogLogger {
	levelVar := &slog.LevelVar{} // zero value is slog.LevelInfo
	return &slogLogger{
		l:        slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: levelVar})),
		levelVar: levelVar,
	}
}

func (s *slogLogger) Debug(msg string, args ...any) { s.l.Debug(msg, args...) }
func (s *slogLogger) Info(msg string, args ...any)  { s.l.Info(msg, args...) }
func (s *slogLogger) Warn(msg string, args ...any)  { s.l.Warn(msg, args...) }
func (s *slogLogger) Error(msg string, args ...any) { s.l.Error(msg, args...) }

func (s *slogLogger) Fatal(msg string, args ...any) {
	s.l.Error(msg, args...)
	os.Exit(1)
}

func (s *slogLogger) With(args ...any) Logger {
	return &slogLogger{l: s.l.With(args...), levelVar: s.levelVar}
}

// SetLevel adjusts the minimum severity this logger (and any logger
// derived from it via With) emits. Satisfies the unexported levelSetter
// interface used by the package-level SetLevel function.
func (s *slogLogger) SetLevel(level Level) {
	s.levelVar.Set(level.slogLevel())
}
