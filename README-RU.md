xlog - надстройка на стандартными логерами log и slog
=====================================================

## Краткий экскурс

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
функцию NewSlog().

Для дополнительных "сахарных" функций реализована обёртка Xlog, которая хранит
указатель на структурный логер `*slog.Logger`, но предоставляет дополнительные
методы типа `Infof()`, `Trace()`, `Fatalf()` и др. Имеются функции - аналоги
соответствующих методом для использования глобального логера.

Вот перечень "сахарных" методов Xlog и соответствующих функций пакета:
```
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
```

## Полезные ссылки на документацию

* https://go.dev/blog/slog - Structured Logging with slog (August 2023)

* https://pkg.go.dev/log - стандартный логер в Go

* https://pkg.go.dev/log/slog - `slog` из стандартной поставки начиная с go-1.21.0

* https://pkg.go.dev/golang.org/x/exp/slog - экспериментальный `slog`

## Как перейти с go1.20.x на go1.21.x и старше?

1. Убрать из секцию import "golang.org/x/exp/slog"

2. Добавить в секцию import "log/slog"

3. В файле "const.go" задать `OLD_SLOG_FIX = false`

Переход на "golang.org/x/exp/slog" (например для go-1.19) делается в обратной
последовательности.
