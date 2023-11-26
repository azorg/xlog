`xlog` - yet another "log/slog"  backend/frontend wrappers and tinted 🌈 `slog.Handler`
=======================================================================================

Package `xlog` implements some wrappers to work with structured logger
[`slog`](https://pkg.go.dev/log/slog) and classic simple logger
[`log`](https://pkg.go.dev/log) too.

Code of xlog.TintHandler based on [`tint`](https://github.com/lmittmann/tint).

![Tinted xlog](https://github.com/azorg/xlog/blob/main/img/xlog-tinted.png "xlog-tinded.png")

```
go get github.com/azorg/xlog
```

## Usage

```go
  conf := NewConf()          // create default config (look xlog.Conf for details)
  conf.Level = "trace"       // set logger level
  conf.Tint = true           // select tinted logger
  conf.Src = true            // add source file:line to log
  conf.TimeTint = "15:04:05" // add custom timestamp
  x := New(conf)             // create xlog with TintHandler
  x.SetDefault()             // set default xlog
	
  err := errors.New("some error")
  count := 12345

  xlog.Trace("Tinted logger xlog.Trace()", "level", conf.Level)
  xlog.Debug("Tinted logger xlog.Debug()")
  xlog.Info("Tinted logger xlog.Info()", "count", count)

  x.Notice("Tinted logger x.Notice()")
  x.Warn("Tinted logger x.Warn()")
  x.Error("Tinted logger x.Error()", Err(err))
	
  sl := x.Slog() // *slog.Logger may used too
  sl.Info("Tinted logger is *slog.Logger sl.Info()", "str", "some string")
```

Look `xlog_test.go` for more examples.

