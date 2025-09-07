// File: "sugar.go"

package xlog

import (
	"context"
	"log/slog" // go>=1.21
	"os"
	// FIXME: "golang.org/x/exp/slog" // экспериментальный пакет для go=1.20 только
)

// Ключ для ошибки в журнале при использовании обертки Err()
const ErrKey = "err"

// LogAttr записывает сообщение с атрибутами в структурированный журнал
// с заданным уровнем журналирования и заданным контекстом
func (c *Logger) LogAttrs(
	ctx context.Context,
	level slog.Level, msg string, attrs ...slog.Attr) error {
	return logAttrs(ctx, c.Logger, level, msg, attrs...)
}

// LogAttr записывает сообщение с атрибутами в структурированный журнал
// по умолчанию с заданным уровнем журналирования и заданным контекстом
func LogAttrs(
	ctx context.Context,
	level slog.Level, msg string, attrs ...slog.Attr) error {
	return logAttrs(ctx, currentClog.Logger, level, msg, attrs...)
}

// Log записывает сообщение в структурированный журнал с заданным
// уровнем журналирования и заданным контекстом
func (c *Logger) Log(
	ctx context.Context,
	level slog.Level, msg string, args ...any) error {
	return logs(ctx, c.Logger, level, msg, args...)
}

// Log записывает сообщение в структурированный журнал по умолчанию
// с заданным уровнем журналирования и заданным контекстом
func Log(
	ctx context.Context,
	level slog.Level, msg string, args ...any) error {
	return logs(ctx, currentClog.Logger, level, msg, args...)
}

// Flood записывает сообщение в журнал (LevelFlood)
func (c *Logger) Flood(msg string, args ...any) {
	logs(context.Background(), c.Logger, LevelFlood, msg, args...)
}

// Flood записывает сообщение в журнал по умолчанию (LevelFlood)
func Flood(msg string, args ...any) {
	logs(context.Background(), currentClog.Logger, LevelFlood, msg, args...)
}

// Trace записывает сообщение в журнал (LevelTrace)
func (c *Logger) Trace(msg string, args ...any) {
	logs(context.Background(), c.Logger, LevelTrace, msg, args...)
}

// Trace записывает сообщение в журнал по умолчанию (LevelTrace)
func Trace(msg string, args ...any) {
	logs(context.Background(), currentClog.Logger, LevelTrace, msg, args...)
}

// Debug записывает сообщение в журнал (LevelDebug)
func (c *Logger) Debug(msg string, args ...any) {
	logs(context.Background(), c.Logger, LevelDebug, msg, args...)
}

// Debug записывает сообщение в журнал по умолчанию (LevelDebug)
func Debug(msg string, args ...any) {
	logs(context.Background(), currentClog.Logger, LevelDebug, msg, args...)
}

// Info записывает сообщение в журнал (LevelInfo)
func (c *Logger) Info(msg string, args ...any) {
	logs(context.Background(), c.Logger, LevelInfo, msg, args...)
}

// Info записывает сообщение в журнал по умолчанию (LevelInfo)
func Info(msg string, args ...any) {
	logs(context.Background(), currentClog.Logger, LevelInfo, msg, args...)
}

// Notice записывает сообщение в журнал (LevelNotice)
func (c *Logger) Notice(msg string, args ...any) {
	logs(context.Background(), c.Logger, LevelNotice, msg, args...)
}

// Notice записывает сообщение в журнал по умолчанию (LevelNotice)
func Notice(msg string, args ...any) {
	logs(context.Background(), currentClog.Logger, LevelNotice, msg, args...)
}

// Warn записывает сообщение в журнал (LevelWarn)
func (c *Logger) Warn(msg string, args ...any) {
	logs(context.Background(), c.Logger, LevelWarn, msg, args...)
}

// Warn записывает сообщение в журнал по умолчанию (LevelWarn)
func Warn(msg string, args ...any) {
	logs(context.Background(), currentClog.Logger, LevelWarn, msg, args...)
}

// Error записывает сообщение в журнал (LevelError)
func (c *Logger) Error(msg string, args ...any) {
	logs(context.Background(), c.Logger, LevelError, msg, args...)
}

// Error записывает сообщение в журнал по умолчанию (LevelError)
func Error(msg string, args ...any) {
	logs(context.Background(), currentClog.Logger, LevelError, msg, args...)
}

