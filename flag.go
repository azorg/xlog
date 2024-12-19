// File: "flag.go"

package xlog

import "flag"

// Command line logger option structure
type Opt struct {
	Level   string // -log <level>
	File    string // -lfile <file>
	SLog    bool   // -slog
	JLog    bool   // -jlog
	TLog    bool   // -tlog
	Src     bool   // -lsrc
	NoSrc   bool   // -lnosrc
	Pkg     bool   // -lpkg
	NoPkg   bool   // -lnopkg
	Func    bool   // -lfunc
	NoFunc  bool   // -lnofunc
	Ext     bool   // -lext
	NoExt   bool   // -lnoext
	Time    bool   // -ltime
	NoTime  bool   // -lnotime
	TimeFmt string // -ltimefmt <fmt>
	NoLevel bool   // -lnolevel
	Color   bool   // -lcolor
	NoColor bool   // -lnocolor
}

// Setup command line logger options
// Usage:
//
//		-log <level>       - log level (flood/trace/debug/info/notice/warm/error/critical)
//		-lfile <file>      - log file path or stdout/stderr
//		-slog              - use structured text logger (slog)
//		-jlog              - use structured JSON logger (slog)
//		-tlog              - use tinted (colorized) logger (tint)
//		-lsrc|-lnosrc      - force on/off log source file name and line number
//		-lpkg|-lnopkg      - force on/off log source directory/file name and line number
//	 -lfunc|-lnofunc    - force on/off log function name
//		-lext|-lnoext      - force enable/disable remove ".go" extension from source file name
//		-ltime|-lnotime    - force on/off timestamp
//		-ltimefmt <format> - override log time format (e.g. 15:04:05.999 or TimeOnly)
//		-lnolevel          - disable log level tag (~level="INFO")
//		-lcolor|-lnocolor  - force enable/disable tinted colors
func NewOpt() *Opt {
	opt := &Opt{}
	flag.StringVar(&opt.Level, "log", "", "override log level (flood/trace/debug/info/warm/error/fatal)")
	flag.StringVar(&opt.File, "lfile", "", "log file path or stdout/stderr")
	flag.BoolVar(&opt.SLog, "slog", false, "use structured text logger (slog)")
	flag.BoolVar(&opt.JLog, "jlog", false, "use structured JSON logger (slog)")
	flag.BoolVar(&opt.TLog, "tlog", false, "use tinted (colorized) logger (tint)")
	flag.BoolVar(&opt.Src, "lsrc", false, "force log source file name and line number")
	flag.BoolVar(&opt.NoSrc, "lnosrc", false, "force off source file name and line number")
	flag.BoolVar(&opt.Pkg, "lpkg", false, "force log source directory/file name and line number")
	flag.BoolVar(&opt.NoPkg, "lnopkg", false, "force off source directory/file name and line number")
	flag.BoolVar(&opt.Func, "lfunc", true, "force enable functions name")
	flag.BoolVar(&opt.NoFunc, "lnofunc", false, "force disable functions name")
	flag.BoolVar(&opt.Ext, "lext", true, "force enable remove '.go' extension from source file name")
	flag.BoolVar(&opt.NoExt, "lnoext", false, "force disable remove '.go' extension from source file name")
	flag.BoolVar(&opt.Time, "ltime", false, "force add timestamp to log")
	flag.BoolVar(&opt.NoTime, "lnotime", false, "force off timestamp")
	flag.StringVar(&opt.TimeFmt, "ltimefmt", "", "override log time format (e.g. 15:04:05.999 or TimeOnly)")
	flag.BoolVar(&opt.NoLevel, "lnolevel", false, `disable log level tag (~level="INFO")`)
	flag.BoolVar(&opt.Color, "lcolor", false, "force enable tinted colors")
	flag.BoolVar(&opt.NoColor, "lnocolor", false, "force disable tinted colors")
	return opt
}

// Add parsed command line options to logger config
func AddOpt(opt *Opt, conf *Conf) {
	if opt.Level != "" {
		conf.Level = opt.Level
	}
	if opt.File != "" {
		conf.File = opt.File
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
	} else if opt.NoSrc {
		conf.Src = false
	}
	if opt.Pkg {
		conf.Src = true
		conf.SrcLong = true
	} else if opt.NoPkg {
		conf.SrcLong = false
	}
	if opt.Func {
		conf.Src = true
		conf.SrcFunc = true
	} else if opt.NoFunc {
		conf.SrcFunc = false
	}
	if opt.Ext {
		conf.Src = true
		conf.NoExt = false
	} else if opt.NoExt {
		conf.NoExt = true
	}
	if opt.Time {
		conf.Time = true
	} else if opt.NoTime {
		conf.Time = false
		conf.TimeTint = TIME_OFF
	}
	if opt.TimeFmt != "" {
		conf.Time = true
		conf.TimeTint = opt.TimeFmt
	}
	if opt.NoLevel {
		conf.NoLevel = true
	}
	if opt.NoColor {
		conf.NoColor = true
	} else if opt.Color {
		conf.NoColor = false
	}
}

// EOF: "flag.go"
