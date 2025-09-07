// File: "default.go"

package xlog

import (
	"log"
	"log/slog" // go>=1.21
	"os"
	// FIXME: "golang.org/x/exp/slog" // экспериментальный пакет для go=1.20 только
)

// Глобальные переменные указывающие на логгеры по умолчанию
var (
	// Исходный (первородный) *log.Logger по умолчанию
	defaultLog *log.Logger = log.Default()

	// Исходный (первородный) *slog.Logger по умолчанию
	defaultSlog *slog.Logger = slog.Default()

	// Текущий (глобальный) логгер (*xlog.Logger)
	currentClog *Logger = Default()
)

// Default() создает логгер по умолчанию на основе
// первородного slog логгера по умолчанию (без каких-либо опций,
// с выводом на stdout), но с возможностью управления уровнем логирования.
func Default() *Logger {
	var level slog.LevelVar
	handler := defaultSlog.Handler() // slog.defaultHandler

	// Создать дополнительный хендлер-обёртку для того, чтобы управлять
	// уровнем логирования стандартного логгера slog
	handler = newLvlHandler(handler, &level)

	// Создать slog логгер с заданным хендлером
	logger := slog.New(handler)

	return &Logger{
		Logger: logger,
		Level:  &level,
		Writer: pipeWriter{os.Stdout},
	}
}

// SetDefault устанавливает заданный логгер текущим
func (c *Logger) SetDefault() {
	currentClog = c
}

// SetDefaultLogs устанавливает данный логгер текущим,
// в том числе log/slog логгером по умолчанию
func (c *Logger) SetDefaultLogs() {
	slog.SetDefault(c.Logger)
	currentClog = c
}

// Current возвращает текущий (глобальный) логгер
func Current() *Logger { return currentClog }

// EOF: "default.go"
