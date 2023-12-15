// File: "conf.go"

package xlog

// Logger configure structure
type Conf struct {
	File     string `json:"file"`      // log file path OR stdout/stderr
	FileMode string `json:"file-mode"` // log file mode (if File is not stdout/stderr)
	Level    string `json:"level"`     // log level (trace/debug/info/warn/error/fatal/silent)
	Slog     bool   `json:"slog"`      // use slog insted standart log (slog.TextHandler)
	JSON     bool   `json:"json"`      // use JSON log (slog.JSONHandler)
	Tint     bool   `json:"tint"`      // use tinted (colorized) log (xlog.TintHandler)
	Time     bool   `json:"time"`      // add timestamp
	TimeUS   bool   `json:"time-us"`   // use timestamp in microseconds
	TimeTint string `json:"time-tint"` // tinted log time format (like time.Kitchen, time.DateTime)
	Src      bool   `json:"src"   `    // log file name and line number
	SrcLong  bool   `json:"src-long"`  // log long file path (directory + file name)
	NoLevel  bool   `json:"no-level"`  // don't print log level tag to log (~level="INFO")
	NoColor  bool   `json:"no-color"`  // disable tinted colors (only if Tint=true)
	Prefix   string `json:"preifix"`   // add prefix to standart log (log=false)
	AddKey   string `json:"add-key"`   // add key to structured log (Slog=true)
	AddValue string `json:"add-value"` // add value to structured log (Slog=true
}

// Create default logger structure
func NewConf() Conf {
	return Conf{
		File:     FILE,
		FileMode: FILE_MODE,
		Level:    LEVEL,
		Slog:     SLOG,
		JSON:     JSON,
		Tint:     TINT,
		Time:     TIME,
		TimeUS:   TIME_US,
		TimeTint: TIME_TINT,
		Src:      SRC,
		SrcLong:  SRC_LONG,
		NoLevel:  NO_LEVEL,
		NoColor:  NO_COLOR,
		Prefix:   PREFIX,
		AddKey:   ADD_KEY,
		AddValue: ADD_VALUE,
	}
}

// EOF: "conf.go"
