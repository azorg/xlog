// File: "tint.go"
// Tinted (colorized) logger based on "github.com/lmittmann/tint" sources

package xlog

import (
	"context"
	"encoding"
	"fmt"
	"io"
	"log/slog" // go>=1.21
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	//"golang.org/x/exp/slog" // depricated for go>=1.21
)

// Time formats
const (
	// Time OFF
	TIME_OFF = ""

	// Default time format of standart logger
	STD_TIME = "2006/01/02 15:04:05"

	// Default time format of standart logger + microseconds
	STD_TIME_US = "2006/01/02 15:04:05.999999"

	// Default time format of standart logger + milliseconds
	STD_TIME_MS = "2006/01/02 15:04:05.999"

	// RFC3339 time format + nanoseconds (slog.TextHandler by default)
	RFC3339Nano = time.RFC3339Nano // "2006-01-02T15:04:05.999999999Z07:00"

	// RFC3339 time format + microseconds
	RFC3339Micro = "2006-01-02T15:04:05.999999Z07:00"

	// RFC3339 time format + milliseconds
	RFC3339Milli = "2006-01-02T15:04:05.999Z07:00"

	// Time only format + microseconds
	TimeOnlyMicro = "15:04:05.999999"

	// Time only format + milliseconds
	TimeOnlyMilli = "15:04:05.999"

	// Default (recomented) time format wuth milliseconds
	DEFAULT_TIME_FORMAT = STD_TIME_MS

	// Default (recomented) time format witn microseconds
	DEFAULT_TIME_FORMAT_US = STD_TIME_US
)

type TintOptions struct {
	// Enable source code location
	AddSource bool

	// Log long file path (directory + file name)
	SourceLong bool

	// Minimum level to log (Default: slog.LevelInfo)
	Level slog.Leveler

	// Off level keys
	NoLevel bool

	// ReplaceAttr is called to rewrite each non-group attribute before it is logged.
	// See https://pkg.go.dev/log/slog#HandlerOptions for details.
	ReplaceAttr func(groups []string, attr slog.Attr) slog.Attr

	// Time format
	TimeFormat string

	// Disable color
	NoColor bool
}

// Tinted (colorized) handler implements a slog.Handler
type TintHandler struct {
	attrsPrefix string
	groupPrefix string
	groups      []string

	mu sync.Mutex
	w  io.Writer

	addSource   bool
	sourceLong  bool
	level       slog.Leveler
	noLevel     bool
	replaceAttr func([]string, slog.Attr) slog.Attr
	timeFormat  string
	noColor     bool
}

// Create new tinted (colorized) handler
func NewTintHandler(w io.Writer, opts *TintOptions) *TintHandler {
	h := &TintHandler{
		w:          w,
		level:      DEFAULT_LEVEL,
		timeFormat: TIME_OFF,
	}

	if opts == nil {
		return h
	}

	h.addSource = opts.AddSource
	h.sourceLong = opts.SourceLong
	h.noLevel = opts.NoLevel
	h.noColor = opts.NoColor
	h.replaceAttr = opts.ReplaceAttr

	if opts.Level != nil {
		h.level = opts.Level
	}

	if opts.TimeFormat != "" {
		h.timeFormat = opts.TimeFormat
	}
	return h
}

