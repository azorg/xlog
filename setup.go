// File: "setup.go"

package xlog

import (
	"context"
	//"io"
	"log"
	"log/slog" // go>=1.21
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
	//"golang.org/x/exp/slog" // deprecated for go>=1.21
)

// Saved loggers, current log level
var (
	defaultLog  *log.Logger  = log.Default()  // initial standard logger
	defaultSlog *slog.Logger = slog.Default() // initial structured logger
	currentXlog *Logger      = Default()      // current global logger
	defaultLock sync.Mutex
)

// Setup standard simple logger
func SetupLog(logger *log.Logger, conf Conf) {
	flag := 0
	if conf.Time {
		flag |= log.LstdFlags
		if conf.TimeUS {
			flag |= log.Lmicroseconds
		}
	}
	if conf.Src {
		if conf.SrcLong {
			flag |= log.Llongfile
		} else {
			flag |= log.Lshortfile
		}
	}
	if conf.Prefix != "" {
		flag |= log.Lmsgprefix
	}

	rotator := newRotator(conf.File, conf.FileMode, &conf.Rotate)
	logger.SetOutput(rotator.Writer())
	logger.SetPrefix(conf.Prefix)
	logger.SetFlags(flag)
}

// Create new configured standard logger
func NewLog(conf Conf) *log.Logger {
	logger := log.New(os.Stdout, "", 0)
	SetupLog(logger, conf)
	return logger
}

// Create new configured structured logger (default/text/JSON/Tinted handler)
func New(conf Conf) *Logger {
	if !conf.Slog && !conf.JSON && !conf.Tint {
		// Don't use Text/JSON/Tint handler, tune standard logger
		return newSlogStd(conf)
	}

	level := Level(ParseLvl(conf.Level))
	rotator := newRotator(conf.File, conf.FileMode, &conf.Rotate)
	var handler slog.Handler

	if conf.Tint {
		// Use Tinted Handler
		opts := &TintOptions{
			Level:       &level,
			AddSource:   conf.Src,
			SourceLong:  conf.SrcLong,
			SourceFunc:  conf.SrcFunc,
			NoExt:       conf.NoExt,
			NoLevel:     conf.NoLevel,
			TimeFormat:  conf.TimeTint,
			NoColor:     conf.NoColor,
			ReplaceAttr: nil,
		}

		if conf.TimeTint == "" {
			if conf.Time {
				opts.TimeFormat = DEFAULT_TIME_FORMAT
				if conf.TimeUS {
					opts.TimeFormat = DEFAULT_TIME_FORMAT_US
				}
			}
		}

		handler = NewTintHandler(rotator.Writer(), opts)
	} else {
		// Use Text/JSON handler
		opts := &slog.HandlerOptions{
			AddSource: conf.Src,
			Level:     &level,
		}

		if conf.JSON {
			conf.SrcFunc = false // force off
		}

		if ADD_LEVELS || !conf.Time || conf.Src || conf.NoExt || conf.NoLevel {
			opts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
				switch a.Key {
				case slog.TimeKey:
					if !conf.Time && len(groups) == 0 {
						return slog.Attr{} // remove timestamp from log
					}
					if conf.TimeUS {
						t := a.Value.Any().(time.Time)
						tstr := t.Format(RFC3339Micro)
						return slog.String(slog.TimeKey, tstr)
					} else if conf.JSON { // FIX difference between Text and JSON handler
						t := a.Value.Any().(time.Time)
						tstr := t.Format(RFC3339Milli)
						return slog.String(slog.TimeKey, tstr)
					}
				case slog.SourceKey:
					src := a.Value.Any().(*slog.Source)
					if src.File == "" { // FIX some bug if slog work as standard logger
						return slog.Attr{}
					}
					if conf.SrcLong { // long: directory + file name
						dir, file := filepath.Split(src.File)
						src.File = filepath.Join(filepath.Base(dir), file)
					} else { // short: only file name
						src.File = path.Base(src.File)
					}
					if conf.NoExt { // remove ".go" extension
						src.File = RemoveGoExt(src.File)
					}
					parts := strings.Split(src.Function, ".")
					if len(parts) != 0 {
						src.Function = parts[len(parts)-1]
					}
					//src.Funcion = GetFuncName(7) // skip=7 (some magic)
					if conf.SrcFunc { // add function name
						if src.Function != "" {
							src.File += ":" + src.Function + "()"
						}
					}
					a.Value = slog.AnyValue(src)
				case slog.LevelKey:
					if conf.NoLevel {
						return slog.Attr{} // remove "level=..." etc
					}
					if ADD_LEVELS { // add TRACE/NOTICE/FATAL/PANIC
						level := a.Value.Any().(slog.Level)
						leveler := Level(level)
						label := leveler.String()
						return slog.String(slog.LevelKey, label)
					}
				} // switch
				return a
			}
		}

		if conf.JSON {
			handler = slog.NewJSONHandler(rotator.Writer(), opts)
		} else {
			handler = slog.NewTextHandler(rotator.Writer(), opts)
		}
	}

	logger := slog.New(handler)
	if conf.AddKey != "" && conf.AddValue != "" {
		logger = logger.With(conf.AddKey, conf.AddValue)
	}

	return &Logger{
		Logger:  logger,
		Level:   &level,
		Rotator: rotator,
	}
}

