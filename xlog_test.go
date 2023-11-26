// File: "xlog_test.go"

package xlog

import (
	"errors"
	"fmt"
	"log"
	"log/slog" // go>=1.21
	"testing"
	"time"
	//"golang.org/x/exp/slog" // depricated for go>=1.21
)

func TestUsage(t *testing.T) {
	fmt.Println(">>> Test Usage")

	conf := NewConf()          // create default config (look xlog.Conf for details)
	conf.Level = "trace"       // set logger level
	conf.Tint = true           // select tinted logger
	conf.Src = true            // add source file:line to log
	conf.TimeTint = "15:04:05" // add custom timestamp
	x := New(conf)             // create xlog with TintHandler
	x.SetDefault()             // set default xlog

	err := errors.New("some error")
	count := 12345

	Trace("Tinted logger xlog.Trace()", "level", conf.Level)
	Debug("Tinted logger xlog.Debug()")
	Info("Tinted logger xlog.Info()", "count", count)

	x.Notice("Tinted logger x.Notice()")
	x.Warn("Tinted logger x.Warn()")
	x.Error("Tinted logger x.Error()", Err(err))

	sl := x.Slog() // *slog.Logger may used too
	sl.Info("Tinted logger is *slog.Logger sl.Info()", "str", "some string")

	fmt.Println()
}

func TestClassic(t *testing.T) {
	fmt.Println(">>> Test Classic log")

	sign := 0x55AA
	err := errors.New("simple error")
	log.Printf("[classic log from box: log.Printf(): sign=0x%04X, err='%v']", sign, err)

	n := NewLog(Conf{})
	n.Print("[n := xlog.NewLog(xlog.Conf{}); n.Print(...)]")

	n = NewLog(Conf{Src: true})
	n.Print("[n := xlog.NewLog(xlog.Conf{Src: true}); n.Print(...)]")

	n = NewLog(Conf{Src: true, Time: true})
	n.Print("[n := xlog.NewLog(xlog.Conf{Src: true, Time: true}); n.Print(...)]")

	n = NewLog(Conf{Time: true, TimeUS: true})
	n.Print("[n := xlog.NewLog(xlog.Conf{Time: true, TimeUS: true}); n.Print(...)]")

	fmt.Println()
}

func TestSlogDefault(t *testing.T) {
	fmt.Println(">>> Test slog default handler")

	count := 12345
	err := errors.New("simple error")

	slog.Info("Default slog.Infof() from box", "count", count, "err", err)

	s := NewSlog(Conf{Level: "debug"})
	s.Debug(`[s := xlog.NewSlog(xlog.Conf{Level: "debug"}); s.Debug(...)]`, "count", count, "err", err)
	s.Info(`[s := xlog.NewSlog(xlog.Conf{Level: "debug"}); s.Info(...)]`, "count", count, "err", err)
	s.Warn(`[s := xlog.NewSlog(xlog.Conf{Level: "debug"}); s.Warn(...)]`, "count", count, "err", err)
	s.Error(`[s := xlog.NewSlog(xlog.Conf{Level: "debug"}); s.Error(...)]`, "count", count, "err", err)

	s = NewSlog(Conf{Src: true})
	s.Info(`[s := xlog.NewSlog(xlog.Conf{Src: true}); s.Info(...)]`, "count", count, "err", err)

	s = NewSlog(Conf{Src: true, SrcLong: true})
	s.Error(`[s := xlog.NewSlog(xlog.Conf{Src: true, SrcLong: true}); s.Error(...)]`, "count", count, "err", err)

	s = NewSlog(Conf{Time: true})
	s.Info(`[s := xlog.NewSlog(xlog.Conf{Time: true}); s.Info(...)]`, "count", count, "err", err)

	s = NewSlog(Conf{Time: true, TimeUS: true})
	s.Warn(`[s := xlog.NewSlog(xlog.Conf{Time: true, TimeUS: true}); s.Warn(...)]`, "count", count, "err", err)

	fmt.Println()
}

func TestSlogText(t *testing.T) {
	fmt.Println(">>> Test slog text handler")

	count := 12345
	err := errors.New("simple error")

	s := NewSlog(Conf{Slog: true})
	s.Info(`[s := xlog.NewSlog(xlog.Conf{Slog: true}); s.Info(...)]`, "count", count, "err", err)

	s = NewSlog(Conf{Slog: true, Src: true})
	s.Info(`[s := xlog.NewSlog(xlog.Conf{Slog: true, Src: true}); s.Info(...)]`, "count", count, "err", err)

	s = NewSlog(Conf{Slog: true, Src: true, SrcLong: true})
	s.Info(`[s := xlog.NewSlog(xlog.Conf{Slog: true, Src: true, SrcLong: true}); s.Info(...)]`,
		"count", count, "err", err)

	s = NewSlog(Conf{Slog: true, Src: true, Time: true})
	s.Info(`[s := xlog.NewSlog(xlog.Conf{Slog: true, Src: true, Time: true}); s.Info(...)]`,
		"count", count, "err", err)

	s = NewSlog(Conf{Slog: true, Time: true, TimeUS: true})
	s.Info(`[s := xlog.NewSlog(xlog.Conf{Slog: true, Time: true, TimeUS: true}); s.Info(...)]`,
		"count", count, "err", err)

	fmt.Println()
}

