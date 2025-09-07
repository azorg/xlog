// File: "const.go"

package main

const (
	APP_NAME         = "xlogscan"         // имя приложения
	APP_DESCRIPTION  = "X-Logger Scanner" // описание приложения
	VERSION_MAJOR    = "x"                // старший номер версии
	VERSION_MINOR    = "y"                // младший номер версии
	VERSION_BUILD    = "z"                // номер сборки
	DEBUG_CTRL_BS    = true               // Подключить обработчик сигнала на Ctrl+\
	LOGROTATE_SIGHUP = true               // Ротация журнала по SIGHUP
)

var Version = VERSION_MAJOR + "." + VERSION_MINOR + "." + VERSION_BUILD
var GitHash string
var BuildTime string

func init() {
	if GitHash != "" {
		Version += "-" + GitHash
	}
}

// EOF: "const.go"
