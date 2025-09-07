// File: "writer.go"

package xlog

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2" // ротатор файлов журналов
)

// Обобщенный интерфейс писателя логов (с ротацией или без, с дублированием журнала
// на stdout/stderr или без, с выдачей в кастомный io.Writer или без).
type Writer interface {
	IsRotatable() bool // проверить возможность ротации
	Rotate() error     // выполнить ротацию журнала (если она возможна)
	Close() error      // закрыть файл журнала
	io.Writer          // интерфейс со стандартным методом Wite([]byte) (int, error)
}

// Писатель логов в никуда (для явного отключения журналирования)
type nullWriter struct{}

// Писатель логов в заданный канал (os.Stdout/os.Stderr)
type pipeWriter struct{ *os.File }

// Писатель логов в заданный файл
type fileWriter struct{ *os.File }

// Писатель логов в заданный файл с ротацией на основе *lumberjack.Logger
type rotatableWriter struct{ *lumberjack.Logger }

// Писатель логов в заданный пользователем io.Writer
type customWriter struct{ io.Writer }

// Писатель логов в заданный файл с дублированием
// журнала в заданный канал (stdout/stderr)
type pipeAndFileWriter struct {
	Pipe     *os.File // os.Stdout/os.Stderr
	*os.File          // обычный файл журанала
}

// Писатель логов в заданный файл с ротацией на основе *lumberjack.Logger с
// дублированием журнала в заданный канал (stdout/stderr)
type pipeAndRotatableWriter struct {
	Pipe               *os.File // os.Stdout/os.Stderr
	*lumberjack.Logger          // ротатор лога
}

// Писатель логов в заданный io.Writer c дублированием
// журнала в заданный канал (stdout/stderr)
type pipeAndCustomWriter struct {
	Pipe      *os.File // os.Stdout/os.Stderr
	io.Writer          // заданный io.Writer
}

// Писатель логов в файл и в заданный io.Writer
type fileAndCustomWriter struct {
	*os.File  // обычный файл журанала
	io.Writer // заданный io.Writer
}

// Писатель логов в заданный файл с ротацией на основе *lumberjack.Logger с
// дублированием журнала в заданный io.Writer
type rotatableAndCustomWriter struct {
	*lumberjack.Logger // ротатор лога
	io.Writer          // заданный io.Writer
}

// Писатель логов в файл, в заданный io.Writer с
// дублированием журнала в заданный канал (stdout/stderr)
type pipeAndFileAndCustomWriter struct {
	Pipe      *os.File // os.Stdout/os.Stderr
	*os.File           // обычный файл журанала
	io.Writer          // заданный io.Writer
}

// Писатель логов в заданный файл с ротацией на основе *lumberjack.Logger с
// дублированием журнала в заданный канал (stdout/stderr) и
// в заданный io.Writer
type pipeAndRotatableAndCustomWriter struct {
	Pipe               *os.File // os.Stdout/os.Stderr
	*lumberjack.Logger          // ротатор лога
	io.Writer                   // заданный io.Writer
}

// Убедится в том, что все писатели логов соответствуют интерфейсу Writer
var _ Writer = nullWriter{}
var _ Writer = pipeWriter{}
var _ Writer = fileWriter{}
var _ Writer = rotatableWriter{}
var _ Writer = customWriter{}
var _ Writer = pipeAndFileWriter{}
var _ Writer = pipeAndRotatableWriter{}
var _ Writer = pipeAndCustomWriter{}
var _ Writer = fileAndCustomWriter{}
var _ Writer = rotatableAndCustomWriter{}
var _ Writer = pipeAndFileAndCustomWriter{}
var _ Writer = pipeAndRotatableAndCustomWriter{}

// getPipe производит выбор канала (os.Stdout, os.Stderr или nil).
// Строка pipeName может принимать следующие значения (регистр
// не имеет значение):
//
//	stdout, os.Stdout, standart - для возврата os.Stdout
//	stderr, os.Stderr, error - для возврата os.Stderr
//	nil, пустая строка, null, none, off - для возврата nil
//
// Если noFile=true, то при пустой строке pipeName возвращается os.Stdout,
// а не nil.
func getPipe(pipeName string, noFile bool) *os.File {
	if pipeName == "" && noFile {
		return os.Stdout // по умолчанию (zero configuration) использовать stdout
	}
	switch strings.ToLower(pipeName) {
	case "stdout", "os.stdout", "standard":
		return os.Stdout
	case "stderr", "os.stderr", "error":
		return os.Stderr
	default: // ~ "", "null", "nil", "none", "off"
		return nil
	}
}

