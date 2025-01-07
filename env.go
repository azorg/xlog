// File: "env.go"

package xlog

import (
	"os"
	"strconv"
	"strings"
)

// Default prefix
const DEFAULT_PREFIX = "LOG_"

// String to bool converter
//
//	true:  true, True, yes, YES, on, 1, 2
//	false: false, FALSE, no, Off, 0, "Abra-Cadabra"
func StringToBool(s string) bool {
	s = strings.ToLower(s)
	if s == "true" || s == "on" || s == "yes" {
		return true
	} else if s == "false" || s == "off" || s == "no" {
		return false
	}
	i, err := strconv.Atoi(s)
	if err == nil {
		return i != 0
	}
	return false // by default
}

// String to int converter
func StringToInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}

// Add settings from environment variables
func Env(conf *Conf, prefixOpt ...string) {
	prefix := DEFAULT_PREFIX
	if len(prefixOpt) != 0 {
		prefix = prefixOpt[0]
	}
	if v := os.Getenv(prefix + "FILE"); v != "" {
		conf.File = v
	}
	if v := os.Getenv(prefix + "FILE_MODE"); v != "" {
		conf.FileMode = v
	}
	if v := os.Getenv(prefix + "LEVEL"); v != "" {
		conf.Level = v
	}
	if v := os.Getenv(prefix + "SLOG"); v != "" {
		conf.Slog = StringToBool(v)
		if conf.Slog {
			conf.Tint = false
			conf.JSON = false
		}
	}
	if v := os.Getenv(prefix + "JSON"); v != "" {
		conf.JSON = StringToBool(v)
		if conf.JSON {
			conf.Tint = false
		}
	}
	if v := os.Getenv(prefix + "TINT"); v != "" {
		conf.Tint = StringToBool(v)
	}
	if v := os.Getenv(prefix + "TIME"); v != "" {
		conf.Time = StringToBool(v)
	}
	if v := os.Getenv(prefix + "TIME_US"); v != "" {
		conf.TimeUS = StringToBool(v)
	}
	if v := os.Getenv(prefix + "TIME_TINT"); v != "" {
		conf.TimeTint = v
	}
	if v := os.Getenv(prefix + "SRC"); v != "" {
		conf.Src = StringToBool(v)
	}
	if v := os.Getenv(prefix + "SRC_LONG"); v != "" {
		conf.SrcLong = StringToBool(v)
	}
	if v := os.Getenv(prefix + "SRC_FUNC"); v != "" {
		conf.SrcFunc = StringToBool(v)
	}
	if v := os.Getenv(prefix + "NO_EXT"); v != "" {
		conf.NoExt = StringToBool(v)
	}
	if v := os.Getenv(prefix + "NO_LEVEL"); v != "" {
		conf.NoLevel = StringToBool(v)
	}
	if v := os.Getenv(prefix + "NO_COLOR"); v != "" {
		conf.NoColor = StringToBool(v)
	}
	if v := os.Getenv(prefix + "PREFIX"); v != "" {
		conf.Prefix = v
	}
	if v := os.Getenv(prefix + "ADD_KEY"); v != "" {
		conf.AddKey = v
	}
	if v := os.Getenv(prefix + "ADD_VALUE"); v != "" {
		conf.AddValue = v
	}
	if v := os.Getenv(prefix + "ROTATE"); v != "" {
		conf.Rotate.Enable = StringToBool(v)
	}
	if v := os.Getenv(prefix + "ROTATE_MAX_SIZE"); v != "" {
		conf.Rotate.MaxSize = StringToInt(v)
	}
	if v := os.Getenv(prefix + "ROTATE_MAX_AGE"); v != "" {
		conf.Rotate.MaxAge = StringToInt(v)
	}
	if v := os.Getenv(prefix + "ROTATE_MAX_BACKUPS"); v != "" {
		conf.Rotate.MaxBackups = StringToInt(v)
	}
	if v := os.Getenv(prefix + "ROTATE_LOCAL_TIME"); v != "" {
		conf.Rotate.LocalTime = StringToBool(v)
	}
	if v := os.Getenv(prefix + "ROTATE_COMPRESS"); v != "" {
		conf.Rotate.Compress = StringToBool(v)
	}
}

// EOF: "env.go"
