// File: "interface.go"

package xlog

import (
	"io"
	"log"
	//"log/slog" // go>=1.21

	"golang.org/x/exp/slog" // deprecated for go>=1.21
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

	// Return standard logger with prefix
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

	// Logf logs at given level as standard logger
	Logf(level slog.Level, format string, args ...any)

	// Floodf logs at LevelFlood as standard logger
	Floodf(format string, args ...any)

	// Tracef logs at LevelTrace as standard logger
	Tracef(format string, args ...any)

	// Debugf logs at LevelDebug as standard logger
	Debugf(format string, args ...any)

	// Infof logs at LevelInfo as standard logger
	Infof(format string, args ...any)

	// Noticef logs at LevelNotice as standard logger
	Noticef(format string, args ...any)

	// Warnf logs at LevelWarn as standard logger
	Warnf(format string, args ...any)

	// Errorf logs at LevelError as standard logger
	Errorf(format string, args ...any)

	// Critf logs at LevelCritical as standard logger
	Critf(format string, args ...any)

	// Fatalf logs at LevelFatal as standard logger and os.Exit(1)
	Fatalf(format string, args ...any)

	// Check log rotation possible
	Rotable() bool

	// Close the existing log file and immediately create a new one
	Rotate() error

	// Close current log file
	Close() error

	// Create log io.Writer
	NewWriter(slog.Level) io.Writer
}

// Ensure *Logger implements Xlogger
var _ Xlogger = (*Logger)(nil)

// EOF: "interface.go"
