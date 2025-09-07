// File: "fields.go"

package xlog

import (
	"log/slog" // go>=1.21
	"time"
	// FIXME: "golang.org/x/exp/slog" // экспериментальный пакет для go=1.20 только
)

// Fields - это простая обертка для наполнения атрибутами записи
// в журнале на основе карт (key/value)
type Fields map[string]any

// FieldsProvider - общий интерфейс для предоставления атрибутов в виде key/value
type FieldsProvider interface {
	Fields() Fields
}

// Fields реализует попутно интерфейс FieldsProvider.
// Данное решение позволяет передавать Fields в качестве FieldsProvider.
func (f Fields) Fields() Fields {
	return f
}

// Args преобразует Fields в последовательность аргументов
// для вызова методов Log, Info, Debug, With и т.п.
// Значения nil исключаются из вывода.
func (fields Fields) Args() []any {
	args := make([]any, 0)
	for key, value := range fields {
		if value != nil && key != "" {
			args = append(args, key, value)
		}
	} // for
	return args
}

// Value - образует slog.Value из any
func Value(value any) slog.Value {
	switch val := value.(type) {
	case string:
		return slog.StringValue(val)
	case int:
		return slog.IntValue(val)
	case int64:
		return slog.Int64Value(val)
	case uint64:
		return slog.Uint64Value(val)
	case float64:
		return slog.Float64Value(val)
	case time.Time:
		return slog.TimeValue(val)
	case time.Duration:
		return slog.DurationValue(val)
	case bool:
		return slog.BoolValue(val)
	default:
		return slog.AnyValue(value)
	} // switch
}

// Attrs преобразует Fields в последовательность []slog.Attr
// для вызова высокоэффективного slog метода LogAttrs.
// Значение nil и пустые строки исключаются.
func (fields Fields) Attrs() []slog.Attr {
	as := make([]slog.Attr, 0)
	for key, value := range fields {
		if value != nil && key != "" { // не выводить <nil> и пустые строки
			if str, ok := value.(string); !ok || str != "" {
				as = append(as, slog.Attr{
					Key:   key,
					Value: Value(value),
				})
			}
		}
	} // for
	return as
}

// Value преобразует Fields в групповое значение с помощью slog.GroupValue()
func (fields Fields) Value() slog.Value {
	as := fields.Attrs()
	return slog.GroupValue(as...)
}

// EOF: "fields.go"
