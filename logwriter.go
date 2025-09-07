// File: "logwriter.go"

package xlog

import (
	"context"
	"io"
	"log/slog" // go>=1.21
	"strings"
	// FIXME: "golang.org/x/exp/slog" // экспериментальный пакет для go=1.20 только
)

// logWriter реализует интерфейс io.Writer для перенаправления
// потока байт в заданный журнал с заданным уровнем логирования
type logWriter struct {
	*slog.Logger
	slog.Level
	context.Context
}

// Проверить, что logWriter реализует интерфейс io.Writer
var _ io.Writer = logWriter{}

// NewLogWiter создает io.Writer на основе заданного slog логгера,
// в который может быть перенаправлен поток байт с заданным уровнем
// логирования. Записываемые в заданный io.Wtiter будут направляться
// в заданный slog.Logger в виде сообщений (атрибуты использоваться не будут).
// Функция может использоваться для построения legacy логгеров на основе
// пакета "log" с перенаправлением журнала в структурированный журнал slog.
func NewLogWriter(
	ctx context.Context,
	logger *slog.Logger, level slog.Level) io.Writer {
	return logWriter{
		Logger:  logger,
		Level:   level,
		Context: ctx,
	}
}

// Write перенаправляет последовательность байт в заданный журнал
func (w logWriter) Write(p []byte) (n int, err error) {
	msg := strings.TrimRight(string(p), "\r\n")
	logs(w.Context, w.Logger, w.Level, msg)
	return len(p), nil
}

// EOF: "logwriter.go"
