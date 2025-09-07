// File: "setup.go"

package xlog

import (
	"io"
	"log"
	"log/slog" // go>=1.21
	// FIXME: "golang.org/x/exp/slog" // экспериментальный пакет для go=1.20 только
)

// SetupLogWithWriter производит настройку стандартного *log.Logger на основе
// унифицированной структуры конфигурации Conf для "Client Logger" с
// направлением вывода в заданный io.Writer вместо заданного файла
func SetupLogWithWriter(logger *log.Logger, conf Conf, writer io.Writer) {
	flag := 0
	if !conf.TimeOff {
		flag |= log.LstdFlags
		if conf.TimeMicro {
			flag |= log.Lmicroseconds
		}
	}
	if conf.Src {
		if conf.SrcPkg {
			flag |= log.Llongfile
		} else {
			flag |= log.Lshortfile
		}
	}
	if conf.Prefix != "" {
		flag |= log.Lmsgprefix
		logger.SetPrefix(conf.Prefix + " ")
	}

	logger.SetOutput(writer)
	logger.SetFlags(flag)
}

// SetupLog производит настройку стандартного *log.Logger на основе
// унифицированной структуры конфигурации Conf для "Client Logger" с
// направлением вывода в файл с ротацией (если предусмотрено конфигурацией).
// Функция вызывает последовательно NewWriter() и SetupLogWithWriter().
func SetupLog(logger *log.Logger, conf Conf) {
	writer := NewWriter(conf.Pipe, conf.File, conf.FileMode, &conf.Rotate, nil)
	SetupLogWithWriter(logger, conf, writer)
}

// Setup - мега функция, которая настраивает все глобальные логгеры
// в соответствии с заданной структурой конфигурации Conf.
// Функция потоко не безопасная!
func Setup(conf Conf) {
	// Настроить стандартный (legacy) логгер
	SetupLog(defaultLog, conf)

	// Настроить структурированный slog логгер
	logger := New(conf)

	// Сохранить структурированный логгер как текущий (глобальный)
	currentClog = logger

	// Настроить глобальный логгер slog.
	// При этом legacy лог автоматически перенаправляется в данный slog.Logger.
	slog.SetDefault(logger.Logger)

	if logFmt := logFormat(conf.Format); logFmt == logFmtDefault {
		// Грязный хук для исключения "зацикливания" логгеров
		// в связи с особенностью реализации slog.SetDefault()
		SetupLog(defaultLog, conf)
	}

	// FIXME: TODO
	// В экспериментальном slog есть ошибка:
	// При добавлении хендлера для управления уровнем
	// логирования некорректно выводятся имена файлов (и строк).
	// Что интересно, в Go v1.21 в log/slog всё исправлено.
	// Используя Go до версии 1.21 (например 1.20) при включении
	// управления уровнем логирования при работе slog через slog.defaultHandler
	// в угоду возможности управления уровнями отключаем вывод файлов и строк.
	// Если "golang.org/x/exp/slog" доработают, то этот FIX можно будет убрать.
	if oldSlogFix { // "runtime.Version() < go1.21.0"
		if currentClog.GetLevel() != DefaultLevel {
			flag := defaultLog.Flags()
			flag = flag &^ (log.Lshortfile | log.Llongfile) // sorry...
			defaultLog.SetFlags(flag)
		}
	}
}

// EOF: "setup.go"