// NewWriter создаёт писатель логов с интерфейсом Writer на основе заданных
// параметров логирования.
//
//	pipeName - имя канала ("stdout", 'stderr", "nil" или пустая строка)
//	fileName - имя файла журнала или пустая строка
//	mode - режим доступа к файлу или пустая строка (например "0644" или "0660")
//	rotate - параметры ротации файла журнала или nil
//	writer - кастомный io.Writer или nil
//
// Возвращаемый Writer осуществляет вывод логов по нескольким направлениям
// в зависимости от заданных параметров (от нуля до трех направлений).
// Возможны следующие направления:
//
//  1. Канал (pipe): stdout или stderr
//  2. Файл журнала (file) с ротацией или без
//  3. Кастомный io.Writer
//
// Данная фабрика может возвращать 12 типов различных реализаций интерфейса
// Writer: от писателя в "/dev/null" до писателя в три направления.
func NewWriter(
	pipeName, fileName, mode string, rotate *RotateConf, writer io.Writer,
) Writer {

	var perm fs.FileMode
	if fileName != "" {
		// Распаковать права доступа (в восьмеричной Unix нотации)
		perm = fileMode(mode)

		// Создать (при необходимости) каталог для файлов журналов
		dir := filepath.Dir(fileName)
		if dir != "" {
			// FIXME: немного магии - права доступа на каталог
			// определяются из прав доступа к файлам
			dirPerm := perm | ((perm & 0044) >> 2) | 0700
			err := os.MkdirAll(dir, dirPerm)
			if err != nil { // в случае ошибки использовать stdout
				fmt.Fprintf(os.Stderr,
					"ERROR: can't create logfile directory: %v\n", err)
				fileName = ""
			}
		}
	}

	var file *os.File
	if fileName != "" {
		// Открыть (создать) файл журнала
		var err error
		file, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, perm)
		if err != nil { // в случае ошибки использовать stdout
			fmt.Fprintf(os.Stderr,
				"ERROR: can't open/create logfile: %v\n", err)
			file = nil
		}
	}

	pipe := getPipe(pipeName, file == nil) // os.Stdout, os.Stderr or nil

	var logger *lumberjack.Logger
	if file != nil && rotate != nil && rotate.Enable { // использовать ротацию
		maxSize := rotate.MaxSize
		if maxSize == 0 { // изменить умолчание для MaxSize от lumberjack
			maxSize = RotateMaxSize
		}

		file.Close() // закрыть файл, т.к. lumberjack откроет его сам
		file = nil

		// Создать объект ротации
		logger = &lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    maxSize,
			MaxAge:     rotate.MaxAge,
			MaxBackups: rotate.MaxBackups,
			LocalTime:  rotate.LocalTime,
			Compress:   rotate.Compress,
		}
	}

	// Вернуть один из 12-ти вариантов реализации интерфейса Writer
	if pipe == nil && file == nil && logger == nil && writer == nil {
		return nullWriter{}
	} else if pipe != nil && file == nil && logger == nil && writer == nil {
		return pipeWriter{pipe}
	} else if pipe == nil && file != nil && writer == nil {
		return fileWriter{file}
	} else if pipe == nil && logger != nil && writer == nil {
		return rotatableWriter{logger}
	} else if pipe == nil && file == nil && logger == nil && writer != nil {
		return customWriter{logger}
	} else if pipe != nil && file != nil && writer == nil {
		return pipeAndFileWriter{Pipe: pipe, File: file}
	} else if pipe != nil && logger != nil && writer == nil {
		return pipeAndRotatableWriter{Pipe: pipe, Logger: logger}
	} else if pipe != nil && file == nil && logger == nil && writer != nil {
		return pipeAndCustomWriter{Pipe: pipe, Writer: writer}
	} else if pipe == nil && file != nil && writer != nil {
		return fileAndCustomWriter{File: file, Writer: writer}
	} else if pipe == nil && logger != nil && writer != nil {
		return rotatableAndCustomWriter{Logger: logger, Writer: writer}
	} else if pipe != nil && file != nil && writer != nil {
		return pipeAndFileAndCustomWriter{Pipe: pipe, File: file, Writer: writer}
	} else { // if pipe != nil && logger != nil && writer != nil {
		return pipeAndRotatableAndCustomWriter{Pipe: pipe, Logger: logger, Writer: writer}
	}
}

