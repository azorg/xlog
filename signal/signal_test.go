// File: "signal_test.go"

package signal

import (
	"testing"

	"github.com/azorg/xlog"
)

// Создать набор опций (*xlog.Opt)
var opt = xlog.NewOpt( /*"log-"*/ )

func TestSetupLog(_ *testing.T) {
	// Заполнить структуру конфигурации логгера
	conf := xlog.Conf{
		Level:  "trace",
		Format: "tint",
	}

	// Обогатить структуру конфигурации переменными окружения
	xlog.Env(&conf, "LOG_")

	// Обогадить conf опциями командной строки
	opt.UpdateConf(&conf)

	// Настроить все глобальные логгеры однотипно
	xlog.Setup(conf)
}

func TestSignal(t *testing.T) {
	Wait()
}

// EOF: "signal_test.go"
