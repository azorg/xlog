// File: "time.go"

package xlog

import (
	"strings"
	"time"
)

// Форматы временной метки для TintHanler'а
const (
	// Метка времени как у стандартного логгера в Go
	StdTime  = "2006/01/02 15:04:05"
	DateTime = time.DateTime // "2006-01-02 15:04:05"

	// Стандартная метка времени с миллисекундами
	StdTimeMilli  = "2006/01/02 15:04:05.999"
	DateTimeMilli = "2006-01-02 15:04:05.999"

	// Стандартная метка времени с микросекундами
	StdTimeMicro  = "2006/01/02 15:04:05.999999"
	DateTimeMicro = "2006-01-02 15:04:05.999999"

	// Формат RFC3339 с наносекундами (slog.TextHandler использует по умолчанию)
	RFC3339Nano = time.RFC3339Nano // "2006-01-02T15:04:05.999999999Z07:00"

	// Формат RFC3339 с микросекундами
	RFC3339Micro = "2006-01-02T15:04:05.999999Z07:00"

	// Формат RFC3339 с миллисекундами
	RFC3339Milli = "2006-01-02T15:04:05.999Z07:00"

	// Только время с миллисекундами
	TimeOnlyMilli = "15:04:05.999"

	// Только время с микросекундами
	TimeOnlyMicro = "15:04:05.999999"

	// Формат простых цифровых часов (по аналогии с time.Kitchen)
	Office = "15:04"

	// Форматы времени дата+время без пробелов (с долями секунд)
	File    = "2006-01-02_15.04.05"
	Home    = "2006-01-02_15.04.05.9"
	Lab     = "2006-01-02_15.04.05.999"
	Science = "2006-01-02_15.04.05.999999"
	Space   = "2006-01-02_15.04.05.999999999"

	// Формат по умолчанию (с миллисекундами)
	defaultTime = StdTimeMilli

	// Формат по умолчанию (с микросекундами)
	defaultTimeMicro = StdTimeMicro

	// Метка времени отключена
	timeOff = ""
)

// Таблица идентификаторов (псевдонимов) форматов временных меток
// используемых для настройки TintHanler'а
var timeFormats = map[string]string{
	// Стандартные форматы из библиотеки Go начиная с версии go1.20
	"Layout":        time.Layout,      // "01/02 03:04:05PM '06 -0700"
	"ANSIC":         time.ANSIC,       // "Mon Jan _2 15:04:05 2006"
	"UnixDate":      time.UnixDate,    // "Mon Jan _2 15:04:05 MST 2006"
	"RubyDate":      time.RubyDate,    // "Mon Jan 02 15:04:05 -0700 2006"
	"RFC822":        time.RFC822,      // "02 Jan 06 15:04 MST"
	"RFC822Z":       time.RFC822Z,     // "02 Jan 06 15:04 -0700"
	"RFC850":        time.RFC850,      // "Monday, 02-Jan-06 15:04:05 MST"
	"RFC1123":       time.RFC1123,     // "Mon, 02 Jan 2006 15:04:05 MST"
	"RFC1123Z":      time.RFC1123Z,    // "Mon, 02 Jan 2006 15:04:05 -0700"
	"RFC3339":       time.RFC3339,     // "2006-01-02T15:04:05Z07:00"
	"RFC3339Nano":   time.RFC3339Nano, // "2006-01-02T15:04:05.999999999Z07:00"
	"Kitchen":       time.Kitchen,     // "3:04PM"
	"Stamp":         time.Stamp,       // "Jan _2 15:04:05"
	"StampMilli":    time.StampMilli,  // "Jan _2 15:04:05.000"
	"StampMicro":    time.StampMicro,  // "Jan _2 15:04:05.000000"
	"StampNano":     time.StampNano,   // "Jan _2 15:04:05.000000000"
	"DateTime":      time.DateTime,    // "2006-01-02 15:04:05"
	"DateTimeMilli": DateTimeMilli,    // "2006-01-02 15:04:05.999"
	"DateTimeMicro": DateTimeMicro,    // "2006-01-02 15:04:05.999999"
	"DateOnly":      time.DateOnly,    // "2006-01-02"
	"TimeOnly":      time.TimeOnly,    // "15:04:05"

	// Дополнительные форматы времени, предоставляемые пакетом xlog
	"StdTime":       StdTime,       // "2006/01/02 15:04:05"
	"StdTimeMilli":  StdTimeMilli,  // "2006/01/02 15:04:05.999"
	"StdTimeMicro":  StdTimeMicro,  // "2006/01/02 15:04:05.999999"
	"RFC3339Micro":  RFC3339Micro,  // "2006-01-02T15:04:05.999999Z07:00"
	"RFC3339Milli":  RFC3339Milli,  // "2006-01-02T15:04:05.999Z07:00"
	"TimeOnlyMicro": TimeOnlyMicro, // "15:04:05.999999"
	"TimeOnlyMilli": TimeOnlyMilli, // "15:04:05.999"
	"Default":       StdTimeMilli,  // "2006/01/02 15:04:05.999"
	"DefaultMicro":  StdTimeMicro,  // "2006/01/02 15:04:05.999999
	"File":          File,          // "2006-01-02_15.04.05"
	"Office":        Office,        // "15:04" по аналогии с "Kitchen" в Go
	"Home":          Home,          // "2006-01-02_15.04.05.9"
	"Lab":           Lab,           // "2006-01-02_15.04.05.999"
	"Science":       Science,       // "2006-01-02_15.04.05.999999"
	"Space":         Space,         // "2006-01-02_15.04.05.999999999"

	// Идентификаторы, которые подразумевают отключение метки времени в журнале
	"off":     timeOff,
	"no":      timeOff,
	"false":   timeOff,
	"disable": timeOff,
	"0":       timeOff,
}

