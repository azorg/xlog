// File: "middleware.go"

package xlog

import (
	"context"
	"log/slog" // go>=1.21
	// FIXME: "golang.org/x/exp/slog" // экспериментальный пакет для go=1.20 только
)

// HandlerFunc - это тип основной функции Handle() интерфейса slog.Handler
// для упаковки и отправки записи журнала в заданный канал/файл
type HandleFunc func(context.Context, slog.Record) error

// Middleware - описывает тип функции для построения
// обёрток для slog.Handler'в (т.н. middleware)
type Middleware func(next HandleFunc) HandleFunc

// MiddlewareFunc - описывает метод обёртку для построения Middleware
//
// Функция может обрабатывать входные значение ctx и r, после чего
// может вызывать следующий в цепочке Middleware или Handle.
type MiddlewareFunc func(
	ctx context.Context, r slog.Record,
	next HandleFunc) error

// Вспомогательная обёртка для создания Handler'ов с Middleware
type MiddlewareHandler struct {
	handler slog.Handler
	mws     []Middleware
}

// Убедиться, что *MiddlewareHandler соответствует интерфейсу slog.Handler
var _ slog.Handler = (*MiddlewareHandler)(nil)

// NewMiddleware создает новый Middlewre на основе заданного метода-обёртки
// с использованием замыканий
//
//	mwf - метод обёртки типа MiddlewareFunc
func NewMiddleware(mwf MiddlewareFunc) Middleware {
	return func(next HandleFunc) HandleFunc {
		return func(ctx context.Context, r slog.Record) error {
			return mwf(ctx, r, next)
		}
	}
}

// Создать обертку MiddlewareHandler'а с заданными middleware(s)
func NewMiddlewareHandler(handler slog.Handler, mws ...Middleware) *MiddlewareHandler {
	return &MiddlewareHandler{handler: handler, mws: mws}
}

// Enabled() требуется для интерфейса slog.Handler
func (h *MiddlewareHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle() требуется для интерфейса slog.Handler
func (h *MiddlewareHandler) Handle(ctx context.Context, r slog.Record) error {
	if len(h.mws) == 0 { // нет middleware (???)
		return h.handler.Handle(ctx, r)
	}

	// Использовать цепочку middleware
	handle := h.mws[len(h.mws)-1](h.handler.Handle)
	for i := len(h.mws) - 2; i >= 0; i-- {
		handle = h.mws[i](handle)
	}
	return handle(ctx, r)
}

// WithAttrs() требуется для интерфейса slog.Handler
func (h *MiddlewareHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return NewMiddlewareHandler(h.handler.WithAttrs(attrs), h.mws...)
}

// WithGroup() требуется для интерфейса slog.Handler
func (h *MiddlewareHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return NewMiddlewareHandler(h.handler.WithGroup(name), h.mws...)
}

// Пример Middleware, который дублирует вывод записей в
// дополнительный логгер (имитация режима "Multi Handler")
//
//	log - дополнительный логгер в который направляются сообщения
func NewMiddlewareMulti(log *Logger) Middleware {
	handler := log.Logger.Handler()
	mwf := func(ctx context.Context, r slog.Record, next HandleFunc) error {
		if handler.Enabled(ctx, r.Level) {
			handler.Handle(ctx, r) // вызвать дополнительный хендлер
		}
		return next(ctx, r)
	}
	return NewMiddleware(mwf)
}

// Пример Middleware, который добавляет (в начало) записи заданные
// дополнительные поля. Пример не корректно работает с группами -
// все поля добавляются в последнюю открытую группу.
// Данная функция приведена скорее для примера использования Middleware,
// чем для практического применения.
func NewMiddlewareWithFields(fields FieldsProvider) Middleware {
	mwf := func(ctx context.Context, r slog.Record, handle HandleFunc) error {
		if fields == nil {
			return handle(ctx, r)
		}

		// Создать новую запись
		rNew := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)

		// Создать новый список трибутов (в начале дополнительные поля)
		attrs := fields.Fields().Attrs()

		// Получить список атрибутов из старой записи
		r.Attrs(func(attr slog.Attr) bool {
			attrs = append(attrs, attr)
			return true
		})

		// Скопировтаь атрибуты в новую запись
		for _, attr := range attrs {
			rNew.AddAttrs(attr)
		}
		return handle(ctx, rNew)
	}
	return NewMiddleware(mwf)
}

// Пример Middleware, который все сообщения с ошибками
// (с атрибутом "err") дополнительно направляет в заданный логгер.
// Данная функция приведена скорее для примера использования Middleware,
// чем для практического применения.
func NewMiddlewareForError(logErr *Logger) Middleware {
	if logErr == nil {
		Panic("logger required (logErr is nil)")
	}
	mwf := func(ctx context.Context, r slog.Record, handle HandleFunc) error {
		hasError := false
		r.Attrs(func(attr slog.Attr) bool {
			if attr.Key == ErrKey { // attr.Key == "err"
				if err, ok := attr.Value.Any().(error); err != nil && ok {
					hasError = true
					return false
				}
			}
			return true
		})
		if hasError {
			// Создать новую запись для logErr "с пометкой"
			rNew := slog.NewRecord(r.Time, r.Level, "ERROR: "+r.Message, r.PC)

			// Скопировтаь атрибуты в новую запись
			r.Attrs(func(attr slog.Attr) bool {
				rNew.AddAttrs(attr)
				return true
			})

			// Отправить запись в специальный журнал logErr
			logErr.Logger.Handler().Handle(ctx, rNew)
		}
		return handle(ctx, r)
	}
	return NewMiddleware(mwf)
}

// Пример middleware, который заменяет значения атрибутов
// passwd, password на ********.
// Данная функция приведена скорее для примера использования Middleware,
// чем для практического применения.
func NewMiddlewareNoPasswd() Middleware {
	mwf := func(ctx context.Context, r slog.Record, handle HandleFunc) error {
		// Создать новую запись
		rNew := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)

		// Скопировтаь атрибуты в новую запись с подменой passwd и password
		r.Attrs(func(attr slog.Attr) bool {
			if (attr.Key == "passwd" || attr.Key == "password") &&
				(attr.Value.Kind() == slog.KindString || attr.Value.Kind() == slog.KindInt64) {
				attr.Value = slog.StringValue("********")
			}
			rNew.AddAttrs(attr)
			return true
		})
		return handle(ctx, rNew)
	}
	return NewMiddleware(mwf)
}

// EOF: "middleware.go"
