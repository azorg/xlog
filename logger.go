// File: "logger.go"

package xlog

import (
	"log"
	"log/slog" // go>=1.21
	"os"
	"strings"
	//"golang.org/x/exp/slog" // deprecated for go>=1.21
)

const ERR_KEY = "err"

// Logger wrapper structure
type Logger struct {
	Logger  *slog.Logger // standard slog logger
	Leveler              // set/get level interface
	Rotator              // log rotate interface
}

// Create logger based on default slog.Logger
func Default() *Logger {
	leveler := Level(DEFAULT_LEVEL)
	return &Logger{
		Logger:  slogDefault(),
		Leveler: &leveler,
		Rotator: newPipe(os.Stdout),
	}
}

// Return current Logger
func Current() *Logger { return currentXlog }

// Return current *slog.Logger
func Slog() *slog.Logger { return currentXlog.Logger }

// Create Logger based on *slog.Logger (*slog.Logger -> Logger)
func X(logger *slog.Logger) *Logger {
	if logger == nil {
		return Default()
	}
	leveler := Level(DEFAULT_LEVEL)
	return &Logger{
		Logger:  logger,
		Leveler: &leveler,
		Rotator: newPipe(os.Stdout),
	}
}

// Create logger that includes the given attributes in each output
func (x *Logger) With(args ...any) *Logger {
	return &Logger{
		Logger:  x.Logger.With(args...),
		Leveler: x.Leveler,
		Rotator: x.Rotator,
	}
}

// Create logger that includes the given attributes in each output
func (x *Logger) WithAttrs(attrs []slog.Attr) *Logger {
	return &Logger{
		Logger:  slog.New(x.Logger.Handler().WithAttrs(attrs)),
		Leveler: x.Leveler,
		Rotator: x.Rotator,
	}
}

// Create logger that starts a group
func (x *Logger) WithGroup(name string) *Logger {
	return &Logger{
		Logger:  x.Logger.WithGroup(name),
		Leveler: x.Leveler,
		Rotator: x.Rotator,
	}
}

// Extract *slog.Logger (*xlog.Logger -> *slog.Logger)
func (x *Logger) Slog() *slog.Logger {
	if x.Logger == nil {
		x.Logger = slog.Default()
	}
	return x.Logger
}

// Set logger as default logger
func (x *Logger) SetDefault() {
	defaultLock.Lock()
	defer defaultLock.Unlock()
	currentXlog = x
}

// Set logger as default xlog/log/slog loggers
func (x *Logger) SetDefaultLogs() {
	defaultLock.Lock()
	defer defaultLock.Unlock()
	slog.SetDefault(x.Slog())
	currentXlog = x
}

// Return log level as int (slog.Level)
func (x *Logger) GetLevel() slog.Level { return x.Leveler.Level() }

// Set log level as int (slog.Level)
func (x *Logger) SetLevel(level slog.Level) { x.Leveler.Update(level) }

// Return log level as string
func (x *Logger) GetLvl() string { return ParseLevel(x.Leveler.Level()) }

// Set log level as string
func (x *Logger) SetLvl(level string) { x.SetLevel(ParseLvl(level)) }

// Use logger as io.Writer: log to level Info
func (x *Logger) Write(p []byte) (n int, err error) {
	msg := strings.TrimRight(string(p), "\r\n")
	logs(x.Logger, LevelInfo, msg)
	//x.Info(msg)
	return len(p), nil
}

// Return standard logger with prefix
func (x *Logger) NewLog(prefix string) *log.Logger {
	return log.New(x, prefix, 0) // use x as io.Writer
}

// Close current log file
func Close() error {
	return currentXlog.Close()
}

// Check current log rotation possible
func Rotable() bool {
	return currentXlog.Rotable()
}

// Rotate current log file
func Rotate() error {
	return currentXlog.Rotate()
}

// Log logs at given level
func (x *Logger) Log(level slog.Level, msg string, args ...any) {
	logs(x.Logger, level, msg, args)
}

// Log logs at given level with default logger
func Log(level slog.Level, msg string, args ...any) {
	logs(currentXlog.Logger, level, msg, args)
}

// Flood logs at LevelFlood
func (x *Logger) Flood(msg string, args ...any) {
	logs(x.Logger, LevelFlood, msg, args...)
}

// Flood logs at LevelFlood with default logger
func Flood(msg string, args ...any) {
	logs(currentXlog.Logger, LevelFlood, msg, args...)
}

// Trace logs at LevelTrace
func (x *Logger) Trace(msg string, args ...any) {
	logs(x.Logger, LevelTrace, msg, args...)
}

// Trace logs at LevelTrace with default logger
func Trace(msg string, args ...any) {
	logs(currentXlog.Logger, LevelTrace, msg, args...)
}

// Debug logs at LevelDebug
func (x *Logger) Debug(msg string, args ...any) {
	logs(x.Logger, LevelDebug, msg, args...)
}

// Debug logs at LevelDebug with default logger
func Debug(msg string, args ...any) {
	logs(currentXlog.Logger, LevelDebug, msg, args...)
}

// Info logs at LevelInfo
func (x *Logger) Info(msg string, args ...any) {
	logs(x.Logger, LevelInfo, msg, args...)
}

// Info logs at LevelInfo with default logger
func Info(msg string, args ...any) {
	logs(currentXlog.Logger, LevelInfo, msg, args...)
}

// Notice logs at LevelNotice
func (x *Logger) Notice(msg string, args ...any) {
	logs(x.Logger, LevelNotice, msg, args...)
}

// Notice logs at LevelNotice with default logger
func Notice(msg string, args ...any) {
	logs(currentXlog.Logger, LevelNotice, msg, args...)
}

