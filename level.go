// File: "level.go"

package xlog

import (
	"context"
	"fmt"
	//"log/slog" // go>=1.21
	"runtime"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slog" // deprecated for go>=1.21
)

// xlog level delivered from slog.Level, implements slog.Leveler interface
type Level slog.Level

// xlog leveler interface
// (slog.Leveler + Update() + String()/ColorString() methods)
type Leveler interface {
	Level() slog.Level   // get log level as slog.Level (implements a slog.Leveler interface)
	Update(slog.Level)   // update log level
	String() string      // get log level as label
	ColorString() string // get log level as color label
}

// Ensure *Level implements Leveler interface
var _ Leveler = (*Level)(nil)

// Log levels delivered from slog.Level
const (
	LevelFlood    = slog.Level(-12) // FLOOD    (-12)
	LevelTrace    = slog.Level(-8)  // TRACE    (-8)
	LevelDebug    = slog.LevelDebug // DEBUG    (-4)
	LevelInfo     = slog.LevelInfo  // INFO     (0)
	LevelNotice   = slog.Level(2)   // NOTICE   (2)
	LevelWarn     = slog.LevelWarn  // WARN     (4)
	LevelError    = slog.LevelError // ERROR    (8)
	LevelCritical = slog.Level(12)  // CRITICAL (12)
	LevelFatal    = slog.Level(16)  // FATAL    (16)
	LevelPanic    = slog.Level(18)  // PANIC    (18)
	LevelSilent   = slog.Level(20)  // SILENT   (20)
)

const DEFAULT_LEVEL = LevelInfo

// Log level as string for setup
const (
	LvlFlood    = "flood"
	LvlTrace    = "trace"
	LvlDebug    = "debug"
	LvlInfo     = "info"
	LvlNotice   = "notice"
	LvlWarn     = "warn"
	LvlError    = "error"
	LvlCritical = "critical"
	LvlFatal    = "fatal"
	LvlPanic    = "panic"
	LvlSilent   = "silent"
)

// Log level tags
const (
	LabelFlood    = "FLOOD"
	LabelTrace    = "TRACE"
	LabelDebug    = "DEBUG"
	LabelInfo     = "INFO"
	LabelNotice   = "NOTICE"
	LabelWarn     = "WARN"
	LabelError    = "ERROR"
	LabelCritical = "CRITICAL"
	LabelFatal    = "FATAL"
	LabelPanic    = "PANIC"
	LabelSilent   = "SILENT"
)

// Lvl -> Level
var parseLvl = map[string]slog.Level{
	LvlFlood:    LevelFlood,
	LvlTrace:    LevelTrace,
	LvlDebug:    LevelDebug,
	LvlInfo:     LevelInfo,
	LvlNotice:   LevelNotice,
	LvlWarn:     LevelWarn,
	LvlError:    LevelError,
	LvlCritical: LevelCritical,
	LvlFatal:    LevelFatal,
	LvlPanic:    LevelPanic,
	LvlSilent:   LevelSilent,
}

// Level -> Lvl
var parseLevel = map[slog.Level]string{
	LevelFlood:    LvlFlood,
	LevelTrace:    LvlTrace,
	LevelDebug:    LvlDebug,
	LevelInfo:     LvlInfo,
	LevelNotice:   LvlNotice,
	LevelWarn:     LvlWarn,
	LevelError:    LvlError,
	LevelCritical: LvlCritical,
	LevelFatal:    LvlFatal,
	LevelPanic:    LvlPanic,
	LevelSilent:   LvlSilent,
}

// Level() returns log level (Level() - implements slog.Leveler interface)
func (lp *Level) Level() slog.Level { return slog.Level(*lp) }

// Update log level
func (lp *Level) Update(level slog.Level) { *lp = Level(level) }

// String() returns a label for the level
func (lp *Level) String() string {
	str := func(base string, delta slog.Level) string {
		if delta == 0 {
			return base
		}
		return fmt.Sprintf("%s%+d", base, delta)
	}

	l := slog.Level(*lp)
	switch {
	case l < LevelTrace:
		return str(LabelFlood, l-LevelFlood)
	case l < LevelDebug:
		return str(LabelTrace, l-LevelTrace)
	case l < LevelInfo:
		return str(LabelDebug, l-LevelDebug)
	case l < LevelNotice:
		return str(LabelInfo, l-LevelInfo)
	case l < LevelWarn:
		return str(LabelNotice, l-LevelNotice)
	case l < LevelError:
		return str(LabelWarn, l-LevelWarn)
	case l < LevelCritical:
		return str(LabelError, l-LevelError)
	case l < LevelFatal:
		return str(LabelCritical, l-LevelCritical)
	case l < LevelPanic:
		return str(LabelFatal, l-LevelFatal)
	case l < LevelSilent:
		return str(LabelPanic, l-LevelPanic)
	default: // l >= LevelSilent
		return str(LabelSilent, l-LevelSilent)
	}
}

// Parse Lvl (string to num: Lvl -> Level)
func ParseLvl(lvl string) slog.Level {
	lvl = strings.ToLower(lvl)
	level, ok := parseLvl[lvl]
	if !ok { // try convert from numeric
		i, err := strconv.Atoi(lvl)
		if err != nil {
			return DEFAULT_LEVEL
		}
		level = slog.Level(i)
	}
	return level
}

// Parse Level (num to string: Level -> Lvl)
func ParseLevel(level slog.Level) string {
	lvl, ok := parseLevel[level]
	if !ok {
		return fmt.Sprintf("%d", int(level))
	}
	return lvl
}

// Return current log level as int (slog.Level)
func GetLevel() slog.Level {
	return currentXlog.GetLevel()
}

// Set current log level as int (slog.Level)
func SetLevel(level slog.Level) {
	currentXlog.SetLevel(level)
}

// Return current log level as string
func GetLvl() string {
	return currentXlog.GetLvl()
}

// Set current log level as string
func SetLvl(level string) {
	currentXlog.SetLvl(level)
}

// Internal wrapper to work with additional log levels
func logs(l *slog.Logger, level slog.Level, msg string, args ...any) {
	ctx := context.Background()
	if !l.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip wrappers
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)
	_ = l.Handler().Handle(ctx, r)
}

// Internal wrapper to work with additional log levels as standard logger
func logf(l *slog.Logger, level slog.Level, format string, args ...any) {
	ctx := context.Background()
	if !l.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip wrappers
	r := slog.NewRecord(time.Now(), level, fmt.Sprintf(format, args...), pcs[0])
	_ = l.Handler().Handle(ctx, r)
}

// EOF: "level.go"
