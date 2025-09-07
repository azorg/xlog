// File: "xlog_test.go"

package xlog

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog" // go>=1.21
	"testing"
	"time"
	// FIXME: "golang.org/x/exp/slog" // экспериментальный пакет для go=1.20 только
)

// Все уровни логирования
var lvls = []string{
	"flood", "trace", "debug", "info", "notice", "warn", "error",
	"crit", "alert", "emerg", "fatal", "panic", "silent",
}

// Создать набор опций (*xlog.Opt)
var opt = NewOpt("log-")

// Тестирование рефлексии
func TestReflect(t *testing.T) {
	conf := Conf{
		Level: "info",
	}
	//Setup(conf)

	fmt.Println(Sprint(conf))
	fmt.Println(Sprint(&conf))

	sl := []string{"AAA", "BBB"}

	fmt.Println(Sprint(sl))
	fmt.Println(Sprint(&sl))

	Info("slice", "someSlice", sl)

	//Fatal("")
}

// Тестирование логгера по умолчанию
func TestDefault(t *testing.T) {
	SetLvl("trace")

	slog.Debug("slog.Debug() by default")
	slog.Info("slog.Info() by default")

	Debug("Debug() by default", "level", GetLvl())
	Info("Info() by default")

	Alert("Alert() by default")
	Emerg("Emerg() by default")

	Notice("some slice", "slice", []string{"a", "b", "c"})

	//Fatal("Fatal() by default")
}

// Самое простое применение
func TestSimple(t *testing.T) {
	// Заполнить структуру конфигурации
	conf := Conf{
		Level: "info",
	}

	// Обогатить структуру конфигурации переменными окружения
	Env(&conf)

	// Обогадить conf опциями командной строки
	opt.UpdateConf(&conf)

	// Настроить все глобальные логгеры однотипно
	Setup(conf)

	// Обновить уровень логирования глобального логгера
	SetLevel(slog.LevelDebug)

	// Использовать глобальные логггеры log/slog
	slog.Debug("debug message", "value", 42)
	slog.Info("simple slog", "logLevel", GetLvl())
	slog.Error("error", "err", errors.New("some error"))
	log.Print("legacy logger")
}

// Вывести в цикле все метки логирования
func TestLevels(t *testing.T) {
	fmt.Println("Color level labels:")
	for _, lvl := range lvls {
		level := LevelFromString(lvl)
		fmt.Printf(" %3d: %s:%s (slog=%s)\n",
			int(level),
			LevelToLabel(level),
			LevelToColorLabel(level),
			level.String())
	}
}

// Проверить вывод логгера по умолчанию,
// но с обработкой переменных окружения
func TestFullDefaultEnv(t *testing.T) {
	// Структура конфигурации по умолчанию
	conf := Conf{Level: "flood"}

	//conf.Prefix = "PREFIX"
	conf.AddKey = "app"
	conf.AddValue = struct {
		Name    string
		Version string
	}{
		Name:    "testApp",
		Version: "0.0.1",
	}

	// Обогатить структуру конфигурации переменными окружения
	Env(&conf)

	// Обогадить conf опциями командной строки
	opt.UpdateConf(&conf)

	// Настроить все глобальные логгеры однотипно
	Setup(conf)

	slog.Info("slog.Info() - default X-logger", "someInt", 42)
	level := LevelFromString(conf.Level)
	WithGroup("conf").Notice("log level from Conf structure", "level", LevelToString(level))
	slog.Info("full X-logger configuration", "conf", conf)
	slog.Info("error message", "err1", errors.New("some error"))

	for _, lvl := range lvls {
		level := LevelFromString(lvl)
		Log(context.Background(), level, "default X-logger", "level", LevelToString(level))
	}

	log.Print("log.Print()")
}

// Проверка еще...
func TestYetAnother(t *testing.T) {
	// Структура конфигурации по умолчанию
	conf := Conf{Level: "trace"}

	// Обогатить структуру конфигурации переменными окружения
	Env(&conf)

	// Обгадить conf опциями командной строки
	opt.UpdateConf(&conf)

	// Создать logger (*xlog.Logger)
	logger := New(conf)
	slogger := logger.Logger

	logger.Notice("Привет, X-Log",
		"version", "1.0.0", "logLevel", logger.GetLvl())
	mylog := slogger.With("app", "helloworld")
	mylog.Info("application started")

	stuff := Fields{
		"user_id":   123,
		"ip":        "192.168.0.1",
		"timestamp": time.Now(),
	}

	mylog.Debug("user login", "stuff", stuff)
}

