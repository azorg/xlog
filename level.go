// File: "level.go"

package xlog

import (
	"context"
	"fmt"
	"log/slog" // go>=1.21
	"runtime"
	"strconv"
	"strings"
	"time"
	// FIXME: "golang.org/x/exp/slog" // экспериментальный пакет для go=1.20 только
)

// Уровни логирования
//
//	DEBUG, INFO, WARN, ERROR - распространенные (стандартные) уровни (slog)
//	TRACE - распространенный уровень для трассировки программ
//	FLOOD - нестандартный уровень вывода избыточных данных в журнал
//	NOTICE, CRIT, ALERT, EMERG - уровни принятые в syslog
//	FATAL - уровень прерывания выполнения (по аналогии с log.Fatal())
//	PANIC - уровень паники
//	SILENT - фиктивный уровень для отключения журналирования вовсе
const (
	LevelFlood  = slog.Level(-12) // FLOOD  (-12)
	LevelTrace  = slog.Level(-8)  // TRACE  (-8)
	LevelDebug  = slog.LevelDebug // DEBUG  (-4)
	LevelInfo   = slog.LevelInfo  // INFO   (0)
	LevelNotice = slog.Level(2)   // NOTICE (2)
	LevelWarn   = slog.LevelWarn  // WARN   (4)
	LevelError  = slog.LevelError // ERROR  (8)
	LevelCrit   = slog.Level(10)  // CRIT   (10)
	LevelAlert  = slog.Level(12)  // ALERT  (12)
	LevelEmerg  = slog.Level(14)  // EMERG  (14)
	LevelFatal  = slog.Level(16)  // FATAL  (16)
	LevelPanic  = slog.Level(18)  // PANIC  (18)
	LevelSilent = slog.Level(20)  // SILENT (20)
)

// Уровень логирования по умолчанию
const DefaultLevel = slog.LevelInfo

// Строковые идентификаторы уровней логирования
// для представления в структуре конфигурации
const (
	LvlFlood  = "flood"
	LvlTrace  = "trace"
	LvlDebug  = "debug"
	LvlInfo   = "info"
	LvlNotice = "notice"
	LvlWarn   = "warn"
	LvlError  = "error"
	LvlCrit   = "crit"
	LvlAlert  = "alert"
	LvlEmerg  = "emerg"
	LvlFatal  = "fatal"
	LvlPanic  = "panic"
	LvlSilent = "silent"
)

// Метки уровней логирования для вывода в журнал
const (
	labelFlood  = "FLOOD"
	labelTrace  = "TRACE"
	labelDebug  = "DEBUG"
	labelInfo   = "INFO"
	labelNotice = "NOTICE"
	labelWarn   = "WARN"
	labelError  = "ERROR"
	labelCrit   = "CRIT"
	labelAlert  = "ALERT"
	labelEmerg  = "EMERG"
	labelFatal  = "FATAL"
	labelPanic  = "PANIC"
	labelSilent = "SILENT"
)

// Таблица преобразования идентификаторов уровней
var parseLvl = map[string]slog.Level{
	LvlFlood:  LevelFlood,
	LvlTrace:  LevelTrace,
	LvlDebug:  LevelDebug,
	LvlInfo:   LevelInfo,
	LvlNotice: LevelNotice,
	LvlWarn:   LevelWarn,
	LvlError:  LevelError,
	LvlCrit:   LevelCrit,
	LvlAlert:  LevelAlert,
	LvlEmerg:  LevelEmerg,
	LvlFatal:  LevelFatal,
	LvlPanic:  LevelPanic,
	LvlSilent: LevelSilent,
}

// Обратная таблица преобразования идентификаторов уровней
var parseLevel = map[slog.Level]string{
	LevelFlood:  LvlFlood,
	LevelTrace:  LvlTrace,
	LevelDebug:  LvlDebug,
	LevelInfo:   LvlInfo,
	LevelNotice: LvlNotice,
	LevelWarn:   LvlWarn,
	LevelError:  LvlError,
	LevelCrit:   LvlCrit,
	LevelAlert:  LvlAlert,
	LevelEmerg:  LvlEmerg,
	LevelFatal:  LvlFatal,
	LevelPanic:  LvlPanic,
	LevelSilent: LvlSilent,
}

// Обратная таблица преобразования метки уровня
// (используется в т.ч. как пополняемая кеш таблица)
var levelFromLabel = map[string]slog.Level{
	labelFlood:  LevelFlood,
	labelTrace:  LevelTrace,
	labelDebug:  LevelDebug,
	labelInfo:   LevelInfo,
	labelNotice: LevelNotice,
	labelWarn:   LevelWarn,
	labelError:  LevelError,
	labelCrit:   LevelCrit,
	labelAlert:  LevelAlert,
	labelEmerg:  LevelEmerg,
	labelFatal:  LevelFatal,
	labelPanic:  LevelPanic,
	labelSilent: LevelSilent,
}

// LevelToLabel преобразует уровень логирования к строке/метке для
// представления в журнале ("INFO", "ERROR" и др.) в стиле slog,
// подобно одноименному методам String() типов slog.Level/slog.LevelVar.
// Поддерживаются дополнительные уровни (TRACE, NOTICE, CRIT и др.).
func LevelToLabel(level slog.Level) string {
	str := func(base string, delta slog.Level) string {
		if delta == 0 {
			return base
		}
		return fmt.Sprintf("%s%+d", base, delta)
	}
	switch {
	case level < LevelTrace:
		return str(labelFlood, level-LevelFlood)
	case level < LevelDebug:
		return str(labelTrace, level-LevelTrace)
	case level < LevelInfo:
		return str(labelDebug, level-LevelDebug)
	case level < LevelNotice:
		return str(labelInfo, level-LevelInfo)
	case level < LevelWarn:
		return str(labelNotice, level-LevelNotice)
	case level < LevelError:
		return str(labelWarn, level-LevelWarn)
	case level < LevelCrit:
		return str(labelError, level-LevelError)
	case level < LevelAlert:
		return str(labelCrit, level-LevelCrit)
	case level < LevelEmerg:
		return str(labelAlert, level-LevelAlert)
	case level < LevelFatal:
		return str(labelEmerg, level-LevelEmerg)
	case level < LevelPanic:
		return str(labelFatal, level-LevelFatal)
	case level < LevelSilent:
		return str(labelPanic, level-LevelPanic)
	default: // level >= LevelSilent
		return str(labelSilent, level-LevelSilent)
	}
}

