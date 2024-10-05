// FIle: "interface.go"

package xlog

import (
	"log"
	//"log/slog" // go>=1.21

	"golang.org/x/exp/slog" // depricated for go>=1.21
)

// Logger interface
type Xlogger interface {
	// Create Logger that includes the given attributes in each output
	With(args ...any) *Logger

	// Create logger that includes the given attributes in each output
	WithAttrs(attrs []slog.Attr) *Logger

	// Create logger that starts a group
	WithGroup(name string) *Logger

	// Extract *slog.Logger (*xlog.Logger -> *slog.Logger)
	Slog() *slog.Logger

	// Set logger as default xlog logger
	SetDefault()

	// Set logger as default xlog/log/slog loggers
	SetDefaultLogs()

	// Return log level as int (slog.Level)
	GetLevel() slog.Level

	// Set log level as int (slog.Level)
	SetLevel(level slog.Level)

	// Return log level as string
	GetLvl() string

	// Set log level as string
	SetLvl(level string)

	// Use xlog as io.Writer: log to level Info
	Write(p []byte) (n int, err error)

	// Return standart logger with prefix
	NewLog(prefix string) *log.Logger

	// Log logs at given level
	Log(level slog.Level, msg string, args ...any)

	// Flood logs at LevelFlood
	Flood(msg string, args ...any)

	// Trace logs at LevelTrace
	Trace(msg string, args ...any)

	// Debug logs at LevelDebug
	Debug(msg string, args ...any)

	// Info logs at LevelInfo
	Info(msg string, args ...any)

	// Notice logs at LevelNotice
	Notice(msg string, args ...any)

	// Warn logs at LevelWarn
	Warn(msg string, args ...any)

	// Error logs at LevelError
	Error(msg string, args ...any)

	// Crit logs at LevelCritical
	Crit(msg string, args ...any)

	// Fatal logs at LevelFatal and os.Exit(1)
	Fatal(msg string, args ...any)

	// Panic logs at LevelPanic and panic
	Panic(msg string)

	// Logf logs at given level as standart logger
	Logf(level slog.Level, format string, args ...any)

	// Floodf logs at LevelFlood as standart logger
	Floodf(format string, args ...any)

	// Tracef logs at LevelTrace as standart logger
	Tracef(format string, args ...any)

	// Debugf logs at LevelDebug as standart logger
	Debugf(format string, args ...any)

	// Infof logs at LevelInfo as standart logger
	Infof(format string, args ...any)

	// Noticef logs at LevelNotice as standart logger
	Noticef(format string, args ...any)

	// Warnf logs at LevelWarn as standart logger
	Warnf(format string, args ...any)

	// Errorf logs at LevelError as standart logger
	Errorf(format string, args ...any)

	// Critf logs at LevelCritical as standart logger
	Critf(format string, args ...any)

	// Fatalf logs at LevelFatal as standart logger and os.Exit(1)
	Fatalf(format string, args ...any)
}

// EOF: "interface.go"
