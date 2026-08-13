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
