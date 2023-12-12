// File: "xlog.go"

package xlog

import (
	"log"
	"log/slog" // go>=1.21
	"os"
	"strings"
	//"golang.org/x/exp/slog" // depricated for go>=1.21
)

const ERR_KEY = "err"

// Xlog wrapper
type Xlog struct{ *slog.Logger }

// Create Xlog based on default slog.Logger
func Default() Xlog { return Xlog{slogDefault()} }

// Return current Xlog
func Current() Xlog { return currentXlog }

// Return current *slog.Logger
func Slog() *slog.Logger { return Current().Slog() }

// Create Xlog based on *slog.Logger (*slog.Logger -> Xlog)
func X(logger *slog.Logger) Xlog {
	if logger == nil {
		return Default()
	}
	return Xlog{logger}
}

// Create new custom Xlog
func New(conf Conf) Xlog {
	return X(NewSlog(conf))
}

// Extract *slog.Logger from Xlog (Xlog -> *slog.Logger)
func (x Xlog) Slog() *slog.Logger {
	if x.Logger == nil {
		x.Logger = slog.Default()
	}
	return x.Logger
}

// Set Xlog logger as default xlog logger
func (x Xlog) SetDefault() {
	defaultLock.Lock()
	defer defaultLock.Unlock()
	currentXlog = x
}

// Set Xlog logger as default xlog/log/slog loggers
func (x Xlog) SetDefaultLogs() {
	defaultLock.Lock()
	defer defaultLock.Unlock()
	slog.SetDefault(x.Slog())
	currentXlog = x
}

// Use xlog as io.Writer: log to level Info
func (x Xlog) Write(p []byte) (n int, err error) {
	msg := strings.TrimRight(string(p), "\r\n")
	logs(x.Logger, LevelInfo, msg)
	//x.Logger.Info(msg)
	return len(p), nil
}

// Return standart logger with prefix
func (x Xlog) NewLog(prefix string) *log.Logger {
	return log.New(x, prefix, 0)
}

// Log logs at given level
func (x Xlog) Log(level Level, msg string, args ...any) {
	logs(x.Logger, level, msg, args)
}

// Log logs at given level with default Xlog
func Log(level Level, msg string, args ...any) {
	logs(currentXlog.Logger, level, msg, args)
}

// Flood logs at LevelFlood
func (x Xlog) Flood(msg string, args ...any) {
	logs(x.Logger, LevelFlood, msg, args...)
}

// Flood logs at LevelFlood with default Xlog
func Flood(msg string, args ...any) {
	logs(currentXlog.Logger, LevelFlood, msg, args...)
}

// Trace logs at LevelTrace
func (x Xlog) Trace(msg string, args ...any) {
	logs(x.Logger, LevelTrace, msg, args...)
}

// Trace logs at LevelTrace with default Xlog
func Trace(msg string, args ...any) {
	logs(currentXlog.Logger, LevelTrace, msg, args...)
}

// Debug logs at LevelDebug
func (x Xlog) Debug(msg string, args ...any) {
	logs(x.Logger, LevelDebug, msg, args...)
}

// Debug logs at LevelDebug with default Xlog
func Debug(msg string, args ...any) {
	logs(currentXlog.Logger, LevelDebug, msg, args...)
}

// Info logs at LevelInfo
func (x Xlog) Info(msg string, args ...any) {
	logs(x.Logger, LevelInfo, msg, args...)
}

// Info logs at LevelInfo with default Xlog
func Info(msg string, args ...any) {
	logs(currentXlog.Logger, LevelInfo, msg, args...)
}

// Notice logs at LevelNotice
func (x Xlog) Notice(msg string, args ...any) {
	logs(x.Logger, LevelNotice, msg, args...)
}

// Notice logs at LevelNotice with default Xlog
func Notice(msg string, args ...any) {
	logs(currentXlog.Logger, LevelNotice, msg, args...)
}

// Warn logs at LevelWarn
func (x Xlog) Warn(msg string, args ...any) {
	logs(x.Logger, LevelWarn, msg, args...)
}

// Warn logs at LevelWarn with default Xlog
func Warn(msg string, args ...any) {
	logs(currentXlog.Logger, LevelWarn, msg, args...)
}

