// File: "test.go"

package main

import (
  "time"
  "errors"
  "math"
  "log/slog"
	"context"
	
  //"github.com/gofrs/uuid"
	"github.com/azorg/xlog"
)

type valuer struct {}

func (_ valuer) LogValue() slog.Value {
  //return slog.StringValue(time.Now().Format(time.RFC3339))
  return slog.StringValue(time.Now().Format(time.RFC3339Nano))
}

// Тестирование групп
func someGroup(log *xlog.Logger) {
  log.Info("исходный журнал", slog.Float64("pi", 3.14), "xyz", 12345678.)
  log1 := log.WithGroup("группа1")
  log1.Info("журнал с группой 1", "e", 2.7)
  log2 := log1.With("app", "grptest")
  log2.Infof("журнал с группой и With(app, ...) x=%d", 123)
  log2.Info("журнал с группой и With(app, ...)", "y", 321)
  log3 := log2.WithGroup("группа2")
  log3.Info("журнал с группой 2", "name", "vova")
  log1.Info("журнал с группой 1 опять", "eee", "EEE")
  log.Info("исходный журнал опять", slog.Float64("pi", 3.14))
  log2.Info("журнал с группой и With(app, ...)", "xyz", 12345678)
	log2.Log(context.Background(), -1, "странный уровень")
}

// Генерация тестового JSON журнала
func test(logConf xlog.Conf) {
	// Создать логгер для ошибок с выводом в "os.Stderr"
	conf := xlog.Conf{
		Pipe: "stderr", Level: "flood",
		Format: "tint", ColorOff: false,
		//Src: true, SrcPkg: true,
		Src: false, SrcPkg: false,
		GoId: false, IdOn: false, SumOn: false, SumAlone: false,
	}
	logErr := xlog.New(conf)
	logErr.Info("logErr.Info")

	// Создать "основной" логгер с middleware
	// Выводить только в формате JSON
  logConf.Format = "JSON"
	logConf.SrcFields = &xlog.Fields{
		"app": "test",
		"ver": "0.0.3",
	}
	log := xlog.New(logConf, xlog.NewMiddlewareForError(logErr))
  //log := xlog.New(logConf)
  
  someGroup(log)

  log.Info("start test #0")
  log.Info("start test #1")

  //log = log.With("application", APP_NAME, "version", Version)
  log = log.With("application", APP_NAME)
  
	//log.Notice("start test #2", "application", APP_NAME)
	log.Notice("start test #2")
	
	//if true {
	//	return //!!!!
	//}
	
	log2 := log.WithGroup("testGroup")

  var v slog.LogValuer = valuer{}
  log2 = log2.With("someTime", v.LogValue())

  s := struct{
    A string
    X string
  }{"abc", "eee"}
  _ = s

  m := map[int]float64{
    0: 0.,
    1: 1.,
  }
  _ = m
	
  err := errors.New("some error")

	//if true {
	//	return //!!!!
	//}

	log.Info("start test #2",
    "s", s, "logConf", logConf, "now", time.Now(), xlog.Err(err),
    "pi", math.Pi, "complex", complex(math.Pi/2, math.E))
	
	//if true {
	//	return //!!!!
	//}
  
  log.With(slog.Group("vector", slog.Int("x", 1), slog.Float64("y", 2.33))).Flood("group", "value", 123)
  log.With(slog.Group("vector", "x", 1.2, "y", 2, slog.Int("z", 0))).Info("with 1", "q", 123)
  log.With(slog.Group("vector", "x", 1.2, "y", 2, "z", 0)).Info("with 2", "q", 123)
  log.With(slog.Group("vector", "x", 1.2, "y", 2), slog.Int("z", 0)).Info("with 3", "q", 123)
  log.With("x", 1, "y", 2).Trace("with vector")
  log.Alert("slog group 1", slog.Group("vec", "x", 1, "y", 2))
  
  log.Crit("slog group 2", slog.Group("vector", "x", 1, "y", 2), "val", "43")
  log.Emerg("val=42", "val42", 42)
  
  log2.Debug("stop test",
    "nowUTC", time.Now().UTC(),
    "e", math.E, "map", m)

	//if true {
	//	return //!!!!
	//}

	fields := xlog.Fields{
		"inetger": 1,
		"str":     "hello_fields",
	}
	log3 := log2.WithFields(fields)
	fields["axtra"] = 1980
	log3.Info("record with fields", "endValue", 9)
}

// EOF: "test.go"
