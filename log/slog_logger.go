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
	levelVar := &slog.LevelVar{}
	levelVar.Set(initialLevelFromEnv().slogLevel())
	return &slogLogger{
		l:        slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: levelVar})),
		levelVar: levelVar,
	}
}

// initialLevelFromEnv reads the LOG_LEVEL environment variable to seed the
// starting log level. It falls back to LevelInfo if the variable is unset
// or holds an unrecognized value.
func initialLevelFromEnv() Level {
	if level, ok := ParseLevel(os.Getenv("LOG_LEVEL")); ok {
		return level
	}
	return LevelInfo
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

// Level returns the logger's current minimum severity. Satisfies the
// unexported levelGetter interface used by the package-level CurrentLevel
// function.
func (s *slogLogger) Level() Level {
	return levelFromSlog(s.levelVar.Level())
}
