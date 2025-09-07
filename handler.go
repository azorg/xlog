// File: "handler.go"

package xlog

import (
	"fmt"
	"io"
	"log/slog" // go>=1.21
	"path"
	"path/filepath"
	"time"
	// FIXME: "golang.org/x/exp/slog" // экспериментальный пакет для go=1.20 только
)

// NewHandler создаёт новый *slog.Handler на основе заданной структуры
// конфигурации conf с выдачей журнала через заданный writer.
// Возвращаемый хендлер в соответствии с конфигурацией будет формировать
// требуемые дополнительные атрибуты (goroutine, logId, logSum).
// Заодно возвращается указатель на slog.LevelVar для возможности
// безопасного управления уровнем логирования в будущем.
//
//	conf - параметры конфигурации логгера
//	writer - писатель журнала
//	mws - обёртки для метода Hanlde() интерфейса slog.Handler
func NewHandler(conf Conf, writer io.Writer, mws ...Middleware) (
	handler slog.Handler, _ *slog.LevelVar) {

	format := logFormat(conf.Format)
	if format == logFmtDefault {
		// Не использовать Text/JSON/Tint handler,
		// немного подстроить стандартный slog-handler по умолчанию
		// для реализации возможности управления уровнем логирования,
		// обогатить перечень атрибутов по умолчанию.
		// При этом нельзя управлять выводом (только stdout).
		return NewStdHandler(conf, mws...)
	}

	var level slog.LevelVar
	level.Set(LevelFromString(conf.Level))

	if format == logFmtTint { // использовать TintHandler
		// Выбрать формат временной метки
		timeFormat := ""
		if conf.TimeFormat == "" {
			if !conf.TimeOff {
				timeFormat = defaultTime
				if conf.TimeMicro {
					timeFormat = defaultTimeMicro
				}
			}
		} else {
			timeFormat, _ = TimeFormat(conf.TimeFormat)
			conf.TimeOff = false // сохранить для настройки idHandler'а
		}

		// Использовать Tinted Handler
		opts := &TintOptions{
			Level:       &level, // slog.Leveler
			AddSource:   conf.Src,
			SourcePkg:   conf.SrcPkg,
			SourceFunc:  conf.SrcFunc,
			NoExt:       !conf.SrcExt,
			NoLevel:     conf.LevelOff,
			TimeUTC:     !conf.TimeLocal,
			TimeFormat:  timeFormat,
			NoColor:     conf.ColorOff,
			ReplaceAttr: nil,
		}

		handler = NewTintHandler(writer, opts)
	} else { // использовать стандартный slog Text/JSON handler
		opts := &slog.HandlerOptions{
			AddSource: conf.Src,
			Level:     &level, // slog.Leveler
		}

		if format == logFmtJSON {
			conf.SrcFunc = false // в JSON режиме не требуется
		}

		opts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
			if format == logFmtJSON {
				// Выдать значение complex128 в JSON журнал как строку
				cval, ok := a.Value.Any().(complex128)
				if ok {
					return slog.String(a.Key, fmt.Sprintf("%v", cval))
				}
			}

			if addLevels { // заменить метод String() типа slog.Level
				level, ok := a.Value.Any().(slog.Level)
				if ok {
					a.Value = slog.StringValue(LevelToLabel(level))
				}
			}

			if len(groups) != 0 { // модифицируем только корневые атрибуты
				return a // вернуть атрибут как есть
			}

			switch a.Key {
			case slog.TimeKey:
				if conf.TimeOff { // удалить временную метку
					return slog.Attr{}
				}

				t, ok := a.Value.Any().(time.Time)
				if !ok {
					return a // вернуть атрибут как есть (это не время)
				}

				if !conf.TimeLocal { // вывести метку времени в UTC
					t = t.UTC()
					var tstr string
					if conf.TimeMicro {
						tstr = t.Format(RFC3339Micro)
					} else {
						tstr = t.Format(RFC3339Milli)
					}
					return slog.String(slog.TimeKey, tstr)
				}

				// Тут небольшая "магия" в связи с тем, что slog.TextHandler и
				// slog.JSONHandler реализованы несколько по разному:
				//  - TextHandler по умолчанию формирует метку времени с миллисекундами;
				//  - JSONHandler по умолчанию формирует метку времени БЕЗ долей секунд.
				// Код ниже "выравнивает" этот перекос.
				if conf.TimeMicro {
					return slog.String(slog.TimeKey, t.Format(RFC3339Micro))
				} else if format == logFmtJSON {
					return slog.String(slog.TimeKey, t.Format(RFC3339Milli))
				}

			case slog.SourceKey:
				src, ok := a.Value.Any().(*slog.Source)
				if ok {
					if src.File == "" { // FIX some bug if slog work as standard logger
						return slog.Attr{}
					}
					if conf.SrcPkg { // directory + file name
						dir, file := filepath.Split(src.File)
						src.File = filepath.Join(filepath.Base(dir), file)
					} else { // only file name
						src.File = path.Base(src.File)
					}
					if !conf.SrcExt { // remove ".go" extension
						src.File = removeGoExt(src.File)
					}
					//src.Function = getFuncName(7) // skip=7 (some magic)
					src.Function = cropFuncName(src.Function)
					if conf.SrcFunc { // add function name (not for JSON)
						if src.Function != "" {
							src.File += ":" + src.Function + "()"
						}
					}

					if format == logFmtJSON && // только для JSON
						conf.SrcFields != nil && len(*conf.SrcFields) != 0 {
						// Обогатить блок "source" заменив slog.Source на Fields
						srcFields := Fields{
							"function": src.Function,
							"file":     src.File,
							"line":     src.Line,
						}
						for k, v := range *conf.SrcFields {
							if v != nil { // пропустить <nil> и пустые строки
								if str, ok := v.(string); !ok || str != "" {
									srcFields[k] = v
								}
							}
						} // for
						a.Value = slog.AnyValue(srcFields)
					} else {
						a.Value = slog.AnyValue(src)
					}
				}

			case slog.LevelKey:
				if conf.LevelOff { // удалить метку уровня
					return slog.Attr{}
				}
			} // switch

			return a
		} // opt.ReplaceAttr

		if format == logFmtText { // logfmt
			handler = slog.NewTextHandler(writer, opts)
		} else { // JSON
			handler = slog.NewJSONHandler(writer, opts)
		}
	}

	if format != logFmtJSON && conf.Src &&
		conf.SrcFields != nil /*&& len(conf.SrcFields.Fields()) != 0*/ {
		// Для текстовых форматов обогатить вывод conf.SrcFields
		mw := NewMiddlewareWithFields(conf.SrcFields)
		ms := make([]Middleware, 0, len(mws)+1)
		ms = append(ms, mw)
		mws = append(ms, mws...)
	}

	// Использовать IdHandler безусловно (для полноценной работы slog.LogValuer'ов)
	if true || conf.GoId || conf.IdOn || conf.SumOn || len(mws) != 0 {
		// Создать дополнительный хендлер-обёртку
		// для добавления в журнал goroutine, logId, logSum
		// и с заданными middleware(s)
		idOpts := &IdOptions{
			GoId:     conf.GoId,
			LogId:    conf.IdOn,
			AddSum:   conf.SumOn,
			SumFull:  conf.SumFull,
			SumTime:  !conf.TimeOff,
			SumChain: conf.SumChain,
			SumAlone: conf.SumAlone,
		}
		handler = NewIdHandler(handler, idOpts, 0x00, mws...) // sum=0x00
	}

	if conf.AddKey != "" && conf.AddValue != nil {
		// Обогатить вывод хендлера заданным дополнительным Key/Value
		attr := slog.Any(conf.AddKey, slog.AnyValue(conf.AddValue))
		handler = handler.WithAttrs([]slog.Attr{attr})
	}

	return handler, &level
}