// Crit записывает сообщение в журнал (LevelCrit)
func (c *Logger) Crit(msg string, args ...any) {
	logs(context.Background(), c.Logger, LevelCrit, msg, args...)
}

// Crit записывает сообщение в журнал по умолчанию (LevelCrit)
func Crit(msg string, args ...any) {
	logs(context.Background(), currentClog.Logger, LevelCrit, msg, args...)
}

// Alert записывает сообщение в журнал (LevelAlert)
func (c *Logger) Alert(msg string, args ...any) {
	logs(context.Background(), c.Logger, LevelAlert, msg, args...)
}

// Alert записывает сообщение в журнал по умолчанию (LevelAlert)
func Alert(msg string, args ...any) {
	logs(context.Background(), currentClog.Logger, LevelAlert, msg, args...)
}

// Emerg записывает сообщение в журнал (LevelEmerg)
func (c *Logger) Emerg(msg string, args ...any) {
	logs(context.Background(), c.Logger, LevelEmerg, msg, args...)
}

// Emerg записывает сообщение в журнал по умолчанию (LevelEmerg)
func Emerg(msg string, args ...any) {
	logs(context.Background(), currentClog.Logger, LevelEmerg, msg, args...)
}

// Fatal записывает сообщение в журнал (LevelFatal)
// и завершает приложение путем вызова os.Exit(1)
func (c *Logger) Fatal(msg string, args ...any) {
	logs(context.Background(), c.Logger, LevelFatal, msg, args...)
	os.Exit(1)
}

// Fatal записывает сообщение в журнал по умолчанию (LevelFatal)
// и завершает приложение путем вызова os.Exit(1)
func Fatal(msg string, args ...any) {
	logs(context.Background(), currentClog.Logger, LevelFatal, msg, args...)
	os.Exit(1)
}

// Fatal записывает сообщение в журнал (LevelPanic)
// и завершает приложение путем вызова panic()
func (c *Logger) Panic(msg string) {
	logs(context.Background(), c.Logger, LevelPanic, msg)
	panic(msg)
}

// Panic записывает сообщение в журнал по умолчанию (LevelPanic)
// и завершает приложение путем вызова panic()
func Panic(msg string) {
	logs(context.Background(), currentClog.Logger, LevelPanic, msg)
	panic(msg)
}

// Log записывает сообщение в традиционный журнал с заданным
// уровнем журналирования
func (c *Logger) Logf(level slog.Level, format string, args ...any) {
	logf(context.Background(), c.Logger, level, format, args...)
}

// Log записывает сообщение в традиционный журнал по умолчанию с заданным
// уровнем журналирования
func Logf(level slog.Level, format string, args ...any) {
	logf(context.Background(), currentClog.Logger, level, format, args...)
}

// Floodf записывает сообщение в традиционный журнал (LevelFlood)
func (c *Logger) Floodf(format string, args ...any) {
	logf(context.Background(), c.Logger, LevelFlood, format, args...)
}

// Floodf записывает сообщение в традиционный журнал по умолчанию (LevelFlood)
func Floodf(format string, args ...any) {
	logf(context.Background(), currentClog.Logger, LevelFlood, format, args...)
}

// Tracef записывает сообщение в традиционный журнал (LevelTrace)
func (c *Logger) Tracef(format string, args ...any) {
	logf(context.Background(), c.Logger, LevelTrace, format, args...)
}

// Tracef записывает сообщение в традиционный журнал по умолчанию (LevelTrace)
func Tracef(format string, args ...any) {
	logf(context.Background(), currentClog.Logger, LevelTrace, format, args...)
}

// Debugf записывает сообщение в традиционный журнал (LevelDebug)
func (c *Logger) Debugf(format string, args ...any) {
	logf(context.Background(), c.Logger, LevelDebug, format, args...)
}

// Debugf записывает сообщение в традиционный журнал по умолчанию (LevelDebug)
func Debugf(format string, args ...any) {
	logf(context.Background(), currentClog.Logger, LevelDebug, format, args...)
}

// Infof записывает сообщение в традиционный журнал (LevelInfo)
func (c *Logger) Infof(format string, args ...any) {
	logf(context.Background(), c.Logger, LevelInfo, format, args...)
}