// Create new configured structured logger (default/text/JSON/Tinted handler)
func NewSlog(conf Conf) *slog.Logger {
	logger := New(conf)
	return logger.Logger
}

// Setup standart and structured default global loggers
func Setup(conf Conf) {
	// Setup standart logger
	l := logDefault()
	SetupLog(l, conf)

	// Setup structured logger
	logger := New(conf)
	slog.SetDefault(logger.Logger)

	// Save log level and update global xlog wrapper
	defaultLock.Lock()
	currentXlog = logger
	defaultLock.Unlock()

	// Repeat setup standart logger (stop loop forever)
	// TODO: why?
	SetupLog(l, conf)

	// FIXME: TODO
	// В экспериментальном slog есть ошибка:
	// При добавлении хендлера для управления уровнем
	// логирования некорректно выводятся имена файлов (и строк).
	// Что интересно, в Go v1.21 в log/slog всё исправлено.
	// Используя Go до версии 1.21 (например 1.20) при включении
	// управления уровнем логирования при работе slog через slog.defaultHandler
	// в угоду возможности управления уровнями отключаем вывод файлов и строк.
	// Если "golang.org/x/exp/slog" доработают это FIX можно будет убрать.
	if OLD_SLOG_FIX { // "runtime.Version() < go1.21.0"
		if currentXlog.Level.Level() != DEFAULT_LEVEL {
			stdlog := logDefault()
			flag := stdlog.Flags()
			flag = flag &^ (log.Lshortfile | log.Llongfile) // sorry...
			stdlog.SetFlags(flag)
		}
	}
}

// Get default standart logger
func logDefault() *log.Logger {
	defaultLock.Lock()
	defer defaultLock.Unlock()
	l := defaultLog
	if l == nil {
		l = log.Default()
		defaultLog = l
	}
	return l
}

// Get default structured logger
func slogDefault() *slog.Logger {
	defaultLock.Lock()
	defer defaultLock.Unlock()
	lg := defaultSlog
	if lg == nil {
		lg = slog.Default()
		defaultSlog = lg
	}
	return lg
}

// Create custom structured logger based on default standard logger
func newSlogStd(conf Conf) *Logger {
	// Setup standard logger
	stdlog := logDefault()
	SetupLog(stdlog, conf)

	level := ParseLvl(conf.Level)
	leveler := Level(level)
	var logger *slog.Logger
	if level == DEFAULT_LEVEL {
		// Don't change default log level
		logger = slogDefault()
	} else {
		// Hook to direct log level
		handler := slogDefault().Handler() // slog.defaultHandler
		handler = newStdHandler(handler, &leveler)
		logger = slog.New(handler)
	}

	if conf.AddKey != "" && conf.AddValue != "" {
		logger = logger.With(conf.AddKey, conf.AddValue)
	}

	return &Logger{
		Logger:  logger,
		Level:   &leveler,
		Rotator: newPipe(os.Stdout),
	}
}

// Help wrapper to direct log level in standard logger mode
type stdHandler struct {
	handler slog.Handler
	level   slog.Leveler
}

// Create logStdHandler with the given level
func newStdHandler(handler slog.Handler, level slog.Leveler) *stdHandler {
	// Optimization: avoid chains of logStdHandlers
	if sh, ok := handler.(*stdHandler); ok {
		handler = sh.handler
	}
	return &stdHandler{handler: handler, level: level}
}

// Enabled() implements Enabled() by reporting whether
func (h *stdHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

// Handle() implements Handler.Handle()
func (h *stdHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.handler.Handle(ctx, r)
}

// WithAttrs() implements Handler.WithAttrs()
func (h *stdHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	//if len(attrs) == 0 {
	//  return h
	//}
	return newStdHandler(h.handler.WithAttrs(attrs), h.level)
}

// WithGroup() implements Handler.WithGroup()
func (h *stdHandler) WithGroup(name string) slog.Handler {
	//if name == "" {
	//  return h
	//}
	return newStdHandler(h.handler.WithGroup(name), h.level)
}

// EOF: "setup.go"
