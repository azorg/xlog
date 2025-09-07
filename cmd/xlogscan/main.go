// File: "main.go"

package main

import (
  "flag"
  "os"
	
	"github.com/azorg/xlog"
	"github.com/azorg/xlog/signal"
)

// Опции командной строки
type Opt struct {
	File  string   // входной файл журнала
  Chain bool     // признак обработки цепочки
}

func main() {
  // Разобрать опции командной строки
	for _, o := range os.Args[1:] {
		if o == "-help" || o == "--help" || o == "help" {
			printUsage()
		} else if o == "-v" || o == "--version" || o == "version" {
			printVersion()
		} else if len(o) > 0 && o[0:1] != "-" {
			break // abort by first command
		}
	}
	
  // Настройки логгера по умолчанию
	logConf := xlog.Conf{
    Format:     "human",
    ColorOff:   false,
    Level:      "flood",
    Src:        true,
    //TimeFormat: "timeOnlyMilli",
  }
	
  // Получить настройки логгера из переменных окружения
	xlog.Env(&logConf)

  opt := &Opt{}
	flag.StringVar(&opt.File, "file", "", "Input log file (use stdin by default)")
	flag.BoolVar(&opt.Chain, "chain", false, "Check chain")
  
  logOpt := xlog.NewOpt()
  flag.Parse()

	// Добавить настройки логгера, заданные в командной строке
	logOpt.UpdateConf(&logConf)

	// Настроить логгер xlog по умолчанию
	xlog.Setup(logConf)

	if DEBUG_CTRL_BS {
		// Подключить обработчик сигнала на Ctrl+\ для безусловного останова приложения
		go func() {
			<-signal.CtrlBS
			xlog.Fatal(APP_NAME + " aborted")
		}()
	}
		
  go func() {
		<-signal.CtrlC
		xlog.Fatal(APP_NAME + " aborted")
	}()

	if LOGROTATE_SIGHUP {
    // Подключить обработчик SIGHUP
		go func() {
			for {
				_, ok := <-signal.SIGHUP
				if !ok {
					return
				}
				if xlog.IsRotatable() {
					xlog.Debug("rotate log by SIGHUP")
					err := xlog.Rotate()
					if err != nil {
						xlog.Crit("can't rotate log", "err", err)
					}
				}
			} // for
		}()
	}
	
  // Разобрать командную строку
	args := flag.Args() // аргументы командной строки без обработанных флагов
	argc := len(args)

	if argc == 0 {
    scan(logConf, opt.File, opt.Chain)
		return
	}
	
  cmd := args[0] // argc != 0
	switch cmd {
	case "scan":
    scan(logConf, opt.File, opt.Chain)
  case "test":
    test(logConf)
  default:
		xlog.Fatal("Unknown command (run with --help option)", "cmd", cmd)
	} // switch
}

// EOF: "main.go"