// Проверка ротации
func TestRotate(t *testing.T) {
	fmt.Println(">>> Test Rotate")

	conf := Conf{} // create default config
	conf.Pipe = "stdout"
	conf.Level = "flood" // set logger level
	conf.Format = "tint" // select tinted logger
	conf.Src = true      // add source file:line to log
	conf.SrcPkg = false  // add package name
	conf.SrcFunc = true  // add function mame to log
	conf.SrcExt = false  // remove ".go" extension
	conf.ColorOff = true // color OFF
	//conf.Time = true
	conf.TimeFormat = "dateTimeMilli" // add custom timestamp

	conf.File = "logs/test.log" // log file
	//conf.FileMode = "0600"

	conf.Rotate.Enable = true
	conf.Rotate.MaxBackups = 3

	// Обогатить структуру конфигурации переменными окружения
	Env(&conf, "LOG_")

	// Обогатить conf опциями командной строки
	opt.UpdateConf(&conf)

	log := New(conf) // create xlog with TintHandler
	//x.SetDefault() // set default xlog

	log.Info("hello 1")
	err := log.Rotate()
	if err != nil {
		log.Error("can't rotate", "err", err)
	}
	log.Info("hello 2")
	time.Sleep(time.Second)
	log.Close()
}

// Общая структура для всех записей в журнале (абстрактно)
type LogSystem struct {
	Os        string `json:"os"`
	OsVersion string `json:"osVersion"`
	Hostname  string `json:"hostanme"`
	Arch      string `json:"arch"`
}

var logSystem = LogSystem{
	Os:        "windows",
	OsVersion: "3.11",
	Hostname:  "main.kremlin.ru",
	Arch:      "e2k",
}

// Некоторый JSON
type SomeJSON struct {
	Num      int     `json:"num"`
	NickName string  `json:"nickName"`
	Value    float64 `json:"value"`
}

// Событие журнала
type LogEvent struct {
	Type   string    `json:"type"`
	Action string    `json:"action"`
	Status string    `json:"status"`
	System LogSystem `json:"system"`
	Other  *SomeJSON `json:"other"`
}

// Контекст записи в журнале
type LogContext struct {
	TraceId       string `json:"traceId"`
	CorrelationId string `json:"correlationID"`
}

// Пример журнала по черновику А.Мате
func TestMate(t *testing.T) {
	// Структура конфигурации по умолчанию
	conf := Conf{Level: "trace"}

	// Обогатить структуру конфигурации переменными окружения
	Env(&conf)

	// Обогатить conf опциями командной строки
	opt.UpdateConf(&conf)

	// Создать logger (*xlog.Logger)
	log := New(conf)
	log = log.With("system", logSystem)

	someJSON := SomeJSON{
		Num:      2022,
		NickName: "V.V.Putin",
		Value:    3.1415926,
	}

	evt := LogEvent{
		Type:   "LOGIN_FAILED",
		Action: "user_login",
		Status: "failure",
		System: logSystem,
		Other:  &someJSON,
	}

	ctx := LogContext{
		TraceId:       "12345678",
		CorrelationId: "87654321",
	}

	log.Notice("login user",
		"event", evt,
		"context", ctx)

	log.Emerg("emergency", "level", slog.Level(100),
		"someJSON", &someJSON, "conf", &conf)

	m := map[string]int{
		"a": 1,
		"b": 2,
	}

	s := []string{"qqq", "www"}

	log.Alert("alert", "map", m, "slice", s)

	a := [...]int{1, 2, 3, 4, 5}

	log.Notice("notice", "array", a, "dt", time.Second)

	//q := int64(0x0102030405)
	//fmt.Printf("q=0x%02X\n", byte(q>>16))
}

// Проверка групп
func TestGroup(t *testing.T) {

	// Структура конфигурации по умолчанию
	conf := Conf{Level: "trace"}

	// Обогатить структуру конфигурации переменными окружения
	Env(&conf)

	// Обогатить conf опциями командной строки
	opt.UpdateConf(&conf)

	// Создать logger (*xlog.Logger)
	log := New(conf)

	log.Info("root", "pi", 3.14)

	log1 := log.WithGroup("group")

	log2 := log1.With("app", "superApp")

	log2.Info("sub info", slog.Group("request", "val", 333))
}

// Проверка middleware
func TestMiddleware(t *testing.T) {
	srcFields := Fields{
		"app":    "test",
		"ver":    "1.2.3",
		"testID": 12345,
	}

	// Создать логгер для ошибок с выводом "os.Stderr"
	conf := Conf{Pipe: "stderr", Level: "debug", Format: "tint", ColorOff: false,
		Src: false, SrcPkg: false, SrcFunc: true, SrcFields: &srcFields,
		GoId: true, IdOn: true, SumOn: true, SumAlone: true,
		AddKey: "errorLog", AddValue: true}
	//Env(&conf)            // обогатить структуру конфигурации переменными окружения
	//opt.UpdateConf(&conf) // обогатить conf опциями командной строки
	logErr := New(conf).With("logErr", true)
	logErr.Info("logErr.Info", "anyValue", 3.1415926)

	// Создать "основной" логгер с middleware
	conf = Conf{Pipe: "", Level: "trace", Format: "slog",
		Src: false, SrcPkg: true, SrcFunc: true, SrcFields: &srcFields,
		GoId: false, IdOn: false, SumOn: true, SumAlone: true,
		AddKey: "addKey", AddValue: "addValue"}
	Env(&conf)            // обогатить структуру конфигурации переменными окружения
	opt.UpdateConf(&conf) // обогатить conf опциями командной строки
	log := New(conf, NewMiddlewareForError(logErr))
	log = log.WithMiddleware(NewMiddlewareNoPasswd())
	log = log.With("pin", 1234)

	srcFields["srcField"] = "extra source"

	log.Info("hello middleware", "passwd", "superSecret")
	log.Error("hello middleware with error", "err", errors.New("fake error"), "passwd", "SECRET-1234")
}

