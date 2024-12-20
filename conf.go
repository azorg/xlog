// File: "conf.go"

package xlog

// Logger configure structure
type Conf struct {
	File     string `json:"file"`      // log file path OR stdout/stderr
	FileMode string `json:"file-mode"` // log file mode (if File is not stdout/stderr)
	Level    string `json:"level"`     // log level (trace/debug/info/warn/error/fatal/silent)
	Slog     bool   `json:"slog"`      // use slog instead standard log (slog.TextHandler)
	JSON     bool   `json:"json"`      // use JSON log (slog.JSONHandler)
	Tint     bool   `json:"tint"`      // use tinted (colorized) log (xlog.TintHandler)
	Time     bool   `json:"time"`      // add timestamp
	TimeUS   bool   `json:"time-us"`   // use timestamp in microseconds
	TimeTint string `json:"time-tint"` // tinted log time format (like time.Kitchen, time.DateTime)
	Src      bool   `json:"src"   `    // log file name and line number
	SrcLong  bool   `json:"src-long"`  // log long file path (directory + file name)
	SrcFunc  bool   `json:"src-func"`  // add function name to log
	NoExt    bool   `json:"no-ext"`    // remove ".go" extension from file name
	NoLevel  bool   `json:"no-level"`  // don't print log level tag to log (~level="INFO")
	NoColor  bool   `json:"no-color"`  // disable tinted colors (only if Tint=true)
	Prefix   string `json:"preifix"`   // add prefix to standard log (log=false)
	AddKey   string `json:"add-key"`   // add key to structured log (Slog=true)
	AddValue string `json:"add-value"` // add value to structured log (Slog=true

	// Log rotate options or nil
	Rotate *RotateOpt `json:"rotate"`
}

// Log rotate options (delivered from lumberjack)
type RotateOpt struct {
	// Enable log rotation
	Enable bool `json:"enable"`

	// Maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `json:"max-size"`

	// Maximum number of days to retain old log files based on the
	// timestamp encoded in their filename. Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `json:"max-age"`

	// Maximum number of old log files to retain. The default/ is to retain
	// all old log files (though MaxAge may still cause them to get deleted)
	MaxBackups int `json:"max-backups"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool `json:"local-time"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `json:"compress"`
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
		SrcFunc:  SRC_FUNC,
		NoExt:    NO_EXT,
		NoLevel:  NO_LEVEL,
		NoColor:  NO_COLOR,
		Prefix:   PREFIX,
		AddKey:   ADD_KEY,
		AddValue: ADD_VALUE,
		Rotate: &RotateOpt{
			Enable:     ROTATE_ENABLE,
			MaxSize:    ROTATE_MAX_SIZE,
			MaxAge:     ROTATE_MAX_AGE,
			MaxBackups: ROTATE_MAX_BACKUPS,
			LocalTime:  ROTATE_LOCAL_TIME,
			Compress:   ROTATE_COMPRESS,
		},
	}
}

// EOF: "conf.go"
