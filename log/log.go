package log

// Logger is the only surface consumers of this module depend on for
// logging. The concrete implementation backing it can change without
// affecting call sites.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
	With(args ...any) Logger
}

var std Logger = newSlogLogger()

// SetLogger replaces the logger used by the package-level functions.
// Call it during application startup to plug in a different Logger
// implementation.
func SetLogger(l Logger) {
	std = l
}

func Debug(msg string, args ...any) { std.Debug(msg, args...) }
func Info(msg string, args ...any)  { std.Info(msg, args...) }
func Warn(msg string, args ...any)  { std.Warn(msg, args...) }
func Error(msg string, args ...any) { std.Error(msg, args...) }
func Fatal(msg string, args ...any) { std.Fatal(msg, args...) }
func With(args ...any) Logger       { return std.With(args...) }

// levelSetter is an optional capability a Logger implementation can add to
// support SetLevel. It's not part of Logger itself so implementations that
// don't need level control aren't forced to implement it.
type levelSetter interface {
	SetLevel(level Level)
}

// SetLevel adjusts the minimum severity the current logger emits. It's a
// no-op if the active Logger doesn't support level control (the default
// logger does).
func SetLevel(level Level) {
	if ls, ok := std.(levelSetter); ok {
		ls.SetLevel(level)
	}
}

// levelGetter is the read-side counterpart to levelSetter, kept separate
// for the same reason: Logger implementations that don't support level
// control aren't forced to implement it.
type levelGetter interface {
	Level() Level
}

// CurrentLevel reports the active logger's current level. ok is false if
// the active Logger doesn't support level reporting (the default logger
// does).
func CurrentLevel() (level Level, ok bool) {
	if lg, ok := std.(levelGetter); ok {
		return lg.Level(), true
	}
	return 0, false
}
