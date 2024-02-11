// File: "flag.go"

package xlog

import "flag"

// Command line logger option structure
type Opt struct {
	Level   string // -log <level>
	SLog    bool   // -slog
	JLog    bool   // -jlog
	TLog    bool   // -tlog
	Src     bool   // -lsrc
	NoSrc   bool   // -lnosrc
	Pkg     bool   // -lpkg
	NoPkg   bool   // -lnopkg
	Time    bool   // -ltime
	NoTime  bool   // -lnotime
	TimeFmt string // -ltimefmt <fmt>
}

// Setup command line logger options
// Usage:
//
//	-log  <level> - Log level (trace/debug/info/notice/warm/error/fatal)
//	-slog         - Use structured text logger (slog)
//	-jlog         - Use structured JSON logger (slog)
//	-tlog         - Use tinted (colorized) logger (tint)
//	-lsrc         - Force log source file name and line number
//	-lnosrc       - Force log source file name and line number
//	-lpkg         - Force log source directory/file name and line number
//	-lnopkg       - Force off source directory/file name and line number
//	-ltime        - Force add timestamp to log
//	-lnotime      - Force off timestamp
//	-ltimefmt     - Timestamp format
func NewOpt() *Opt {
	opt := &Opt{}
	flag.StringVar(&opt.Level, "log", "", "Override log level (flood/trace/debug/info/warm/error/fatal)")
	flag.BoolVar(&opt.SLog, "slog", false, "Use structured text logger (slog)")
	flag.BoolVar(&opt.JLog, "jlog", false, "Use structured JSON logger (slog)")
	flag.BoolVar(&opt.TLog, "tlog", false, "Use tinted (colorized) logger (tint)")
	flag.BoolVar(&opt.Src, "lsrc", false, "Force log source file name and line number")
	flag.BoolVar(&opt.NoSrc, "lnosrc", false, "Force off source file name and line number")
	flag.BoolVar(&opt.Pkg, "lpkg", false, "Force log source directory/file name and line number")
	flag.BoolVar(&opt.NoPkg, "lnopkg", false, "Force off source directory/file name and line number")
	flag.BoolVar(&opt.Time, "ltime", false, "Force add timestamp to log")
	flag.BoolVar(&opt.NoTime, "lnotime", false, "Force off timestamp")
	flag.StringVar(&opt.TimeFmt, "ltimefmt", "", "Override log timestamp format (e.g. 15:04:05.999)")
	return opt
}

// Add parsed command line options to logger config
func AddOpt(opt *Opt, conf *Conf) {
	if opt.Level != "" {
		conf.Level = opt.Level
	}
	if opt.SLog {
		conf.Slog = true
		conf.Tint = false
		conf.JSON = false
	}
	if opt.JLog {
		conf.JSON = true
		conf.Tint = false
	}
	if opt.TLog {
		conf.Tint = true
	}
	if opt.Src {
		conf.Src = true
	}
	if opt.NoSrc {
		conf.Src = false
	}
	if opt.Pkg {
		conf.Src = true
		conf.SrcLong = true
	}
	if opt.NoPkg {
		conf.SrcLong = false
	}
	if opt.Time {
		conf.Time = true
	}
	if opt.NoTime {
		conf.Time = false
	}
	if opt.TimeFmt != "" {
		conf.Time = true
		conf.TimeTint = opt.TimeFmt
	}
}

// EOF: "flag.go"