// LevelToColorLabel преобразует уровень логирования к строке/метке для
// представления в журнале с применением Escape/Ansi символов подсветки.
// Используется в TintHandler'е, если в структуре конфигурации Conf
// заданы Format="tinted" и Color=true.
// Поддерживаются дополнительные уровни (FLOOD, EMERG, ALERT и др.).
func LevelToColorLabel(level slog.Level) string {
	str := func(ansi, base string, delta slog.Level) string {
		if delta == 0 {
			return ansi + base + ansiReset
		}
		return fmt.Sprintf("%s%s%+d"+ansiReset, ansi, base, delta)
	}
	switch {
	case level < LevelTrace:
		return str(ansiFlood, labelFlood, level-LevelFlood)
	case level < LevelDebug:
		return str(ansiTrace, labelTrace, level-LevelTrace)
	case level < LevelInfo:
		return str(ansiDebug, labelDebug, level-LevelDebug)
	case level < LevelNotice:
		return str(ansiInfo, labelInfo, level-LevelInfo)
	case level < LevelWarn:
		return str(ansiNotice, labelNotice, level-LevelNotice)
	case level < LevelError:
		return str(ansiWarn, labelWarn, level-LevelWarn)
	case level < LevelCrit:
		return str(ansiError, labelError, level-LevelError)
	case level < LevelAlert:
		return str(ansiCrit, labelCrit, level-LevelCrit)
	case level < LevelEmerg:
		return str(ansiAlert, labelAlert, level-LevelAlert)
	case level < LevelFatal:
		return str(ansiEmerg, labelEmerg, level-LevelEmerg)
	case level < LevelPanic:
		return str(ansiFatal, labelFatal, level-LevelFatal)
	case level < LevelSilent:
		return str(ansiPanic, labelPanic, level-LevelPanic)
	default: // level >= LevelSilent
		return str(ansiPanic, labelSilent, level-LevelSilent)
	}
}

// LevelFromString преобразует строку идентификатор уровня логирования
// ("debug", "info", "0" и др.), используемый в структуре конфигурации
// к slog.Level. Функция не чувствительна в регистру. Уровень логирования
// может быть задан как строкой, так и десятичным целым числом.
func LevelFromString(level string) slog.Level {
	level = strings.ToLower(level)
	lvl, ok := parseLvl[level]
	if !ok {
		i, err := strconv.Atoi(level)
		if err != nil {
			return DefaultLevel
		}
		return slog.Level(i)
	}
	return lvl
}

// LevelToString преобразовывает численное значение уровня логгирования
// slog.Level к представлению виде строки в структуре конфигурации.
// Если задан не известный уровень, то возвращается его десятичное
// представление.
func LevelToString(level slog.Level) string {
	lvl, ok := parseLevel[level]
	if !ok {
		return fmt.Sprintf("%d", int(level))
	}
	return lvl
}

// LevelFromLabel преобразует метку уровня логирования (INFO, WARN, ...)
// в численное значение. Функция может быть востребована для парсинга логов.
// Входное значением может быть вида "ERROR+2", принятого в slog.
func LevelFromLabel(label string) slog.Level {
	label = strings.ToUpper(label)
	level, ok := levelFromLabel[label] // поиск значения в кеше
	if ok {
		return level
	}
	for level = slog.Level(-20); level <= slog.Level(20); level++ {
		if label == strings.ToUpper(LevelToLabel(level)) {
			levelFromLabel[label] = level // кешировать сведения
			return level
		}
	}
	return DefaultLevel
}

// logAttrs функция обёртка для реализации метода LogAttr для xlog.Logger
func logAttrs(
	ctx context.Context, log *slog.Logger,
	level slog.Level, msg string, attrs ...slog.Attr) error {

	if !log.Enabled(ctx, level) {
		return nil
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip wrappers

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.AddAttrs(attrs...)

	return log.Handler().Handle(ctx, r)
}

// logs - функция обёртка для реализации дополнительных уровней логирования
// с использованием структурированного логирования (Log, Trace, Notice, ...).
// Данная функция ДОЛЖНА вызываться из функций обёрток (см. "shugar.go").
func logs(
	ctx context.Context, log *slog.Logger,
	level slog.Level, msg string, args ...any) error {

	if !log.Enabled(ctx, level) {
		return nil
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip wrappers

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)

	return log.Handler().Handle(ctx, r)
}

// logf - функция обёртка для реализации дополнительных уровней логирования
// с использованием традиционного логирования (Logf, Infof, Debugf,
// Noticef, ...). Данная функция ДОЛЖНА вызываться из функций обёрток
// (см. "shugar.go").
func logf(
	ctx context.Context, log *slog.Logger,
	level slog.Level, format string, args ...any) error {

	if !log.Enabled(ctx, level) {
		return nil
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip wrappers

	r := slog.NewRecord(time.Now(), level, fmt.Sprintf(format, args...), pcs[0])

	return log.Handler().Handle(ctx, r)
}

// EOF: "level.go"
