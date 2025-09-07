// File: "new.go"

package xlog

import (
	"io"
	"log"
	"log/slog" // go>=1.21
	"os"
	// FIXME: "golang.org/x/exp/slog" // экспериментальный пакет для go=1.20 только
)

// NewLog создает стандартный (legacy) логгер и настраивает его с учётом
// унифицированной структуры конфигурации Conf для X-logger
func NewLog(conf Conf) *log.Logger {
	logger := log.New(os.Stdout, "", 0)
	SetupLog(logger, conf)
	return logger
}

// New - создаёт новый X-logger на основе заданной структуры
// конфигурации Conf с возможностью вывода журнала в заданный канал
// (stdout/stderr) и/или с записью журнала в заданный файл с возможностью
// ротации (при необходимости), с формированием журнала в заданном формате
// (JSON/Text).
// Функция является простой надстройкой над NewWriter() и NewEx().
//
//	conf - обобщенная конфигурация логгера
//	mws - обёртки (middleware) для метода Hanlde() интерфейса slog.Handler
func New(conf Conf, mws ...Middleware) *Logger {
	w := NewWriter(conf.Pipe, conf.File, conf.FileMode, &conf.Rotate, nil)
	return NewEx(conf, w, mws...)
}

// NewWithWriter - создаёт новый X-logger на основе заданной структуры
// конфигурации Conf с возможностью вывода журнала в заданный канал
// (stdout/stderr) и/или с записью журнала в заданный файл с возможностью
// ротации (при необходимости), с формированием журнала в заданном формате
// (JSON/Text), а если указан writer, то с отправкой журналов
// в заданный io.Writer.
// Функция является простой надстройкой над NewWriter() и NewEx().
//
//	conf - обобщенная конфигурация логгера
//	writer - заданный писатель логов
//	mws - обёртки (middleware) для метода Hanlde() интерфейса slog.Handler
func NewWithWriter(conf Conf, writer io.Writer, mws ...Middleware) *Logger {
	w := NewWriter(conf.Pipe, conf.File, conf.FileMode, &conf.Rotate, writer)
	return NewEx(conf, w, mws...)
}

// NewEx - создаёт новый X-logger на основе заданной структуры
// конфигурации Conf с выдачей журнала через заданный Writer.
// Поля Pipe и File структуры conf при это игнорируются,
// направления вывода журнала определяются только заданным writer.
// Может быть задана произвольная цепочка Middleware.
//
//	conf - обобщенная конфигурация логгера
//	writer - заданный писатель логов (в т.ч. возможно с ротацией)
//	mws - обёртки (middleware) для метода Hanlde() интерфейса slog.Handler
func NewEx(conf Conf, writer Writer, mws ...Middleware) *Logger {
	var handler slog.Handler
	var level *slog.LevelVar

	fmt := logFormat(conf.Format)
	if fmt == logFmtDefault {
		// Для стандартного хендлера вывод возможен только на stdout
		writer = pipeWriter{os.Stdout}
	}

	// Создать slog.Handler и вернуть *slog.LeverVar
	handler, level = NewHandler(conf, writer, mws...)

	return &Logger{
		Logger: slog.New(handler),
		Level:  level,
		Writer: writer,
	}
}

// FromSlog создает X-logger из *slog.Logger.
// У результирующего "суррогатного" логгера нет возможности изменять
// уровень логирования, но можно использовать его "сахарные"
// методы типа Trace() или Noticef().
func FromSlog(logger *slog.Logger) *Logger {
	if logger == nil {
		return Default()
	}
	var level slog.LevelVar
	return &Logger{
		Logger: logger,
		Level:  &level,
		Writer: nullWriter{},
	}
}

// EOF: "setup.go"
