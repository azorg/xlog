// File: "env.go"

package xlog

import (
	"os"
	"strconv"
	"strings"
)

// Префикс для переменных окружения по умолчанию
const DefaultEnvPrefix = "LOG_"

// StringToBool преобразует строку в булево значение.
// Допустимы следующие варианты входных значений:
//
//	true:  "true", "True", "yes", "YES", "on", "1", "2", "99"
//	false: "false", "FALSE", "no", "Off", "0", "Abra-Cadabra"
//
// В случае ошибки (по умолчанию) возвращается false.
// Функция используется при обработке переменных окружения и флагов.
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

// StrintToInt преобразует строку к целому числу.
// В случае ошибки возвращается 0.
// Функция используется при обработке переменных окружения и флагов.
func StringToInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}

// Env - обогащает структуру конфигурации логгера данными
// из переменных окружения.
//
//	conf - обогащаемая структура конфигурации логгера
//	prefixOpt - опциональный префикс (по умолчанию "LOG_")
//
// Анализируются следующие переменные окружения, соответствующие
// (возможно с инверсией) полям структуры Conf:
//
//	LOG_LEVEL       (string/int: "debug", "trace", "error", "0", "-20"...)
//	LOG_PIPE        (string: "stdout", "stderr", "null")
//	LOG_FILE        (string: ~"logs/app.log")
//	LOG_FILE_MODE   (string: ~"0640")
//	LOG_FORMAT      (string: "json", "logfmt", "tinted", "default")
//	LOG_GOID        (bool)
//	LOG_ID          (bool)
//	LOG_SUM         (bool)
//	LOG_SUM_FULL    (bool)
//	LOG_SUM_CHAIN   (bool)
//	LOG_SUM_ALONE   (bool)
//	LOG_TIME        (bool)
//	LOG_TIME_LOCAL  (bool)
//	LOG_TIME_MICRO  (bool)
//	LOG_TIME_FORMAT (string: "default", "lab", "15.04.05"...)
//	LOG_SRC         (bool)
//	LOG_SRC_PKG     (bool)
//	LOG_SRC_FUNC    (bool)
//	LOG_SRC_EXT     (bool)
//	LOG_COLOR       (bool)
//	LOG_LEVEL_OFF   (bool)
//	LOG_ROTATE      (bool)
//	LOG_ROTATE_MAX_SIZE    (int: мегабайт)
//	LOG_ROTATE_MAX_AGE     (int: суток)
//	LOG_ROTATE_MAX_BACKUPS (int: число файлов)
//	LOG_ROTATE_LOCAL_TIME  (bool)
//	LOG_ROTATE_COMPRESS    (bool)
//
// Для получения bool значений используется функция StringToBool(),
// допускаются определенны "вольности", кроме традиционных true/false.
//
// Если соответствующая переменная окружения не найдена
// (или имеет значение в виде пустой строки), то соответствующее
// поле структуры Conf не модифицируется.
//
// Типовое использование:
//
//	conf := xlog.Conf{}   // подготовить структуру конфигурации логгера
//	xlog.Env(&conf)       // обогатить структуру конфигурации переменными окружения
func Env(conf *Conf, prefixOpt ...string) {
	prefix := DefaultEnvPrefix
	if len(prefixOpt) != 0 {
		prefix = prefixOpt[0]
	}
	if v := os.Getenv(prefix + "LEVEL"); v != "" {
		conf.Level = v
	}
	if v := os.Getenv(prefix + "PIPE"); v != "" {
		conf.Pipe = v
	}
	if v := os.Getenv(prefix + "FILE"); v != "" {
		conf.File = v
	}
	if v := os.Getenv(prefix + "FILE_MODE"); v != "" {
		conf.FileMode = v
	}
	if v := os.Getenv(prefix + "FORMAT"); v != "" {
		conf.Format = v
	}
	if v := os.Getenv(prefix + "GOID"); v != "" {
		conf.GoId = StringToBool(v)
	}
	if v := os.Getenv(prefix + "ID"); v != "" {
		conf.IdOn = StringToBool(v)
	}
	if v := os.Getenv(prefix + "SUM"); v != "" {
		conf.SumOn = StringToBool(v)
	}
	if v := os.Getenv(prefix + "SUM_FULL"); v != "" {
		conf.SumFull = StringToBool(v)
	}
	if v := os.Getenv(prefix + "SUM_CHAIN"); v != "" {
		conf.SumChain = StringToBool(v)
	}
	if v := os.Getenv(prefix + "SUM_ALONE"); v != "" {
		conf.SumAlone = StringToBool(v)
	}
	if v := os.Getenv(prefix + "TIME"); v != "" {
		conf.TimeOff = !StringToBool(v)
		if conf.TimeOff {
			conf.TimeFormat = ""
		}
	}
	if v := os.Getenv(prefix + "TIME_LOCAL"); v != "" {
		conf.TimeLocal = StringToBool(v)
	}
	if v := os.Getenv(prefix + "TIME_MICRO"); v != "" {
		conf.TimeMicro = StringToBool(v)
	}
	if v := os.Getenv(prefix + "TIME_FORMAT"); v != "" {
		conf.TimeFormat = v
	}
	if v := os.Getenv(prefix + "SRC"); v != "" {
		conf.Src = StringToBool(v)
	}
	if v := os.Getenv(prefix + "SRC_PKG"); v != "" {
		conf.SrcPkg = StringToBool(v)
	}
	if v := os.Getenv(prefix + "SRC_FUNC"); v != "" {
		conf.SrcFunc = StringToBool(v)
	}
	if v := os.Getenv(prefix + "SRC_EXT"); v != "" {
		conf.SrcExt = StringToBool(v)
	}
	if v := os.Getenv(prefix + "COLOR"); v != "" {
		conf.ColorOff = !StringToBool(v)
	}
	if v := os.Getenv(prefix + "LEVEL_OFF"); v != "" {
		conf.LevelOff = StringToBool(v)
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
