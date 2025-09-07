// File: "tint.go"

package xlog

import (
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"log/slog" // go>=1.21
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
	// FIXME: "golang.org/x/exp/slog" // экспериментальный пакет для go=1.20 только
)

// Символ перевода строки
const newLineChar = '\n'

// Структура конфигурации для создания TintHandler'а
type TintOptions struct {
	// Управление уровнем журналирования
	// (по умолчанию DefaultLevel=slog.LevelInfo).
	// Сообщения ниже заданного уровня не выводятся в журнал.
	Level slog.Leveler

	// Добавить ссылки на исходный код (файл:строка)
	AddSource bool

	// Добавить имя пакета к имени файла исходного текста (пакет/файл:строка)
	SourcePkg bool

	// Добавить имя функций/методов
	SourceFunc bool

	// Удалить расширение ".go" из имени файла
	NoExt bool

	// Отключить вывод метки уровня журналирования
	NoLevel bool

	// Отключить подсветку и исключить вывод ANSI/Escape символов
	NoColor bool

	// Выводить метку времени в UTC
	TimeUTC bool

	// Формат метки времени
	TimeFormat string

	// Функция "хук" для подмены атрибутов перед формированием
	// записей в журнале.
	// См. https://pkg.go.dev/log/slog#HandlerOptions for details.
	ReplaceAttr func(groups []string, attr slog.Attr) slog.Attr
}

// Структура данных TintHandler, соответствующего интерфейсу slog.Handler.
// TintHandler - это минималистский slog.Handler с подсветкой на основе
// исходников с "github.com/lmittmann/tint"
// (https://github.com/lmittmann/tint/blob/main/handler.go).
//
// Особенности исходного проекта:
//   - MIT License
//   - отсутствие зависимостей (только slog)
//   - компактный выходной формат схож с zerolog.ConsoleWriter и slog.TextHandler
//   - "тонирование" меток и ключей атрибутов
//   - подсветка ошибок
//   - всю подсветку на основе ANSII символов можно отключить
//   - поддержка `ReplaceAtt` как у slog.TextHandler/slog.JSONHandler
//
// Что изменено в рамках xlog:
//   - упрощена подкраска ошибок
//   - добавлен вывод имени пакета/функции (по опциям: sourcePkg/source/Func)
//   - есть возможность отключить метку уровня (noLevel)
//   - добавлена возможность исключения вывода расширения файла ".go"
//   - добавлена возможность вывода метки времени в UTC
//   - добавлена возможность вывода вложенных структур,
//     в т.ч. по указателям (см. функцию String())
//   - несколько улучшено представление чисел с плавающей точкой в журнале (как в JSON)
//   - время в атрибутах выводится в формате time.RFC3339Nano (а не просто как t.String())
type TintHandler struct {
	mu sync.Mutex
	w  io.Writer

	level      slog.Leveler
	addSource  bool
	sourcePkg  bool
	sourceFunc bool
	noExt      bool
	noLevel    bool
	noColor    bool
	timeUTC    bool
	timeFormat string

	attrsPrefix string
	groupPrefix string
	groups      []string

	replaceAttr func([]string, slog.Attr) slog.Attr
}

// Убедиться, что *TintHandler реализует интерфейс slog.Handler
var _ slog.Handler = (*TintHandler)(nil)

// Создать новый Tinted хендлер, соответствующий slog.Handler'у
func NewTintHandler(w io.Writer, opts *TintOptions) *TintHandler {
	h := &TintHandler{
		w:          w,
		level:      DefaultLevel,
		timeFormat: timeOff,
	}

	if opts == nil {
		return h
	}

	h.addSource = opts.AddSource
	h.sourcePkg = opts.SourcePkg
	h.sourceFunc = opts.SourceFunc
	h.noExt = opts.NoExt
	h.noLevel = opts.NoLevel
	h.noColor = opts.NoColor
	h.timeUTC = opts.TimeUTC
	h.replaceAttr = opts.ReplaceAttr

	if opts.Level != nil {
		h.level = opts.Level
	}

	if opts.TimeFormat != "" {
		h.timeFormat = opts.TimeFormat
	}
	return h
}