// Infof записывает сообщение в традиционный журнал по умолчанию (LevelInfo)
func Infof(format string, args ...any) {
	logf(context.Background(), currentClog.Logger, LevelInfo, format, args...)
}

// Noticef записывает сообщение в традиционный журнал (LevelNotice)
func (c *Logger) Noticef(format string, args ...any) {
	logf(context.Background(), c.Logger, LevelNotice, format, args...)
}

// Noticef записывает сообщение в традиционный журнал по умолчанию (LevelNotice)
func Noticef(format string, args ...any) {
	logf(context.Background(), currentClog.Logger, LevelNotice, format, args...)
}

// Warnf записывает сообщение в традиционный журнал (LevelWarn)
func (c *Logger) Warnf(format string, args ...any) {
	logf(context.Background(), c.Logger, LevelWarn, format, args...)
}

// Warnf записывает сообщение в традиционный журнал по умолчанию (LevelWarn)
func Warnf(format string, args ...any) {
	logf(context.Background(), currentClog.Logger, LevelWarn, format, args...)
}

// Errorf записывает сообщение в традиционный журнал (LevelError)
func (c *Logger) Errorf(format string, args ...any) {
	logf(context.Background(), c.Logger, LevelError, format, args...)
}

// Errorf записывает сообщение в традиционный журнал по умолчанию (LevelError)
func Errorf(format string, args ...any) {
	logf(context.Background(), currentClog.Logger, LevelError, format, args...)
}

// Critf записывает сообщение в традиционный журнал (LevelCrit)
func (c *Logger) Critf(format string, args ...any) {
	logf(context.Background(), c.Logger, LevelCrit, format, args...)
}

// Critf записывает сообщение в традиционный журнал по умолчанию (LevelCrit)
func Critf(format string, args ...any) {
	logf(context.Background(), currentClog.Logger, LevelCrit, format, args...)
}

// Alertf записывает сообщение в традиционный журнал (LevelAlert)
func (c *Logger) Alertf(format string, args ...any) {
	logf(context.Background(), c.Logger, LevelAlert, format, args...)
}

// Alertf записывает сообщение в традиционный журнал по умолчанию (LevelAlert)
func Alertf(format string, args ...any) {
	logf(context.Background(), currentClog.Logger, LevelAlert, format, args...)
}

// Emergf записывает сообщение в традиционный журнал (LevelEmerg)
func (c *Logger) Emergf(format string, args ...any) {
	logf(context.Background(), c.Logger, LevelEmerg, format, args...)
}

// Emergf записывает сообщение в традиционный журнал по умолчанию (LevelEmerg)
func Emergf(format string, args ...any) {
	logf(context.Background(), currentClog.Logger, LevelEmerg, format, args...)
}

// Fatalf записывает сообщение в традиционный журнал (LevelFatal)
// и завершает приложение путем вызова os.Exit(1)
func (c *Logger) Fatalf(format string, args ...any) {
	logf(context.Background(), c.Logger, LevelFatal, format, args...)
	os.Exit(1)
}

// Fatalf записывает сообщение в традиционный журнал по умолчанию (LevelFatal)
// и завершает приложение путем вызова os.Exit(1)
func Fatalf(format string, args ...any) {
	logf(context.Background(), currentClog.Logger, LevelFatal, format, args...)
	os.Exit(1)
}

// Err возвращает slog.Attr с ключом "err" если err != nil
// или возвращает "пустой" атрибут, если err == nil.
// Таким образом можно логировать сообщения и исключать
// не информативные записи типа "err=nil".
func Err(err error) slog.Attr {
	if err == nil {
		return slog.Any("", nil)
	}
	return slog.Any(ErrKey, err)
}

// String возвращает slog.Attr, если  key != "" и value != ""
// или иначе возвращает "пустой" атрибут.
// Данная обёртка позволяет исключить из журнала пустые строки.
func String(key, value string) slog.Attr {
	if value == "" || key == "" {
		return slog.Any("", nil)
	}
	return slog.String(key, value)
}

// Int возвращает slog.Attr, если key != "" и value != 0
// или иначе возвращает "пустой" атрибут.
// Данная обёртка позволяет исключить из журнала нулевые значения.
func Int(key string, value int) slog.Attr {
	if value == 0 || key == "" {
		return slog.Any("", nil)
	}
	return slog.Int(key, value)
}

// EOF: "sugar.go"
