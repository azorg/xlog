/*
Пакет xlog реализует простую надстройку на стандартными логерами log и slog.

Логер slog включен в стандартную поставку Go начиная с версии 1.21 ("log/slog").
До этого логер представлен экспериментальным пакетом "golang.org/x/exp/slog".

Структуры данных:

	Conf - Обобщенная структура конфигурации логера, имеет JSON тэги

	Logger - Структура/обёртка над slog для добавления методов типа Debugf/Noticef/Errorf/Trace

Интерфейсы:

	Xlogger - интерфейс к структуре Logger (приведен для наглядности API)

	Leveler - интерфейс управления уровнем журналирования

Функции настройки конфигурации:

	NewConf() - заполнить обобщенную структуру конфигурации логгера значениями по умолчанию

	SetupLog() - настроить стандартный логгер в соответствии с заданной структурой конфигурации

	NewLog() - создать стандартный логгер log.Logger в соответствии со структурой конфигурации

	NewSlog() - создать структурированный логгер slog.Logger в соответствии со структурой конфигурации

	NewSlogEx() - создать структурированный логгер slog.Logger и вернуть интерфейс Leveler

	Setup() - настроить стандартный и структурированный логгеры по умолчанию в соответствии с структурой конфигурации

	GetLevel() - вернуть текущий уровень журналирования

	GetLvl() - вернуть текущий уровень журналирования в виде строки

Функции для работы с надстройкой Logger:

	Default() - Создать логер на основе исходного slog.Deafult()

	Current() - Вернуть текущий глобальный логер

	Slog() - Вернуть текущий глобальный логер slog.Logger

	X() - Создать логер на основе логера slog (для доступа к "сахарным" методам)

	New() - Cоздать новый логер с заданными параметрами конфигурации

Методы для работы с Logger (методы интерфейса Xlogger):

	With() - Создать дочерний логгер с дополнительными атрибутами

	WithAttrs() - Создать дочерний логгер с дополнительными атрибутами

	WithGroup() - Создать дочерний логгер с группировкой ключей

	Slog() - Обратное преобразование *xlog.Logger -> *slog.Logger

	SetDefault() - Установить логер как xlog по умолчанию

	SetDefaultLogs() - Установить логер как log/slog/xlog по умолчанию

	GetLevel() - получить текуший уровень журналирования (slog.Level)

	SetLevel(l) - обновить текущий уровень журналирования (slog.Level)

	GetLvl() - получить текущий уровень журналирования в виде строки

	SetLvl() - обновить текущий уровень журналирования в виде строки

	Write(p []byte) (n int, err error) - метод для соответствия io.Writer

	NewLog(prefix string) *log.Logger - вернуть стандартный логгер с префиксом

Методы для использования xlog.Logger с дополнительными уровнями:

	Log(level slog.Level, msg string, args ...any)
	Flood(msg string, args ...any)
	Trace(msg string, args ...any)
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Notice(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Crit(msg string, args ...any)
	Fatal(msg string, args ...any)
	Panic(msg string)

	Logf(level slog.Level, format string, args ...any)
	Floodf(format string, args ...any)
	Tracef(format string, args ...any)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Noticef(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Critf(format string, args ...any)
	Fatalf(format string, args ...any)

Примечание: имеются аналогичные глобальные функции в пакете
для использования глобального логера.

Вспомогательные функции работы с уровнями журналирования:

	ParseLvl(lvl string) slog.Level - получить уровень из строки типа "debug"

	ParseLevel(level slog.Level) string - преобразовать уровень к строке

Методы интерфейса Leveler:

	Level() slog.Level - получить уровень журналирования (реализация интерфейса slog.Leveler)

	Update(slog.Level) - обновить уровень журналирования

	String() string - сформировать метку для журнала

	ColorString() string - сформировать метку для журнала с ANSI/Escape подкраской
*/
package xlog

// EOF: "doc.go"