// clone копирует хендлер
func (h *TintHandler) clone() *TintHandler {
	return &TintHandler{
		w:           h.w,
		level:       h.level,
		addSource:   h.addSource,
		sourcePkg:   h.sourcePkg,
		sourceFunc:  h.sourceFunc,
		noExt:       h.noExt,
		noLevel:     h.noLevel,
		noColor:     h.noColor,
		timeUTC:     h.timeUTC,
		timeFormat:  h.timeFormat,
		attrsPrefix: h.attrsPrefix,
		groupPrefix: h.groupPrefix,
		groups:      h.groups,
		replaceAttr: h.replaceAttr,
	}
}

// Метод Enabled() реализует интерфейс slog.Handler
func (h *TintHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

// Подготовить текстовый буфер для записи в журнал
func (h *TintHandler) format(r slog.Record) []byte {
	// Сформировать буфер "sync pool"
	buf := newBuffer()
	defer buf.Free()

	rep := h.replaceAttr

	// Добавить метку времени
	if h.timeFormat != timeOff && !r.Time.IsZero() {
		val := r.Time.Round(0) // strip monotonic to match Attr behavior
		if rep == nil {
			h.appendTime(buf, r.Time)
			buf.WriteByte(' ')
		} else if a := rep(nil /* groups */, slog.Time(slog.TimeKey, val)); a.Key != "" {
			if a.Value.Kind() == slog.KindTime {
				h.appendTime(buf, a.Value.Time())
			} else {
				h.appendValue(buf, a.Value, false)
			}
			buf.WriteByte(' ')
		}
	}

	// Добавить метку уровня сообщения
	if !h.noLevel {
		if rep == nil {
			h.appendLevel(buf, r.Level)
			buf.WriteByte(' ')
		} else if a := rep(nil /* groups */, slog.Any(slog.LevelKey, r.Level)); a.Key != "" {
			h.appendValue(buf, a.Value, false)
			buf.WriteByte(' ')
		}
	}

	// Добавить ссылку на исходный текст
	if h.addSource {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			src := &slog.Source{
				Function: f.Function,
				File:     f.File,
				Line:     f.Line,
			}
			if !h.sourcePkg {
				src.File = path.Base(src.File) // only file name
			}
			if h.noExt { // remove ".go" extension
				src.File = removeGoExt(src.File)
			}
			if h.sourceFunc { // add function name
				funcName := getFuncName(5) // skip=5 (some magic)
				if funcName != "" {
					src.File += ":" + funcName + "()"
				}
			}
			if rep == nil {
				h.appendSource(buf, src)
				buf.WriteByte(' ')
			} else if a := rep(nil /* groups */, slog.Any(slog.SourceKey, src)); a.Key != "" {
				h.appendValue(buf, a.Value, false)
				buf.WriteByte(' ')
			}
		}
	}

	// Записать текст сообщения
	if rep == nil {
		buf.WriteString(r.Message)
		buf.WriteByte(' ')
	} else if a := rep(nil /* groups */, slog.String(slog.MessageKey, r.Message)); a.Key != "" {
		h.appendValue(buf, a.Value, false)
		buf.WriteByte(' ')
	}

	// Записать атрибуты префикса
	if len(h.attrsPrefix) > 0 {
		buf.WriteString(h.attrsPrefix)
	}

	// Записать атрибуты сообщения
	r.Attrs(func(attr slog.Attr) bool {
		h.appendAttr(buf, attr, h.groupPrefix, h.groups)
		return true
	})

	return *buf
}

// Форматировать slog запись в строку (только для экспериментов).
// Данный метод не требуется для реализации интерфейса slog.Handler.
// FIXME: Удалить данный артефакт.
func (h *TintHandler) formatToString(r slog.Record) string {
	buf := h.format(r)

	size := len(buf)
	if size == 0 {
		return ""
	}

	// Исключить последний пробел
	return string(buf[:size-1])
}

// Метод Handle() реализует интерфейс slog.Handler
func (h *TintHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := h.format(r)
	if len(buf) == 0 {
		return nil
	}

	// Заменить последний пробел на символ перевода строки
	buf[len(buf)-1] = newLineChar

	// Произвести запись буфера в выходной канал/файл
	h.mu.Lock()
	_, err := h.w.Write(buf)
	h.mu.Unlock()
	return err
}

// Метод WithAttrs() реализует интерфейс slog.Handler
func (h *TintHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := h.clone()

	buf := newBuffer()
	defer buf.Free()

	// Write attributes to buffer
	for _, attr := range attrs {
		h.appendAttr(buf, attr, h.groupPrefix, h.groups)
	}
	h2.attrsPrefix = h.attrsPrefix + string(*buf)
	return h2
}

