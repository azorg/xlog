// File: "scan.go"

package main

import (
  "os"
  "encoding/json"
  "time"
  "fmt"
  
	"github.com/azorg/xlog"
)

// Сканировать файл журнала с целью проверки контрольных сумм
//
//  logConf - конфигурация логгера
//  fileName - имя файла сканируемого журнала (по умолчанию stdin)
//  sumChain - признак обработки цепочек
func scan(logConf xlog.Conf, fileName string, sumChain bool) {
  // Сделать вывод журнала "человеческим" (ничего лишнего)
  logConf.IdOn = false
  logConf.SumOn = false
  logConf.GoId = false
	xlog.Setup(logConf)

  xlog.Info("start scan", "app", APP_NAME, "version", Version,
    xlog.String("file", fileName), "chain", sumChain)

  file := os.Stdin
  if fileName != "" {
    var err error
    file, err = os.Open(fileName)
    if err != nil {
      xlog.Fatal("can't open log file", "err", err, "file", fileName)
      return
    }
  }

  sum := uint16(0)
  dec := json.NewDecoder(file)
  recCnt := int64(0) // счетчик записей
  errCnt := int64(0) // счетчик ошибок

  for dec.More() {
    // Распарсить JSON запись журнала
    rec := map[string]any{}
    err := dec.Decode(&rec)
    if err != nil {
      xlog.Crit("can't decode JSON from log file", "err",
        err, "file", fileName)
      return
    }

    // Распарсить запись и вычислить контрольную сумму
    res, err := xlog.ChecksumVerify(logConf.SumFull, rec)
    recCnt++

		logId := ""
		if !res.LogId.IsNil() {
			logId = res.LogId.String()
		}

		resTime := ""
    if !res.Time.IsZero() {
			resTime = res.Time.Format(time.RFC3339Nano)
		}

    log := xlog.WithGroup("res").With(
      "cnt", recCnt,
      xlog.String("time", resTime),
      "level", xlog.LevelToLabel(res.Level),
      "msg", res.Message,
      xlog.Int(xlog.GoKey, res.Goroutine),
      xlog.String("logId", logId))

    if len(res.Source) != 0 {
      log = log.With("source", res.SourceToString())
    }

    if err != nil {
      errCnt++
      log.Error("can't' verify record", "err", err, "errCnt", errCnt)
      continue
    }

    if res.Sum ^ sum != res.LogSum {
      errCnt++
      log.Error("bad log check sum",
        "logSum", fmt.Sprintf("%04x", res.LogSum),
        "sum", fmt.Sprintf("%04x", res.Sum),
        "errCnt", errCnt)
      
      if sumChain { // пропробовать выполнить коррекцию
        sum = res.LogSum
      }
    } else {
      log.Trace("scan record",
        "logSum", fmt.Sprintf("%04x", res.LogSum),
        "sum", fmt.Sprintf("%04x", res.Sum ^ sum))
    
      if sumChain {
        sum = res.LogSum
      }
    }
  } // for

  xlog.Info("finish scan", "recCnt", recCnt, "errCnt", errCnt)
}

// EOF: "scan.go"
