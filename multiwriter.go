// File: "multiwriter.go"

package xlog

import "io"

// MultiWriter - это обёртка для io.Writer для направления журналов по
// нескольким направлениям. Может использоваться, если необходимо
// отправлять журналы по нескольким направлениям одновременно
// (например syslog + NATS).
type MultiWriter []io.Writer

// Убедится в том, что MultiWriter соответствуют интерфейсу io.Writer
var _ io.Writer = NewMultiWriter()

// NewMultiWriter cоздает "пустой" MutliWriter.
// Далее в MultiWriter могут быть добавлены io.Writer'а с помощью
// метода Add.
func NewMultiWriter() MultiWriter { return MultiWriter([]io.Writer{}) }

// Add добавляет заданный io.Writer в MultiWriter
func (mw MultiWriter) Add(w io.Writer) { mw = append(mw, w) }

// Write реализует интерфейс io.Writer для MultiWriter'а.
// Производится последовательная запись данных data во все
// io.Writer'ы MultiWriter'а. Ошибки не возвращаются.
func (mw MultiWriter) Write(data []byte) (int, error) {
	for _, w := range mw {
		w.Write(data)
	}
	return len(data), nil // всегда возвращаем успех
}

// EOF: "multiwriter.go"