// Enabled() implements slog.Handler interface
func (h *TintHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

// Handle() implements slog.Handler interface
func (h *TintHandler) Handle(_ context.Context, r slog.Record) error {
	// Get a buffer from the sync pool
	buf := NewBuffer()
	defer buf.Free()

	rep := h.replaceAttr

	// Write time
	if h.timeFormat != TIME_OFF && !r.Time.IsZero() {
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

	// Write level
	if !h.noLevel {
		if rep == nil {
			h.appendLevel(buf, r.Level)
			buf.WriteByte(' ')
		} else if a := rep(nil /* groups */, slog.Any(slog.LevelKey, r.Level)); a.Key != "" {
			h.appendValue(buf, a.Value, false)
			buf.WriteByte(' ')
		}
	}

	// Write source
	if h.addSource {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			src := &slog.Source{
				Function: f.Function,
				File:     f.File,
				Line:     f.Line,
			}
			if !h.sourceLong {
				src.File = path.Base(src.File) // only file name
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

	// Write message
	if rep == nil {
		buf.WriteString(r.Message)
		buf.WriteByte(' ')
	} else if a := rep(nil /* groups */, slog.String(slog.MessageKey, r.Message)); a.Key != "" {
		h.appendValue(buf, a.Value, false)
		buf.WriteByte(' ')
	}

	// Write handler attributes
	if len(h.attrsPrefix) > 0 {
		buf.WriteString(h.attrsPrefix)
	}

	// Write attributes
	r.Attrs(func(attr slog.Attr) bool {
		h.appendAttr(buf, attr, h.groupPrefix, h.groups)
		return true
	})

	if len(*buf) == 0 {
		return nil
	}
	(*buf)[len(*buf)-1] = '\n' // replace last space with newline

	h.mu.Lock()
	defer h.mu.Unlock()

	_, err := h.w.Write(*buf)
	return err
}

// WithAttrs() implements slog.Handler interface
func (h *TintHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := h.clone()

	buf := NewBuffer()
	defer buf.Free()

	// write attributes to buffer
	for _, attr := range attrs {
		h.appendAttr(buf, attr, h.groupPrefix, h.groups)
	}
	h2.attrsPrefix = h.attrsPrefix + string(*buf)
	return h2
}

// WithGroup() implements slog.Handler interface
func (h *TintHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := h.clone()
	h2.groupPrefix += name + "."
	h2.groups = append(h2.groups, name)
	return h2
}

func (h *TintHandler) clone() *TintHandler {
	return &TintHandler{
		attrsPrefix: h.attrsPrefix,
		groupPrefix: h.groupPrefix,
		groups:      h.groups,
		w:           h.w,
		addSource:   h.addSource,
		sourceLong:  h.sourceLong,
		level:       h.level,
		replaceAttr: h.replaceAttr,
		timeFormat:  h.timeFormat,
		noColor:     h.noColor,
	}
}

func (h *TintHandler) appendTime(buf *Buffer, t time.Time) {
	if TINT_ALIGN_TIME { // slow (append zeros)
		time := t.Format(h.timeFormat)
		time += strings.Repeat("0", len(h.timeFormat)-len(time))
		if h.noColor {
			buf.WriteString(time)
		} else {
			buf.WriteString(AnsiTime)
			buf.WriteString(time)
			buf.WriteString(AnsiReset)
		}
	} else { // fast
		if h.noColor {
			*buf = t.AppendFormat(*buf, h.timeFormat)
		} else {
			buf.WriteString(AnsiTime)
			*buf = t.AppendFormat(*buf, h.timeFormat)
			buf.WriteString(AnsiReset)
		}
	}
}

func (h *TintHandler) appendLevel(buf *Buffer, level slog.Level) {
	xl := Level(level)
	if h.noColor {
		buf.WriteString(xl.String())
	} else {
		buf.WriteString(xl.ColorString())
	}
}

func (h *TintHandler) appendSource(buf *Buffer, src *slog.Source) {
	if !h.noColor {
		buf.WriteString(AnsiSource)
		defer buf.WriteString(AnsiReset)
	}
	dir, file := filepath.Split(src.File)
	buf.WriteString(filepath.Join(filepath.Base(dir), file))
	buf.WriteByte(':')
	buf.WriteString(strconv.Itoa(src.Line))
}

func (h *TintHandler) appendAttr(buf *Buffer, attr slog.Attr,
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
		if attr.Key != "" {
			groupsPrefix += attr.Key + "."
			groups = append(groups, attr.Key)
		}
		for _, groupAttr := range attr.Value.Group() {
			h.appendAttr(buf, groupAttr, groupsPrefix, groups)
		}
	} else if err, ok := attr.Value.Any().(error); err != nil && ok {
		// appen error
		h.appendError(buf, attr.Key, err, groupsPrefix)
		buf.WriteByte(' ')
	} else {
		h.appendKey(buf, attr.Key, groupsPrefix)
		h.appendValue(buf, attr.Value, true)
		buf.WriteByte(' ')
	}
}

func (h *TintHandler) appendKey(buf *Buffer, key, groups string) {
	if !h.noColor {
		buf.WriteString(AnsiKey)
		defer buf.WriteString(AnsiReset)
	}
	appendString(buf, groups+key, true)
	buf.WriteByte('=')
}

func (h *TintHandler) appendValue(buf *Buffer, v slog.Value, quote bool) {
	switch v.Kind() {
	case slog.KindString:
		appendString(buf, v.String(), quote)
	case slog.KindInt64:
		*buf = strconv.AppendInt(*buf, v.Int64(), 10)
	case slog.KindUint64:
		*buf = strconv.AppendUint(*buf, v.Uint64(), 10)
	case slog.KindFloat64:
		*buf = strconv.AppendFloat(*buf, v.Float64(), 'g', -1, 64)
	case slog.KindBool:
		*buf = strconv.AppendBool(*buf, v.Bool())
	case slog.KindDuration:
		appendString(buf, v.Duration().String(), quote)
	case slog.KindTime:
		appendString(buf, v.Time().String(), quote)
	case slog.KindAny:
		switch cv := v.Any().(type) {
		case slog.Level:
			h.appendLevel(buf, cv)
		case encoding.TextMarshaler:
			data, err := cv.MarshalText()
			if err != nil {
				break
			}
			appendString(buf, string(data), quote)
		case *slog.Source:
			h.appendSource(buf, cv)
		default:
			appendString(buf, fmt.Sprint(v.Any()), quote)
		}
	}
}

func (h *TintHandler) appendError(buf *Buffer, key string, err error, groupsPrefix string) {
	buf.WriteStringIf(!h.noColor, AnsiErrKey)
	appendString(buf, groupsPrefix+key, true)
	buf.WriteByte('=')
	buf.WriteStringIf(!h.noColor, AnsiErrVal)
	appendString(buf, err.Error(), true)
	buf.WriteStringIf(!h.noColor, AnsiReset)
}

func appendString(buf *Buffer, s string, quote bool) {
	if quote && needsQuoting(s) {
		*buf = strconv.AppendQuote(*buf, s)
	} else {
		buf.WriteString(s)
	}
}

func needsQuoting(s string) bool {
	if len(s) == 0 {
		return true
	}
	for _, r := range s {
		if unicode.IsSpace(r) || r == '"' || r == '=' || !unicode.IsPrint(r) {
			return true
		}
	}
	return false
}

// EOF: "tint.go"