// Проверка Fields
func TestFields(t *testing.T) {
	// Настролить логгер(ы)
	conf := Conf{Level: "debug"}
	Env(&conf)            // обогатить структуру конфигурации переменными окружения
	opt.UpdateConf(&conf) // обогатить conf опциями командной строки
	Setup(conf)

	Slog().LogAttrs(context.Background(),
		slog.LevelInfo, "hello fields (Attrs)",
		Fields{
			"app": "test",
			"now": time.Now(),
			"num": 65536,
		}.Attrs()...)

	Error(
		"error with fields (Args)",
		"err", errors.New("error"),
		slog.Group(
			"header",
			Fields{
				"type": "magic",
				"size": 16384,
			}.Args()...))

	log := WithAttrs(Fields{
		"application": "testApplication",
		"version":     "1.2.3",
	}.Attrs())

	log.Notice("notice with fields (Args)", "someFlag", true)

	log.Debug(
		"debug with fields",
		Fields{
			"now":        time.Now(),
			"superValue": 123,
			"str":        "some string",
			"pi":         3.1415926,
		}.Args()...)

	logXY := WithFields(Fields{
		"x": 1,
		"y": 2,
	})

	logXY.Debug("vector", "z", 3)
}

// Проверка WithFields
func TestWithFields(t *testing.T) {
	// Настролить логгер(ы)
	conf := Conf{
		Level: "debug",
		Src:   true, SrcFields: &Fields{"src": true},
	}
	Env(&conf)            // обогатить структуру конфигурации переменными окружения
	opt.UpdateConf(&conf) // обогатить conf опциями командной строки
	Setup(conf)

	fields := Fields{"eee": 321}
	fields["app"] = "app"
	fields["cnt"] = 123.456
	log := With("fff", fields)
	fields["eee"] = "-+-"
	log.Info("with fields")

	//if true {
	//	return
	//}

	log = log.With("withSingle", "one")
	log = log.WithGroup("grp")
	log = log.With("withGroup", "grp")
	log.Info("text WithFields()", "pi-pi", 3.1415926)
}

// Проверка Fields #2
func TestFields2(t *testing.T) {
	// Настролить логгер
	srcFields := Fields{"src": false}

	conf := Conf{
		Level: "debug",
		Src:   true, SrcFields: &srcFields,
	}

	Env(&conf)            // обогатить структуру конфигурации переменными окружения
	opt.UpdateConf(&conf) // обогатить conf опциями командной строки

	log := New(conf)

	log.Info("hello #1 (src must false)")

	srcFields["src"] = true

	log.Info("hello #2 (src must true)")

	v := Fields{"x": 1, "y": 2}

	log = log.WithFields(Fields{"a": 1980})
	log = log.WithGroup("v")
	log = log.WithFields(v)

	//log.Info("grp value", "grp", v)
	log.Info("hello #3", "z", 3)

}

// Проверка "Multi Handler"
func TestMultiHandler(t *testing.T) {

	// Настройки журнала для вывода в JSON файл с ротацией
	conf := Conf{
		Level:     "debug",
		File:      "logs/file.log",
		FileMode:  "0600",
		Format:    "json", // JSON handler
		GoId:      true,
		IdOn:      true,
		SumOn:     true,
		SumFull:   true,
		SumAlone:  true,
		TimeLocal: false, // UTC
		Src:       true,
		SrcPkg:    true,
		SrcFunc:   true,
		SrcFields: &Fields{
			"id":   "itcagent", // идентификатор сервиса
			"host": "arm1.dev.corp.com",
		},
		Rotate: RotateConf{
			Enable:     true,
			MaxSize:    5,     // MB
			MaxAge:     7,     // days
			MaxBackups: 100,   // number
			LocalTime:  false, // UTC
			Compress:   true,
		},
	}

	// Создать логгер (хендлер) для вывода в файл
	log := New(conf)

	// Настройки "человеческого" журнала для вывода на stdout
	conf = Conf{
		Level:      "trace",
		Pipe:       "stdout",
		Format:     "tint", // tinted handler
		GoId:       false,
		IdOn:       false,
		SumOn:      false,
		TimeLocal:  true,    // Local time
		TimeFormat: "space", // 2006-01-02_15.04.05.999999999
		Src:        true,
		SrcPkg:     false,
		SrcFunc:    false,
		ColorOff:   false,
	}

	// Создать "Multi Handler" логгер (stdout/text + file/JSON)
	log = New(conf, NewMiddlewareMulti(log))

	log.Debug("Hello, Multi Handler!", "cnt", 1) // попадет в JSON файл и на stdout
	log.Trace("Hello, Multi Handler!", "cnt", 2) // попадет только на stdout
	log.Flood("Hello, Multi Handler!", "cnt", 3) // будет пропущено
}

// EOF: "xlog_test.go"
