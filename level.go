// File: "level.go"

package xlog

import (
	"context"
	"fmt"
	"log/slog" // go>=1.21
	"runtime"
	"time"
	//"golang.org/x/exp/slog" // depricated for go>=1.21
)

// xlog level delivered from slog.Level, implements slog.Leveler interface
type Level slog.Level

// Log levels delivered from slog.Level
const (
	LevelTrace  = Level(slog.Level(-8))  // TRACE  (-8)
	LevelDebug  = Level(slog.LevelDebug) // DEBUG  (-4)
	LevelInfo   = Level(slog.LevelInfo)  // INFO   (0)
	LevelNotice = Level(slog.Level(2))   // NOTICE (2)
	LevelWarn   = Level(slog.LevelWarn)  // WARN   (4)
	LevelError  = Level(slog.LevelError) // ERROR  (8)
	LevelFatal  = Level(slog.Level(12))  // FATAL  (12)
	LevelPanic  = Level(slog.Level(16))  // PANIC  (16)
	LevelSilent = Level(slog.Level(20))  // SILENT (20)
)

const DEFAULT_LEVEL = LevelInfo

// Log level as string for setup
type Lvl string

const (
	LvlTrace  Lvl = "trace"
	LvlDebug      = "debug"
	LvlInfo       = "info"
	LvlNotice     = "notice"
	LvlWarn       = "warn"
	LvlError      = "error"
	LvlFatal      = "fatal"
	LvlPanic      = "panic"
	LvlSilent     = "silent"
)

// Log level tags
const (
	LabelTrace  = "TRACE"
	LabelDebug  = "DEBUG"
	LabelInfo   = "INFO"
	LabelNotice = "NOTICE"
	LabelWarn   = "WARN"
	LabelError  = "ERROR"
	LabelFatal  = "FATAL"
	LabelPanic  = "PANIC"
	LabelSilent = "SILENT"
)

// Lvl -> Level
var parseLvl = map[Lvl]Level{
	LvlTrace:  LevelTrace,
	LvlDebug:  LevelDebug,
	LvlInfo:   LevelInfo,
	LvlNotice: LevelNotice,
	LvlWarn:   LevelWarn,
	LvlError:  LevelError,
	LvlFatal:  LevelFatal,
	LvlPanic:  LevelPanic,
	LvlSilent: LevelSilent,
}

// Level -> Lvl
var parseLevel = map[Level]Lvl{
	LevelTrace:  LvlTrace,
	LevelDebug:  LvlDebug,
	LevelInfo:   LvlInfo,
	LevelNotice: LvlNotice,
	LevelWarn:   LvlWarn,
	LevelError:  LvlError,
	LevelFatal:  LvlFatal,
	LevelPanic:  LvlPanic,
	LevelSilent: LvlSilent,
}

// Level() returns the receiver (it implements slog.Leveler interface)
func (l Level) Level() slog.Level { return slog.Level(l) }

// String() returns a label for the level
func (l Level) String() string {
	str := func(base string, delta Level) string {
		if delta == 0 {
			return base
		}
		return fmt.Sprintf("%s%+d", base, delta)
	}

	switch {
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
	case l < LevelFatal:
		return str(LabelError, l-LevelError)
	case l < LevelPanic:
		return str(LabelFatal, l-LevelFatal)
	case l < LevelSilent:
		return str(LabelPanic, l-LevelPanic)
	default: // l >= LevelSilent
		return str(LabelSilent, l-LevelSilent)
	}
}

// Parse Lvl (Lvl -> Level)
func ParseLvl(lvl Lvl) Level {
	level, ok := parseLvl[lvl]
	if !ok {
		return DEFAULT_LEVEL
	}
	return level
}

// Parse Level (Level -> Lvl)
func ParseLevel(level Level) Lvl {
	lvl, ok := parseLevel[level]
	if !ok {
		return Lvl(fmt.Sprintf("%d", int(level)))
	}
	return lvl
}

// Return current log level as int (slog.Level)
func GetLevel() Level {
	defaultLock.Lock()
	defer defaultLock.Unlock()
	return currentLevel
}

// Return current log level as string (xlog.Lvl)
func GetLvl() Lvl {
	level := GetLevel()
	return ParseLevel(level)
}

// Internal wrapper to work with additional log levels
func logs(l *slog.Logger, level Level, msg string, args ...any) {
	if l == nil {
		l = slog.Default()
	}
	ctx := context.Background()
	if !l.Enabled(ctx, slog.Level(level)) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip wrappers
	r := slog.NewRecord(time.Now(), slog.Level(level), msg, pcs[0])
	r.Add(args...)
	_ = l.Handler().Handle(ctx, r)
}

// Internal wrapper to work with additional log levels as standart logger
func logf(l *slog.Logger, level Level, format string, args ...any) {
	if l == nil {
		l = slog.Default()
	}
	ctx := context.Background()
	if !l.Enabled(ctx, slog.Level(level)) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip wrappers
	r := slog.NewRecord(time.Now(), slog.Level(level),
		fmt.Sprintf(format, args...), pcs[0])
	_ = l.Handler().Handle(ctx, r)
}

// EOF: "level.go"