func TestSlogJSON(t *testing.T) {
	fmt.Println(">>> Test slog JSON handler")

	count := 12345
	err := errors.New("simple error")

	s := NewSlog(Conf{JSON: true})
	s.Info(`[s := xlog.NewSlog(xlog.Conf{JSON: true}); s.Info(...)]`, "count", count, "err", err)

	s = NewSlog(Conf{JSON: true, Src: true})
	s.Info(`[s := xlog.NewSlog(xlog.Conf{JSON: true, Src: true}); s.Info(...)]`, "count", count, "err", err)

	s = NewSlog(Conf{JSON: true, Src: true, SrcLong: true})
	s.Info(`[s := xlog.NewSlog(xlog.Conf{JSON: true, Src: true, SrcLong: true}); s.Info(...)]`,
		"count", count, "err", err)

	s = NewSlog(Conf{JSON: true, Src: true, Time: true})
	s.Info(`[s := xlog.NewSlog(xlog.Conf{JSON: true, Src: true, Time: true}); s.Info(...)]`,
		"count", count, "err", err)

	s = NewSlog(Conf{JSON: true, Time: true, TimeUS: true})
	s.Info(`[s := xlog.NewSlog(xlog.Conf{JSON: true, Time: true, TimeUS: true}); s.Info(...)]`,
		"count", count, "err", err)

	fmt.Println()
}

func TestTintHandler(t *testing.T) {
	fmt.Println(">>> Test tinted handler")

	err := errors.New("some error")
	str := "some-string-value"
	cnt := 123

	x := New(Conf{Tint: true})
	x.Info(`[x := xlog.New(xlog.Conf{Tint: true}); x.Info(...)]`,
		"cnt", cnt, "str", str, "err", err)

	x = New(Conf{Tint: true, Src: true})
	x.Warn(`[x := xlog.New(xlog.Conf{Tint: true, Src: true}); x.Warn(...)]`)

	x = New(Conf{Tint: true, Src: true, SrcLong: true})
	x.Warn(`[x := xlog.New(xlog.Conf{Tint: true, Src: true, SrcLong: true}); x.Warn(...)]`)

	x = New(Conf{Tint: true, Time: true})
	x.Error(`[x := xlog.New(xlog.Conf{Tint: true, Time: true}); x.Error(...)]`)

	x = New(Conf{Tint: true, Time: true, TimeUS: true})
	x.Error(`[x := xlog.New(xlog.Conf{Tint: true, Time: true, TimeUS: true}); x.Error(...)]`)

	x = New(Conf{Tint: true, TimeTint: time.Kitchen})
	x.Notice(`[x := xlog.New(xlog.Conf{Tint: true, TimeTint: time.Kitchen}); x.Notice(...)]`)

	x = New(Conf{Tint: true, TimeTint: time.TimeOnly})
	x.Notice(`[x := xlog.New(xlog.Conf{Tint: true, TimeTint: time.TimeOnly}); x.Notice(...)]`)

	x = New(Conf{Tint: true, TimeTint: TimeOnlyMicro})
	x.Notice(`[x := xlog.New(xlog.Conf{Tint: true, TimeTint: xlog.TimeOnlyMicro}); x.Notice(...)]`)

	x = New(Conf{Tint: true, TimeTint: "15:04:05.999999999"})
	x.Notice(`[x := xlog.New(xlog.Conf{Tint: true, TimeTint: "15:04:05.999999999"}); x.Notice(...)]`)

	x = New(Conf{Tint: true, Level: "trace"})
	x.Trace(`[x := xlog.New(xlog.Conf{Tint: true, Level: "trace"}); x.Trace(...)]`, Err(err))

	fmt.Println()
}

func TestSlogToXlog(t *testing.T) {
	fmt.Println(">>> Test Xlog <-> *slog.Logger")

	s := NewSlog(Conf{Tint: true, Level: "trace"}) // *slog.Logger
	x := X(s)                                      // *slog.Logger -> xlog.Xlog
	l := x.Slog()                                  // xlog.Xlog -> *slog.Logger

	x.Slog().Info("x.Slog().Info()")
	X(x.Slog()).Trace("xlog.X(x.Slog()).Trace")
	l.Warn("l.Warn()")

	fmt.Println()
}

func TestSetDefault(t *testing.T) {
	fmt.Println(">>> Test xlog.SetDefault() and xlog.SetDefaultLogs()")

	x := New(Conf{Tint: true, Level: "trace", Time: true, TimeUS: true})
	x.SetDefault()
	Notice("xlog.Notice() after x.SetDefault()")

	slog.Info("slog.Info() by default")
	log.Print("log.Print() by default")

	x.SetDefaultLogs()
	slog.Info("slog.Info() after x.SetDefaultLogs()")
	log.Print("log.Print() after x.SetDefaultLogs()")

	fmt.Println()
}

// EOF: "xlog_test.go"