// Метод WithGroup() реализует интерфейс slog.Handler
func (h *TintHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := h.clone()
	h2.groupPrefix += name + "."
	h2.groups = append(h2.groups, name)
	return h2
}

// appendTime добавляет метку времени в буфер
func (h *TintHandler) appendTime(buf *buffer, t time.Time) {
	if h.timeUTC {
		t = t.UTC()
	}

	if tintAlignTime { // FIXME: slow (may append zeros)
		fmt := h.timeFormat
		fmtLen := len(fmt)
		time := t.Format(fmt)
		if len(fmt) != 0 && fmt[fmtLen-1] == '9' { // FIXME: bad magic code
			if addZeros := fmtLen - len(time); addZeros > 0 {
				time += strings.Repeat("0", addZeros)
			}
		}
		if h.noColor {
			buf.WriteString(time)
		} else {
			buf.WriteString(ansiTime)
			buf.WriteString(time)
			buf.WriteString(ansiReset)
		}
	} else { // fast
		if h.noColor {
			*buf = t.AppendFormat(*buf, h.timeFormat)
		} else {
			buf.WriteString(ansiTime)
			*buf = t.AppendFormat(*buf, h.timeFormat)
			buf.WriteString(ansiReset)
		}
	}
}

// appendLevel добавляет метку уровня журналирования в буфер
func (h *TintHandler) appendLevel(buf *buffer, level slog.Level) {
	if h.noColor {
		buf.WriteString(LevelToLabel(level))
	} else {
		buf.WriteString(LevelToColorLabel(level))
	}
}

// appensSource добавляет ссылку на исходный текст в буфер
func (h *TintHandler) appendSource(buf *buffer, src *slog.Source) {
	if !h.noColor {
		buf.WriteString(ansiSource)
		defer buf.WriteString(ansiReset)
	}

	dir, file := filepath.Split(src.File)
	buf.WriteString(filepath.Join(filepath.Base(dir), file))
	buf.WriteByte(':')
	buf.WriteString(strconv.Itoa(src.Line))
}

// appendAttr добавляет запись ключ/значение в буфер
func (h *TintHandler) appendAttr(buf *buffer, attr slog.Attr,
	groupsPrefix string, groups []string) {

	attr.Value = attr.Value.Resolve()
	if rep := h.replaceAttr; rep != nil && attr.Value.Kind() != slog.KindGroup {
		attr = rep(groups, attr)
		attr.Value = attr.Value.Resolve()
	}

	if attr.Equal(slog.Attr{}) {
		return // skip empty
	}

	if attr.Value.Kind() == slog.KindGroup {
		// Append group
		if attr.Key != "" {
			groupsPrefix += attr.Key + "."
			groups = append(groups, attr.Key)
		}
		for _, groupAttr := range attr.Value.Group() {
			h.appendAttr(buf, groupAttr, groupsPrefix, groups)
		}
		return
	}

	if err, ok := attr.Value.Any().(error); err != nil && ok {
		// Append error
		h.appendError(buf, attr.Key, err, groupsPrefix)
		buf.WriteByte(' ')
		return
	}

	h.appendKey(buf, attr.Key, groupsPrefix)
	h.appendValue(buf, attr.Value, true)
	buf.WriteByte(' ')
}

// apendKey добавляет "ключ=" в буфер
func (h *TintHandler) appendKey(buf *buffer, key, groups string) {
	buf.WriteStringIf(!h.noColor, ansiKey)
	appendString(buf, groups+key, true, !h.noColor)
	buf.WriteByte('=')
	buf.WriteStringIf(!h.noColor, ansiReset)
}