// Метод IsRotatable возвращает признак возможности ротации логов
func (_ nullWriter) IsRotatable() bool                      { return false }
func (_ pipeWriter) IsRotatable() bool                      { return false }
func (_ fileWriter) IsRotatable() bool                      { return false }
func (_ customWriter) IsRotatable() bool                    { return false }
func (_ pipeAndFileWriter) IsRotatable() bool               { return false }
func (_ rotatableWriter) IsRotatable() bool                 { return true }
func (_ pipeAndRotatableWriter) IsRotatable() bool          { return true }
func (_ pipeAndCustomWriter) IsRotatable() bool             { return false }
func (_ fileAndCustomWriter) IsRotatable() bool             { return false }
func (_ rotatableAndCustomWriter) IsRotatable() bool        { return true }
func (_ pipeAndFileAndCustomWriter) IsRotatable() bool      { return false }
func (_ pipeAndRotatableAndCustomWriter) IsRotatable() bool { return true }

// Метод Rotate производить ротацию логов, если она возможна
func (_ nullWriter) Rotate() error                 { return nil } // do nothing
func (_ pipeWriter) Rotate() error                 { return nil } // do nothing
func (_ fileWriter) Rotate() error                 { return nil } // do nothing
func (_ customWriter) Rotate() error               { return nil } // do nothing
func (_ pipeAndFileWriter) Rotate() error          { return nil } // do nothing
func (_ pipeAndCustomWriter) Rotate() error        { return nil } // do nothing
func (_ fileAndCustomWriter) Rotate() error        { return nil } // do nothing
func (_ pipeAndFileAndCustomWriter) Rotate() error { return nil } // do nothing

// Метод Close производит закрытие файла журнала, если он есть
func (_ nullWriter) Close() error                      { return nil } // do nothing
func (_ pipeWriter) Close() error                      { return nil } // do nothing
func (_ customWriter) Close() error                    { return nil } // do nothing
func (w pipeAndFileWriter) Close() error               { return w.File.Close() }
func (w pipeAndRotatableWriter) Close() error          { return w.Logger.Close() }
func (_ pipeAndCustomWriter) Close() error             { return nil } // do nothing
func (w pipeAndFileAndCustomWriter) Close() error      { return w.File.Close() }
func (w pipeAndRotatableAndCustomWriter) Close() error { return w.Logger.Close() }

// Метод Write для nullWriter не делает ничего
func (_ nullWriter) Write(_ []byte) (int, error) { return 0, nil }

// Метод Write для pipeAndFileWriter пишет в заданный канал и в заданный файл.
// Ошибка возвращается только по результату записи в файл (ошибка записи в pipe
// умалчивается).
func (w pipeAndFileWriter) Write(b []byte) (n int, err error) {
	w.Pipe.Write(b)
	return w.File.Write(b)
}

// Метод Write для pipeAndRotatableWriter пишел в заданный канал и
// в заданный файл с ротацией.
// Ошибка возвращается только по результату записи в файл (ошибка записи в pipe
// умалчивается).
func (w pipeAndRotatableWriter) Write(b []byte) (n int, err error) {
	w.Pipe.Write(b)
	return w.Logger.Write(b)
}

// Метод Write для pipeAndCustomWriter пишет в заданный канал и
// заданный io.Writer. Ошибка возвращается только по результату
// записи в io.Writer (ошибка записи в pipe умалчивается).
func (w pipeAndCustomWriter) Write(b []byte) (n int, err error) {
	w.Pipe.Write(b)
	return w.Writer.Write(b)
}

// Метод Write для fileAndCustomWriter пишет в заданный файл и
// заданный io.Writer. Ошибка возвращается только по результату
// записи в io.Writer (ошибка записи в файл умалчивается).
func (w fileAndCustomWriter) Write(b []byte) (n int, err error) {
	w.File.Write(b)
	return w.Writer.Write(b)
}

// Метод Write для rotatableAndCustomWriter пишет в заданный файл и
// заданный io.Writer. Ошибка возвращается только по результату
// записи в io.Writer (ошибка записи в файл умалчивается).
func (w rotatableAndCustomWriter) Write(b []byte) (n int, err error) {
	w.Logger.Write(b)
	return w.Writer.Write(b)
}

// Метод Write для pipeAndFileAndCustomWriter пишет в заданный канал,
// в заданный файл и в заданный io.Writer.
// Ошибка возвращается только по результату записи в io.Writer
// (ошибки записи в канал/файл умалчиваются).
func (w pipeAndFileAndCustomWriter) Write(b []byte) (n int, err error) {
	w.Pipe.Write(b)
	w.File.Write(b)
	return w.Writer.Write(b)
}

// Метод Write для pipeAndRotatableAndCustomWriter пишет в заданный канал,
// в заданный файл и в заданный io.Writer.
// Ошибка возвращается только по результату записи в io.Writer
// (ошибки записи в канал/файл умалчиваются).
func (w pipeAndRotatableAndCustomWriter) Write(b []byte) (n int, err error) {
	w.Pipe.Write(b)
	w.Logger.Write(b)
	return w.Writer.Write(b)
}

// EOF: "writer.go"
