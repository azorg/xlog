// File: "sprint.go"

package xlog

import (
	"fmt"
	"reflect"
	"time"
)

// Sprint - преобразует структуру данных в строку в формате
// близком к стандартному формату "%+v", но с обработкой
// указателей и вложенных структур.
// Используется рефлексия.
// Опционально могут поддерживаться JSON теги (если UseJSONTags=true)
// Функция используется для отображения структур данных в TintHandler'е.
//
//	val - значение произвольного типа, включая структуры, указатели
//	на структуры, ошибки и карты.
func Sprint(val any) string {
	return sprint("", val)
}

// Sprint - преобразует структуру данных в строку в формате
// близком к стандартному формату "%+v", но с обработкой
// указателей и вложенных структур.
// Используется рефлексия.
// Опционально могут поддерживаться JSON теги (если UseJSONTags=true)
//
//	prefix - префикс "&", если данная структура была доступна по указателю
func sprint(prefix string, val any) string {
	switch v := val.(type) {
	case error:
		return v.Error()

	case time.Time:
		return v.Format(RFC3339Micro)

	case time.Duration:
		return fmt.Sprintf("%v", v)
	}

	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		if v.IsZero() {
			return "<nil>"
		}
		elem := v.Elem()
		return sprint("&", elem.Interface())

	case reflect.Struct:
		buf := newBuffer()
		defer buf.Free()
		buf.WriteString(prefix + "{")
		delim := ""
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			name := field.Name
			if UseJSONTags {
				if tag := field.Tag.Get("json"); tag != "" {
					name = tag
				}
			}
			value := sprint("", v.Field(i).Interface())
			buf.WriteString(delim + name + ":" + value)
			delim = " "
		} // for
		buf.WriteString("}")
		return buf.String()

	default:
		return fmt.Sprintf("%+v", val)
	}
}

// EOF: "sprint.go"