// appendValue добавляет значение после "key=" в буфер
func (h *TintHandler) appendValue(buf *buffer, v slog.Value, quote bool) {
	switch v.Kind() {
	case slog.KindString:
		appendString(buf, v.String(), quote, !h.noColor)
	case slog.KindInt64:
		*buf = strconv.AppendInt(*buf, v.Int64(), 10)
	case slog.KindUint64:
		*buf = strconv.AppendUint(*buf, v.Uint64(), 10)
	case slog.KindFloat64:
		//*buf = strconv.AppendFloat(*buf, v.Float64(), 'g', -1, 64) // так плохо
		// Небольшое улучшение представления чисел с плавающей точкой в журнале
		data, _ := json.Marshal(v.Float64())
		*buf = append(*buf, data...)
	case slog.KindBool:
		*buf = strconv.AppendBool(*buf, v.Bool())
	case slog.KindDuration:
		appendString(buf, v.Duration().String(), quote, !h.noColor)
	case slog.KindTime:
		//appendString(buf, v.Time().String(), quote, !h.noColor) // так плохо
		appendString(buf, v.Time().Format(time.RFC3339Nano), quote, !h.noColor) // как в JSON
	case slog.KindAny:
		defer func() {
			// Copied from log/slog/handler.go.
			if r := recover(); r != nil {
				// If it panics with a nil pointer, the most likely cases are
				// an encoding.TextMarshaler or error fails to guard against nil,
				// in which case "<nil>" seems to be the feasible choice.
				//
				// Adapted from the code in fmt/print.go.
				if v := reflect.ValueOf(v.Any()); v.Kind() == reflect.Pointer && v.IsNil() {
					appendString(buf, "<nil>", false, false) // quote=false color=false
					return
				}

				// Otherwise just print the original panic message.
				appendString(buf, fmt.Sprintf("!PANIC: %v", r), true, !h.noColor)
			}
		}()

		switch cv := v.Any().(type) {
		case slog.Level:
			h.appendLevel(buf, cv)
		case encoding.TextMarshaler:
			data, err := cv.MarshalText()
			if err != nil {
				break
			}
			appendString(buf, string(data), quote, !h.noColor)
		case *slog.Source:
			h.appendSource(buf, cv)
		default:
			// Оригинальный код:
			//appendString(buf, fmt.Sprintf("%+v", cv), quote, !h.noColor)

			// Модернизированный код для форматирования структур
			appendString(buf, Sprint(cv), quote, !h.noColor)
		} // switch
	}
}

// appendError добавляет "err=..." в буффер
func (h *TintHandler) appendError(buf *buffer, key string, err error, groupsPrefix string) {
	buf.WriteStringIf(!h.noColor, ansiErrKey)
	appendString(buf, groupsPrefix+key, true, !h.noColor)
	buf.WriteByte('=')
	buf.WriteStringIf(!h.noColor, ansiErrVal)
	appendString(buf, err.Error(), true, !h.noColor)
	buf.WriteStringIf(!h.noColor, ansiReset)
}

// appendString добавляет строку, при необходимости в кавычках в буфер
func appendString(buf *buffer, s string, quote, color bool) {
	if quote && !color {
		// Trim ANSI escape sequences
		var inEscape bool
		s = cut(s, func(r rune) bool {
			if r == ansiEsc {
				inEscape = true
			} else if inEscape && unicode.IsLetter(r) {
				inEscape = false
				return true
			}

			return inEscape
		})
	}

	quote = quote && needsQuoting(s)
	switch {
	case color && quote:
		s = strconv.Quote(s)
		s = strings.ReplaceAll(s, `\x1b`, string(ansiEsc))
		buf.WriteString(s)
	case !color && quote:
		*buf = strconv.AppendQuote(*buf, s)
	default:
		buf.WriteString(s)
	}
}

func cut(s string, f func(r rune) bool) string {
	var res []rune
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError {
			break
		}
		if !f(r) {
			res = append(res, r)
		}
		i += size
	}
	return string(res)
}

// needsQuoting проверят нужно ли экранировать строку кавычками.
// Скопировано из log/slog/text_handler.go (ранее была более простая версия).
func needsQuoting(s string) bool {
	if len(s) == 0 {
		return true
	}
	for i := 0; i < len(s); {
		b := s[i]
		if b < utf8.RuneSelf {
			// Quote anything except a backslash that would need quoting in a
			// JSON string, as well as space and '='
			if b != '\\' && (b == ' ' || b == '=' || !safeSet[b]) {
				return true
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError || unicode.IsSpace(r) || !unicode.IsPrint(r) {
			return true
		}
		i += size
	}
	return false
}

// Copied from log/slog/json_handler.go.
//
// safeSet is extended by the ANSI escape code "\u001b".
var safeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      true,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      true,
	'=':      true,
	'>':      true,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
	'\u001b': true,
}

// EOF: "tint.go"