// Создать модернизированный хендлер стандартного slog логгера по умолчанию
// с возможностью управления уровнем логирования и с формированием
// дополнительных атрибутов (goroutine, logId, logSum) и middleware
//
//	conf - параметры конфигурации логгера
//	mws - обёртки для метода Hanlde() интерфейса slog.Handler
func NewStdHandler(conf Conf, mws ...Middleware) (slog.Handler, *slog.LevelVar) {
	// Настроить стандартный (legacy) логгер
	SetupLog(defaultLog, conf)

	var level slog.LevelVar
	level.Set(LevelFromString(conf.Level))

	handler := defaultSlog.Handler() // slog.defaultHandler

	// Создать дополнительный хендлер-обёртку для того, чтобы управлять
	// уровнем логирования стандартного логгера slog
	handler = newLvlHandler(handler, &level)

	if conf.Src &&
		conf.SrcFields != nil && len(*conf.SrcFields) != 0 {
		// Обогатить вывод данных conf.SrcFields
		mw := NewMiddlewareWithFields(conf.SrcFields)
		ms := make([]Middleware, 0, len(mws)+1)
		ms = append(ms, mw)
		mws = append(ms, mws...)
	}

	if conf.GoId || conf.IdOn || conf.SumOn || len(mws) != 0 {
		// Создать дополнительный хендлер-обёртку
		// для добавления в журнал goroutine, logId, logSum
		// и с заданными middleware(s)
		idOpts := &IdOptions{
			GoId:     conf.GoId,
			LogId:    conf.IdOn,
			AddSum:   conf.SumOn,
			SumFull:  conf.SumFull,
			SumTime:  !conf.TimeOff,
			SumChain: conf.SumChain,
			SumAlone: conf.SumAlone,
		}
		handler = NewIdHandler(handler, idOpts, 0x00, mws...) // sum=0x00
	}

	if conf.AddKey != "" && conf.AddValue != nil {
		// Обогатить вывод хендлера заданным дополнительным Key/Value
		attr := slog.Any(conf.AddKey, slog.AnyValue(conf.AddValue))
		handler = handler.WithAttrs([]slog.Attr{attr})
	}

	return handler, &level
}

// EOF: "handler.go"