// Преобразовать идентификаторы временных форматов к нижнему регистру
func init() {
	m := make(map[string]string)
	for k, v := range timeFormats {
		m[strings.ToLower(k)] = v
	}
	timeFormats = m
}

// TimeFormat возвращает строку форматирования временной метки
// в соответствии с заданной строкой идентификатором.
// Идентификатор (псевдоним) обрабатывается без учета регистра.
// Если псевдоним не найден, то входная строка используется
// как формат времени в Go нотации (аля "2006-01-02 15:04:05").
//
// На выходе ok=true, если псевдоним найден.
// Функция используется только для Tinted хендлера.
// Допустимы следующие псевдонимы:
//
//	Стандартные форматы из библиотеки Go начиная с версии go1.20:
//	"Layout":        time.Layout       "01/02 03:04:05PM '06 -0700"
//	"ANSIC":         time.ANSIC        "Mon Jan _2 15:04:05 2006"
//	"UnixDate":      time.UnixDate     "Mon Jan _2 15:04:05 MST 2006"
//	"RubyDate":      time.RubyDate     "Mon Jan 02 15:04:05 -0700 2006"
//	"RFC822":        time.RFC822       "02 Jan 06 15:04 MST"
//	"RFC822Z":       time.RFC822Z      "02 Jan 06 15:04 -0700"
//	"RFC850":        time.RFC850       "Monday, 02-Jan-06 15:04:05 MST"
//	"RFC1123":       time.RFC1123      "Mon, 02 Jan 2006 15:04:05 MST"
//	"RFC1123Z":      time.RFC1123Z     "Mon, 02 Jan 2006 15:04:05 -0700"
//	"RFC3339":       time.RFC3339      "2006-01-02T15:04:05Z07:00"
//	"RFC3339Nano":   time.RFC3339Nano  "2006-01-02T15:04:05.999999999Z07:00"
//	"Kitchen":       time.Kitchen      "3:04PM"
//	"Stamp":         time.Stamp        "Jan _2 15:04:05"
//	"StampMilli":    time.StampMilli   "Jan _2 15:04:05.000"
//	"StampMicro":    time.StampMicro   "Jan _2 15:04:05.000000"
//	"StampNano":     time.StampNano    "Jan _2 15:04:05.000000000"
//	"DateTime":      time.DateTime     "2006-01-02 15:04:05"
//	"DateTimeMilli": DateTimeMilli     "2006-01-02 15:04:05.999"
//	"DateTimeMicro": DateTimeMicro     "2006-01-02 15:04:05.999999"
//	"DateOnly":      time.DateOnly     "2006-01-02"
//	"TimeOnly":      time.TimeOnly     "15:04:05"
//
//	Дополнительные форматы времени, предоставляемые пакетом xlog:
//	"StdTime":       StdTime        "2006/01/02 15:04:05"
//	"StdTimeMilli":  StdTimeMilli   "2006/01/02 15:04:05.999"
//	"StdTimeMicro":  StdTimeMicro   "2006/01/02 15:04:05.999999"
//	"RFC3339Micro":  RFC3339Micro   "2006-01-02T15:04:05.999999Z07:00"
//	"RFC3339Milli":  RFC3339Milli   "2006-01-02T15:04:05.999Z07:00"
//	"TimeOnlyMicro": TimeOnlyMicro  "15:04:05.999999"
//	"TimeOnlyMilli": TimeOnlyMilli  "15:04:05.999"
//	"Default":       StdTimeMilli   "2006/01/02 15:04:05.999"
//	"DefaultMicro":  StdTimeMicro   "2006/01/02 15:04:05.999999
//	"File":          File           "2006-01-02_15.04.05"
//	"Office":        Office         "15:04" по аналогии с "Kitchen" в Go
//	"Home":          Home           "2006-01-02_15.04.05.9"
//	"Lab":           Lab            "2006-01-02_15.04.05.999"
//	"Science":       Science        "2006-01-02_15.04.05.999999"
//	"Space":         Space          "2006-01-02_15.04.05.999999999"
//
//	Идентификаторы, которые подразумевают отключение метки времени в журнале:
//	"" (пустая строка), "off", "no", "false", "disable", "0"
func TimeFormat(alias string) (format string, ok bool) {
	format, ok = timeFormats[strings.ToLower(alias)]
	if ok {
		return format, true // идентификатор формата найден
	}
	return alias, false // вернуть как есть
}

// EOF: "time.go"