// Warn logs at LevelWarn
func (x *Logger) Warn(msg string, args ...any) {
	logs(x.Logger, LevelWarn, msg, args...)
}

// Warn logs at LevelWarn with default logger
func Warn(msg string, args ...any) {
	logs(currentXlog.Logger, LevelWarn, msg, args...)
}

// Error logs at LevelError
func (x *Logger) Error(msg string, args ...any) {
	logs(x.Logger, LevelError, msg, args...)
}

// Error logs at LevelError with default logger
func Error(msg string, args ...any) {
	logs(currentXlog.Logger, LevelError, msg, args...)
}

// Crit logs at LevelCritical
func (x *Logger) Crit(msg string, args ...any) {
	logs(x.Logger, LevelCritical, msg, args...)
}

// Crit logs at LevelCritical with default logger
func Crit(msg string, args ...any) {
	logs(currentXlog.Logger, LevelCritical, msg, args...)
}

// Fatal logs at LevelFatal and os.Exit(1)
func (x *Logger) Fatal(msg string, args ...any) {
	logs(x.Logger, LevelFatal, msg, args...)
	os.Exit(1)
}

// Fatal logs at LevelFatal with default logger and os.Exit(1)
func Fatal(msg string, args ...any) {
	logs(currentXlog.Logger, LevelFatal, msg, args...)
	os.Exit(1)
}

// Panic logs at LevelPanic and panic
func (x *Logger) Panic(msg string) {
	logs(x.Logger, LevelPanic, msg)
	panic(msg)
}

// Panic logs at LevelPanic with default logger and panic
func Panic(msg string) {
	logs(currentXlog.Logger, LevelPanic, msg)
	panic(msg)
}

// Logf logs at given level as standard logger
func (x *Logger) Logf(level slog.Level, format string, args ...any) {
	logf(x.Logger, level, format, args...)
}

// Logf logs at given level as standard logger with default logger
func Logf(level slog.Level, format string, args ...any) {
	logf(currentXlog.Logger, level, format, args...)
}

// Floodf logs at LevelFlood as standard logger
func (x *Logger) Floodf(format string, args ...any) {
	logf(x.Logger, LevelFlood, format, args...)
}

// Floodf logs at LevelFlood as standard logger with default logger
func Floodf(format string, args ...any) {
	logf(currentXlog.Logger, LevelFlood, format, args...)
}

// Tracef logs at LevelTrace as standard logger
func (x *Logger) Tracef(format string, args ...any) {
	logf(x.Logger, LevelTrace, format, args...)
}

// Tracef logs at LevelTrace as standard logger with default logger
func Tracef(format string, args ...any) {
	logf(currentXlog.Logger, LevelTrace, format, args...)
}

// Debugf logs at LevelDebug as standard logger
func (x *Logger) Debugf(format string, args ...any) {
	logf(x.Logger, LevelDebug, format, args...)
}

// Debugf logs at LevelDebug as standard logger with default logger
func Debugf(format string, args ...any) {
	logf(currentXlog.Logger, LevelDebug, format, args...)
}

// Infof logs at LevelInfo as standard logger
func (x *Logger) Infof(format string, args ...any) {
	logf(x.Logger, LevelInfo, format, args...)
}

// Infof logs at LevelInfo as standard logger with default logger
func Infof(format string, args ...any) {
	logf(currentXlog.Logger, LevelInfo, format, args...)
}

// Noticef logs at LevelNotice as standard logger
func (x *Logger) Noticef(format string, args ...any) {
	logf(x.Logger, LevelNotice, format, args...)
}

// Noticef logs at LevelNotice as standard logger with default logger
func Noticef(format string, args ...any) {
	logf(currentXlog.Logger, LevelNotice, format, args...)
}

// Warnf logs at LevelWarn as standard logger
func (x *Logger) Warnf(format string, args ...any) {
	logf(x.Logger, LevelWarn, format, args...)
}

// Warnf logs at LevelWarn as standard logger with default logger
func Warnf(format string, args ...any) {
	logf(currentXlog.Logger, LevelWarn, format, args...)
}

// Errorf logs at LevelError as standard logger
func (x *Logger) Errorf(format string, args ...any) {
	logf(x.Logger, LevelError, format, args...)
}

// Errorf logs at LevelError as standard logger with default logger
func Errorf(format string, args ...any) {
	logf(currentXlog.Logger, LevelError, format, args...)
}

// Critf logs at LevelCritical as standard logger
func (x *Logger) Critf(format string, args ...any) {
	logf(x.Logger, LevelCritical, format, args...)
}

// Critf logs at LevelCritical as standard logger with default logger
func Critf(format string, args ...any) {
	logf(currentXlog.Logger, LevelCritical, format, args...)
}

// Fatalf logs at LevelFatal as standard logger and os.Exit(1)
func (x *Logger) Fatalf(format string, args ...any) {
	logf(x.Logger, LevelFatal, format, args...)
	os.Exit(1)
}

// Fatalf logs at LevelFatal as standard logger with default logger
// and os.Exit(1)
func Fatalf(format string, args ...any) {
	logf(currentXlog.Logger, LevelFatal, format, args...)
	os.Exit(1)
}

// Err() returns slog.Attr with "err" key if err != nil
func Err(err error) slog.Attr {
	if err == nil {
		return slog.Any("", nil)
	}
	return slog.Any(ERR_KEY, err)
}

// String return slog.Attr if key != "" and value != ""
func String(key, value string) slog.Attr {
	if value == "" || key == "" {
		return slog.Any("", nil)
	}
	return slog.String(key, value)
}

// Integer return slog.Attr if key != "" and value != 0
func Int(key string, value int) slog.Attr {
	if value == 0 || key == "" {
		return slog.Any("", nil)
	}
	return slog.Int(key, value)
}

// EOF: "logger.go"
