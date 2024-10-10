// File: "time.go"

package xlog

import (
	"strings"
	"time"
)

// Time formats
const (
	// Time OFF
	TIME_OFF = ""

	// Default time format of standart logger
	STD_TIME = "2006/01/02 15:04:05"
	STD_LOG  = "2006-01-02 15:04:05"

	// Default time format of standart logger + milliseconds
	STD_TIME_MS = "2006/01/02 15:04:05.999"
	STD_LOG_MS  = "2006-01-02 15:04:05.999"

	// Default time format of standart logger + microseconds
	STD_TIME_US = "2006/01/02 15:04:05.999999"
	STD_LOG_US  = "2006-01-02 15:04:05.999999"

	// RFC3339 time format + nanoseconds (slog.TextHandler by default)
	RFC3339Nano = time.RFC3339Nano // "2006-01-02T15:04:05.999999999Z07:00"

	// RFC3339 time format + microseconds
	RFC3339Micro = "2006-01-02T15:04:05.999999Z07:00"

	// RFC3339 time format + milliseconds
	RFC3339Milli = "2006-01-02T15:04:05.999Z07:00"

	// Time only format + microseconds
	TimeOnlyMicro = "15:04:05.999999"

	// Time only format + milliseconds
	TimeOnlyMilli = "15:04:05.999"

	// Time format for file names (no spaces, no ":", sorted by date/time)
	FILE_TIME_FORMAT = "2006-01-02_15.04.05"

	// Compromise time format (no spaces, no ":")
	COMPROMISE_TIME_FORMAT_DS = "2006-01-02_15.04.05.9"
	COMPROMISE_TIME_FORMAT    = "2006-01-02_15.04.05.999"
	COMPROMISE_TIME_FORMAT_US = "2006-01-02_15.04.05.999999"
	COMPROMISE_TIME_FORMAT_NS = "2006-01-02_15.04.05.999999999"

	// Digital clock
	CLOCK_TIME_FORMAT = "15:04"

	// Default (recomented) time format wuth milliseconds
	DEFAULT_TIME_FORMAT = STD_TIME_MS

	// Default (recomented) time format witn microseconds
	DEFAULT_TIME_FORMAT_US = STD_TIME_US
)

// Time format aliases
var timeFormat = map[string]string{
	// From "time" Go package:
	"Layout":      time.Layout,      // "01/02 03:04:05PM '06 -0700"
	"ANSIC":       time.ANSIC,       // "Mon Jan _2 15:04:05 2006"
	"UnixDate":    time.UnixDate,    // "Mon Jan _2 15:04:05 MST 2006"
	"RubyDate":    time.RubyDate,    // "Mon Jan 02 15:04:05 -0700 2006"
	"RFC822":      time.RFC822,      // "02 Jan 06 15:04 MST"
	"RFC822Z":     time.RFC822Z,     // "02 Jan 06 15:04 -0700"
	"RFC850":      time.RFC850,      // "Monday, 02-Jan-06 15:04:05 MST"
	"RFC1123":     time.RFC1123,     // "Mon, 02 Jan 2006 15:04:05 MST"
	"RFC1123Z":    time.RFC1123Z,    // "Mon, 02 Jan 2006 15:04:05 -0700"
	"RFC3339":     time.RFC3339,     // "2006-01-02T15:04:05Z07:00"
	"RFC3339Nano": time.RFC3339Nano, // "2006-01-02T15:04:05.999999999Z07:00"
	"Kitchen":     time.Kitchen,     // "3:04PM"
	"Stamp":       time.Stamp,       // "Jan _2 15:04:05"
	"StampMilli":  time.StampMilli,  // "Jan _2 15:04:05.000"
	"StampMicro":  time.StampMicro,  // "Jan _2 15:04:05.000000"
	"StampNano":   time.StampNano,   // "Jan _2 15:04:05.000000000"
	"DateTime":    time.DateTime,    // "2006-01-02 15:04:05"
	"DateOnly":    time.DateOnly,    // "2006-01-02"
	"TimeOnly":    time.TimeOnly,    // "15:04:05"

	// xlog ideas:
	"StdTime":       STD_TIME,                  // "2006/01/02 15:04:05"
	"StdTimeMilli":  STD_TIME_MS,               // "2006/01/02 15:04:05.999"
	"StdTimeMicro":  STD_TIME_US,               // "2006/01/02 15:04:05.999999"
	"StdLog":        STD_LOG,                   // "2006-01-02 15:04:05"
	"StdLogMilli":   STD_LOG_MS,                // "2006-01-02 15:04:05.999"
	"StdLogMicro":   STD_LOG_US,                // "2006-01-02 15:04:05.999999"
	"RFC3339Micro":  RFC3339Micro,              // "2006-01-02T15:04:05.999999Z07:00"
	"RFC3339Milli":  RFC3339Milli,              // "2006-01-02T15:04:05.999Z07:00"
	"TimeOnlyMicro": TimeOnlyMicro,             // "15:04:05.999999"
	"TimeOnlyMilli": TimeOnlyMilli,             // "15:04:05.999"
	"Default":       STD_TIME_MS,               // "2006/01/02 15:04:05.999"
	"DefaultMicro":  STD_TIME_US,               // "2006/01/02 15:04:05.999999
	"Office":        CLOCK_TIME_FORMAT,         // "15:04" like "Kitchen"
	"File":          FILE_TIME_FORMAT,          // "2006-01-02_15.04.05"
	"Home":          COMPROMISE_TIME_FORMAT_DS, // "2006-01-02_15.04.05.9"
	"Lab":           COMPROMISE_TIME_FORMAT,    // "2006-01-02_15.04.05.999"
	"Science":       COMPROMISE_TIME_FORMAT_US, // "2006-01-02_15.04.05.999999"
	"Space":         COMPROMISE_TIME_FORMAT_NS, // "2006-01-02_15.04.05.999999999"

	// Off time formats
	"off":     TIME_OFF,
	"no":      TIME_OFF,
	"false":   TIME_OFF,
	"disable": TIME_OFF,
	"0":       TIME_OFF,
}

// Convert time formats to lower case
func init() {
	m := make(map[string]string)
	for k, v := range timeFormat {
		m[strings.ToLower(k)] = v
	}
	timeFormat = m
}

// Return time format by alias
func TimeFormat(alias string) (format string, ok bool) {
	format, ok = timeFormat[strings.ToLower(alias)]
	if ok {
		return format, true // alias found
	}
	return alias, false // return as-is
}

// EOF: "time.go"
