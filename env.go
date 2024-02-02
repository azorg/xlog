// File: "env.go"

package xlog

import (
	"os"
	"strconv"
	"strings"
)

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

// Add settings from eviroment variables
func Env(conf *Conf) {
	if v := os.Getenv("XLOG_FILE"); v != "" {
		conf.File = v
	}
	if v := os.Getenv("XLOG_FILE_MODE"); v != "" {
		conf.FileMode = v
	}
	if v := os.Getenv("XLOG_LEVEL"); v != "" {
		conf.Level = v
	}
	if v := os.Getenv("XLOG_SLOG"); v != "" {
		conf.Slog = StringToBool(v)
	}
	if v := os.Getenv("XLOG_JSON"); v != "" {
		conf.JSON = StringToBool(v)
	}
	if v := os.Getenv("XLOG_TINT"); v != "" {
		conf.Tint = StringToBool(v)
	}
	if v := os.Getenv("XLOG_TIME"); v != "" {
		conf.Time = StringToBool(v)
	}
	if v := os.Getenv("XLOG_TIME_US"); v != "" {
		conf.TimeUS = StringToBool(v)
	}
	if v := os.Getenv("XLOG_TIME_TINT"); v != "" {
		conf.TimeTint = v
	}
	if v := os.Getenv("XLOG_SRC"); v != "" {
		conf.Src = StringToBool(v)
	}
	if v := os.Getenv("XLOG_SRC_LONG"); v != "" {
		conf.SrcLong = StringToBool(v)
	}
	if v := os.Getenv("XLOG_NO_LEVEL"); v != "" {
		conf.NoLevel = StringToBool(v)
	}
	if v := os.Getenv("XLOG_NO_COLOR"); v != "" {
		conf.NoColor = StringToBool(v)
	}
	if v := os.Getenv("XLOG_PREFIX"); v != "" {
		conf.Prefix = v
	}
	if v := os.Getenv("XLOG_ADD_KEY"); v != "" {
		conf.AddKey = v
	}
	if v := os.Getenv("XLOG_ADD_VALUE"); v != "" {
		conf.AddValue = v
	}
}

// EOF: "env.go"
