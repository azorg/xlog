// File: "flag.go"

package xlog

import "flag"

// Command line logger options
type Opt struct {
	LogLvl  string // log level (trace/debug/info/warn/error/fatal)
	SLog    bool   // use structured text loger (slog)
	JLog    bool   // use structured JSON loger (slog)
	TLog    bool   // use tinted (colorized) logger (tint)
	LogSrc  bool   // force log source file name and line number
	LogTime bool   // force add time to log
}

// Setup command line logger options
// Usage:
//
//	-log  <level> - Log level (trace/debug/info/notice/warm/error/fatal)
//	-slog         - Use structured text logger (slog)
//	-jlog         - Use structured JSON logger (slog)
//	-tlog         - Use tinted (colorized) logger (tint)
//	-lsrc         - Force log source file name and line number
//	-ltime        - Force add time to log
func NewOpt() *Opt {
	opt := &Opt{}
	flag.StringVar(&opt.LogLvl, "log", "", "Override log level (trace/debug/info/warm/error/fatal)")
	flag.BoolVar(&opt.SLog, "slog", false, "Use structured text logger (slog)")
	flag.BoolVar(&opt.JLog, "jlog", false, "Use structured JSON logger (slog)")
	flag.BoolVar(&opt.TLog, "tlog", false, "Use tinted (colorized) logger (tint)")
	flag.BoolVar(&opt.LogSrc, "lsrc", false, "Force log source file name and line number")
	flag.BoolVar(&opt.LogTime, "ltime", false, "Force add time to log")
	return opt
}

// Add parsed command line options to logger config
func AddOpt(opt *Opt, conf *Conf) {
	if opt.LogLvl != "" {
		conf.Level = Lvl(opt.LogLvl)
	}
	if opt.SLog {
		conf.Slog = true
	}
	if opt.JLog {
		conf.JSON = true
	}
	if opt.TLog {
		conf.Tint = true
	}
	if opt.LogSrc {
		conf.Src = true
	}
	if opt.LogTime {
		conf.Time = true
	}
}

// EOF: "flag.go"
