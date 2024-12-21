xlog - надстройка на стандартными логерами log и slog
=====================================================

## Краткий экскурс

slog - это "конструктор", "концепция", которая позволяет построить нужный именно
Вам логгер в Go приложениях начиная с версии go-1.21. Для более ранних версий
имеется одноименный пакет "golang.org/x/exp/slog".

Пакет `xlog` - это пример использования slog с возможными его расширениями
(имеется дополнительный tinted handler, для "пестрой" раскраски с помощью
ANSI/Escape последовательностей, добавлены дополнительные уровни журналирования
типа Trace, Notice, Critical, Fatal и т.п.).

Заложены "мостики" единообразного поведения стандартного логгра "log" при работе
через настроенный логгер slog/xlog.

Для подробностей см. "doc.go".

В ряде случаев из всего пакета может быть полезна только одна функция Setup(),
которая на основе структуры конфигурации (которая содержит JSON-теги для
непосредственной сериализации) может настроить нужное поведение глобального
стандартного Go логера "log" и структурного логера "slog".

Можно предположить вариант создания структуры конфигурации "по умолчанию"
с помощью функции NewConf(), внесения необходимых изменений (управление
уровнем журналирования например) и последующим вызовом Setup().
Например так:
```
  conf = xlog.NewConf() // создать структура конфигурации по умолчанию
  conf.Level = "debug" // задать уровень отладки
  conf.Tint = true // включить структурный цветной логер
  conf.Src = true // включить в журнал имена файлов и номера строк исходных текстов
  xlog.Setup(conf) // настроить глобальный логер (xlog + slog + log)
```

### Возможные конфигурации логера
| Slog  | JSON  | Tint  | Описание                        |
|:-----:|:-----:|:-----:|:--------------------------------|
| false | false | false | логер по умолчанию              |
| true  | false | false | структурный TextHandler         |
| x     | true  | false | структурный JSONHandler         |
| x     | x     | true  | структурный цветной TintHandler |

Если планируется конфигурировать только стандартный логер "log", то достаточно
применения функций NewLog()/SetupLog() заместо Setup().

Для конфигурирования только структурного логера slog стоит использовать
функцию NewSlog(). Функция получает структуру конфигурации и возвращает
стандатный `*slog.Logger`. Расширенная функция NewSlogEx() возвращает кроме
slog.Logger ещё и xlog.Leveler - интерфейс, который может использоваться
для изменения уровня журналирования после исходной инициализации.

Для дополнительных "сахарных" функций реализована обёртка xlog.Logger,
которая хранит указатель на структурный логер `*slog.Logger` и xlog.Leveler,
и предоставляет дополнительные методы типа `Infof()`, `Trace()`, `Fatalf()` и др.
Имеются функции - аналоги соответствующих методом для использования
глобального логера. Через функции и методы SetLevel()/SetLvl() можно
изменять уровень логирования "на лету".

Вот перечень "сахарных" методов xlog.Logger и соответствующих функций пакета:
```
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

  GetLevel() slog.Level
  SetLevel(level slog.Level)
  GetLvl() string
  SetLvl(level string)
```

## Ротация логов
Используется популярный пакет https://github.com/natefinch/lumberjack

## Полезные ссылки на документацию

* https://go.dev/blog/slog - Structured Logging with slog (August 2023)

* https://pkg.go.dev/log - стандартный логер в Go

* https://pkg.go.dev/log/slog - `slog` из стандартной поставки начиная с go-1.21.0

* https://pkg.go.dev/golang.org/x/exp/slog - `slog` для go-1.20 и младше

* https://github.com/golang/example/blob/master/slog-handler-guide/README.md

## Как перейти с go1.20.x на go1.21.x и старше?

1. Убрать из секцию import "golang.org/x/exp/slog"

2. Добавить в секцию import "log/slog"

3. В файле "const.go" задать `OLD_SLOG_FIX = false`

Переход на "golang.org/x/exp/slog" (для go-1.20) делается в обратной
последовательности.

