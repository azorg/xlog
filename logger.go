// File: "logger.go"

package xlog

import (
	"log/slog" // go>=1.21
	// FIXME: "golang.org/x/exp/slog" // экспериментальный пакет для go=1.20 только
)

// Logger описывает полную структуру управления логгером.
type Logger struct {
	// Logger - указатель на стандартный slog-логгер.
	// Допустимо (и рекомендуется) использовать сигнатуры стандартного
	// пакета "log/slog" при организации процесса логгирования в приложениях,
	// потому именно данный указатель является основной целью функций (фабрик)
	// инициализирующих структуру Logger.
	*slog.Logger

	// Level - структура для безопасного управления уровнем логирования.
	// Для изменения уровня логирования в процессе
	// выполнения программы может использоваться её метод Set().
	// Метод Level() возвращает текущий уровень логирования.
	Level *slog.LevelVar

	// Writer - это обобщенный интерфейс писателя логов.
	// Данный интерфейс имеет метод Rotate() который может
	// использоваться для принудительной ротации логов в приложении,
	// например, по сигналу SIGHUP.
	// Метод IsRotatable() возвращает признак возможности ротации файла журнала.
	// Интерфейс содержит так же стандартный метод Write(),
	// который может использоваться для прямой записи в канал и/или
	// в файл журнала.
	// Пакет xlog имеет фабрику для данных интерфейсов NewWriter().
	Writer
}

// Slog возвращает текущий указатель *slog.Logger
func (c *Logger) Slog() *slog.Logger { return c.Logger }

// Slog возвращает указатель *slog.Logger из текущего
// (глобального) логгера
func Slog() *slog.Logger { return currentClog.Logger }

// GetLevel возвращает текущий уровень логирования -
// обёртка для вызова c.Leveler.Level()
func (c *Logger) GetLevel() slog.Level { return c.Level.Level() }

// GetLevel возвращает текущий уровень логирования для глобального логгера
func GetLevel() slog.Level { return currentClog.GetLevel() }

// SetLevel обновляет уровень логирования -
// обёртка для вызова c.Leveler.Update()
func (c *Logger) SetLevel(level slog.Level) { c.Level.Set(level) }

// SetLevel обновляет уровень логирования для глобального логгера
func SetLevel(level slog.Level) { currentClog.SetLevel(level) }

// GetLvl возвращает текущий уровень логирования в виде строки
// вида "info", "debug" и т.п.
func (c *Logger) GetLvl() string { return LevelToString(c.GetLevel()) }

// GetLvl возвращает текущий уровень логирования глобального логгера в виде
// строки вида "info", "debug" и т.п.
func GetLvl() string { return currentClog.GetLvl() }

// SetLvl обновляет уровень логирования на основе строки идентификатора
// типа "trace", "error" и т.п.
func (c *Logger) SetLvl(level string) { c.SetLevel(LevelFromString(level)) }

// SetLvl обновляет уровень логирования глобального логгера на
// основе строки идентификатора типа "trace", "error" и т.п.
func SetLvl(level string) { currentClog.SetLvl(level) }

// Rotate производит ротацию файла журнала (если это возможно) -
// обёртка для вызова c.Writer.Rotable()/c.Writer.Rotate().
// К примеру, в реальных приложениях возможна организации ротация
// логов по сигналу SIGHUP.
func (c *Logger) Rotate() error {
	if !c.Writer.IsRotatable() {
		return ErrNotRotatable // ротация не предусмотрена конфигурацией
	}
	return c.Writer.Rotate()
}

// Rotate производит ротацию файла журнала (если это возможно)
// для глобального логгера.
// К примеру, в реальных приложениях возможна организации ротация
// логов по сигналу SIGHUP.
func Rotate() error {
	return currentClog.Rotate()
}

// IsRotatable возвращает признак возможности ротации файла журнала
func (c *Logger) IsRotatable() bool { return c.Writer.IsRotatable() }

// IsRotatable возвращает признак возможности ротации файла журнала
// для лобального логгера
func IsRotatable() bool { return currentClog.IsRotatable() }

// SlogWithFields создает дочерний *slog.Logger с добавлением заданных атрибутов
//
//	log - исходной slog логгер
//	fields - интерфейс для получения дополнительных атрибутов
func SlogWithFields(log *slog.Logger, fields ...FieldsProvider) *slog.Logger {
	if len(fields) == 0 {
		return log
	}
	as := make([]slog.Attr, 0)
	for i := range fields {
		as = append(as, fields[i].Fields().Attrs()...)
	}
	handler := log.Handler().WithAttrs(as)
	return slog.New(handler)
}

// With создает дочерний логгер с добавлением заданных атрибутов
func (c *Logger) WithFields(fields ...FieldsProvider) *Logger {
	return &Logger{
		Logger: SlogWithFields(c.Logger, fields...),
		Level:  c.Level,
		Writer: c.Writer,
	}
}

// WithFields создает дочерний логгер с добавлением заданных атрибутов
// на основе глобального логгера.
func WithFields(fields FieldsProvider) *Logger {
	return currentClog.WithFields(fields)
}

// WithAttrs создает дочерний логгер c добавлением заданных атрибутов
func (c *Logger) WithAttrs(attrs []slog.Attr) *Logger {
	return &Logger{
		Logger: slog.New(c.Logger.Handler().WithAttrs(attrs)),
		Level:  c.Level,
		Writer: c.Writer,
	}
}

// WithAttrs создает дочерний логгер c добавлением заданных атрибутов
// на основе глобального логгера
func WithAttrs(attrs []slog.Attr) *Logger {
	return currentClog.WithAttrs(attrs)
}

// With создает дочерний логгер с добавлением заданных атрибутов.
// Метод аналогичен методу With для *slog.Logger.
func (c *Logger) With(args ...any) *Logger {
	return &Logger{
		Logger: c.Logger.With(args...),
		Level:  c.Level,
		Writer: c.Writer,
	}
}

// With создает дочерний логгер с добавлением заданых атрибутов
// на основе глобального логгера
func With(args ...any) *Logger {
	return currentClog.With(args...)
}

// WithGroup создает дочерний логгер с добавлением группы.
// Метод аналогичен методу WithGroup для *slog.Logger.
func (c *Logger) WithGroup(name string) *Logger {
	return &Logger{
		Logger: c.Logger.WithGroup(name),
		Level:  c.Level,
		Writer: c.Writer,
	}
}

// WithGroup создает дочерний логгер с добавлением группы
// на основе глобального логгера
func WithGroup(name string) *Logger {
	return currentClog.WithGroup(name)
}

// WithMiddleware создает дочерний логгер c добавлением Middleware
func (c *Logger) WithMiddleware(mws ...Middleware) *Logger {
	if len(mws) == 0 {
		return c
	}
	handler := NewMiddlewareHandler(c.Handler(), mws...)
	return &Logger{
		Logger: slog.New(handler),
		Level:  c.Level,
		Writer: c.Writer,
	}
}

// WithMiddleware создает дочерний логгер c добавлением Middleware
// на основе глобального логгера
func WithMiddleware(mws ...Middleware) *Logger {
	return currentClog.WithMiddleware(mws...)
}

// EOF: "logger.go"
