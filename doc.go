/*
Пакет xlog реализует простую надстройку на стандартными логерами log и slog.

Логер slog включен в стандартную поставку Go начиная с версии 1.21 ("log/slog").
До этого логер представлен экспериментальным пакетом "golang.org/x/exp/slog".
Экспериментальный пакет содержит ошибки (ждем когда исправят, но пока не спешат).

Структуры данных:

	Conf - Обобщенная структура конфигурации логера, имеет JSON тэги

	Xlog - Структура/обёртка над slog для добавления методов типа Debugf/Noticef/Errorf/Trace

Функции настройки конфигурации:

	NewConf() - заполнить обобщенную структуру конфигурации логера значениями по умолчанию

	SetupLog() - настроить стандартный логер в соответствии с заданной структурой конфигурации

	NewLog() - создать стандартный логер log.Logger в соответствии со структурой конфигурации

	NewSlog() - создать структурированный логер slog.Logger в соответствии со структурой конфигурации

	Setup() - настроить стандартный и структурированный логеры по умолчанию в соответствии с структурой конфигурации

	GetLevel() - вернуть текущий уровень журналирования

	GetLvl() - вернуть текущий уровень журналирования в виде строки

Функции для работы с надстройкой Xlog:

	Default() - Создать логер Xlog на основе исходного slog.Deafult()

	Current() - Вернуть текущий глобальный логер Xlog

	Slog() - Вернуть текущий глобальный логер slog.Logger

	X() - Создать логер Xlog на основе логера slog

	New() - Cоздать новый логер Xlog с заданными параметрами конфигурации

Методы для работы с Xlog:

	Slog() - Обратное преобразование Xlog -> *slog.Logger

	SetDefault() - Установить логер как xlog по умолчанию

	SetDefaultLogs() - Установить логер как log/slog/xlog по умолчанию

Методы для использования Xlog с дополнительными уровнями:

	Log(level Level, msg string, args ...any)
	Trace(msg string, args ...any)
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Notice(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
	Panic(msg string)

	Logf(evel Level, format string, args ...any)
	Tracef(format string, args ...any)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Noticef(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)

Примечание: имеются аналогичные глобальные функции в пакете
для использования глобального логера.
*/
package xlog

// EOF: "doc.go"
