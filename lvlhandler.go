// File: "lvlhandler.go"

package xlog

import (
	"context"
	"log/slog" // go>=1.21
	// FIXME: "golang.org/x/exp/slog" // экспериментальный пакет для go=1.20 только
)

// Вспомогательная обёртка для управления уровнем логирования стандартного
// slog логгера по умолчанию
type lvlHandler struct {
	handler slog.Handler
	leveler slog.Leveler
}

// Убедиться, что *lvlHandler соответствует интерфейсу slog.Handler
var _ slog.Handler = (*lvlHandler)(nil)

// Создать обертку slog.Handler'а с возможностью управления уровнем журналирования
func newLvlHandler(handler slog.Handler, leveler slog.Leveler) slog.Handler {
	// Optimization: avoid chains of logStdHandlers
	if sh, ok := handler.(*lvlHandler); ok {
		handler = sh.handler
	}
	return &lvlHandler{handler: handler, leveler: leveler}
}

// Enabled() требуется для интерфейса slog.Handler
func (h *lvlHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.leveler.Level()
}

// Handle() требуется для интерфейса slog.Handler
func (h *lvlHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.handler.Handle(ctx, r)
}

// WithAttrs() требуется для интерфейса slog.Handler
func (h *lvlHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return newLvlHandler(h.handler.WithAttrs(attrs), h.leveler)
}

// WithGroup()) требуется для интерфейса slog.Handler
func (h *lvlHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return newLvlHandler(h.handler.WithGroup(name), h.leveler)
}

// EOF: "lvlhandler.go"
