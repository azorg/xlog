// File: "shortcuts.go"
// Здесь будут храниться шаблоны наиболее часто используемых конфигураций.

package xlog

// OriginalConf - конфигурация максимально приближенная к стандартному
// slog.TextHandler "из коробки" и совместимая с некоторыми проектами
func OriginalConf() Conf {
	return Conf{
		Level:     "info",
		Format:    "logfmt",
		IdOn:      false,
		SumOn:     false,
		Src:       true,
		SrcPkg:    true,
		SrcFunc:   true,
		TimeLocal: true,
	}
}

// SetupOriginalConf устанавливает конфигурацию логгера максимально
// близкую к стандартному slog.TextHandler "из коробки" и совместимую
// с некоторыми проектами
func SetupOriginalConf() Conf {
	conf := OriginalConf()
	Env(&conf)
	Setup(conf)
	return conf
}

// EOF: "shortcuts.go"