// Error logs at LevelError
func (x Xlog) Error(msg string, args ...any) {
	logs(x.Logger, LevelError, msg, args...)
}

// Error logs at LevelError with default Xlog
func Error(msg string, args ...any) {
	logs(currentXlog.Logger, LevelError, msg, args...)
}

// Fatal logs at LevelFatal and os.Exit(1)
func (x Xlog) Fatal(msg string, args ...any) {
	logs(x.Logger, LevelFatal, msg, args...)
	os.Exit(1)
}

// Fatal logs at LevelFatal with default Xlog and os.Exit(1)
func Fatal(msg string, args ...any) {
	logs(currentXlog.Logger, LevelFatal, msg, args...)
	os.Exit(1)
}

// Panic logs at LevelPanic and panic
func (x Xlog) Panic(msg string) {
	logs(x.Logger, LevelPanic, msg)
	panic(msg)
}

// Panic logs at LevelPanic with default Xlog and panic
func Panic(msg string) {
	logs(currentXlog.Logger, LevelPanic, msg)
	panic(msg)
}

// Logf logs at given level as standart logger
func (x Xlog) Logf(level Level, format string, args ...any) {
	logf(x.Logger, level, format, args...)
}

// Logf logs at given level as standart logger with default Xlog
func Logf(level Level, format string, args ...any) {
	logf(currentXlog.Logger, level, format, args...)
}

// Floodf logs at LevelFlood as standart logger
func (x Xlog) Floodf(format string, args ...any) {
	logf(x.Logger, LevelFlood, format, args...)
}

// Floodf logs at LevelFlood as standart logger with default Xlog
func Floodf(format string, args ...any) {
	logf(currentXlog.Logger, LevelFlood, format, args...)
}

// Tracef logs at LevelTrace as standart logger
func (x Xlog) Tracef(format string, args ...any) {
	logf(x.Logger, LevelTrace, format, args...)
}

// Tracef logs at LevelTrace as standart logger with default Xlog
func Tracef(format string, args ...any) {
	logf(currentXlog.Logger, LevelTrace, format, args...)
}

// Debugf logs at LevelDebug as standart logger
func (x Xlog) Debugf(format string, args ...any) {
	logf(x.Logger, LevelDebug, format, args...)
}

// Debugf logs at LevelDebug as standart logger with default Xlog
func Debugf(format string, args ...any) {
	logf(currentXlog.Logger, LevelDebug, format, args...)
}

// Infof logs at LevelInfo as standart logger
func (x Xlog) Infof(format string, args ...any) {
	logf(x.Logger, LevelInfo, format, args...)
}

// Infof logs at LevelInfo as standart logger with default Xlog
func Infof(format string, args ...any) {
	logf(currentXlog.Logger, LevelInfo, format, args...)
}

// Noticef logs at LevelNotice as standart logger
func (x Xlog) Noticef(format string, args ...any) {
	logf(x.Logger, LevelNotice, format, args...)
}

// Noticef logs at LevelNotice as standart logger with default Xlog
func Noticef(format string, args ...any) {
	logf(currentXlog.Logger, LevelNotice, format, args...)
}

// Warnf logs at LevelWarn as standart logger
func (x Xlog) Warnf(format string, args ...any) {
	logf(x.Logger, LevelWarn, format, args...)
}

// Warnf logs at LevelWarn as standart logger with default Xlog
func Warnf(format string, args ...any) {
	logf(currentXlog.Logger, LevelWarn, format, args...)
}

// Errorf logs at LevelError as standart logger
func (x Xlog) Errorf(format string, args ...any) {
	logf(x.Logger, LevelError, format, args...)
}

// Errorf logs at LevelError as standart logger with default Xlog
func Errorf(format string, args ...any) {
	logf(currentXlog.Logger, LevelError, format, args...)
}

// Fatalf logs at LevelFatal as standart logger and os.Exit(1)
func (x Xlog) Fatalf(format string, args ...any) {
	logf(x.Logger, LevelFatal, format, args...)
	os.Exit(1)
}

// Fatalf logs at LevelFatal as standart logger with default Xlog
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

// EOF: "xlog.go"
