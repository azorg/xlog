// File: "signal_wait.go"

package signal

import (
	"fmt"

	"github.com/azorg/xlog"
)

// Debug wait
func Wait() bool {
	fmt.Println(`press Ctrl+C to resume or Ctrl+\ to abort`)
	select {
	case <-CtrlC:
		xlog.Info(`resume application by Ctrl+C (SIGINT)`)
		return true

	case <-CtrlBS:
		xlog.Fatal(`abort application by Ctrl+\ (SIGQUIT)`)

	case <-SIGHUP:
		xlog.Info(`SIGHUP received`)

	case <-SIGTERM:
		xlog.Info(`SIGTERM received`)
	}
	return false
}

// EOF: "signal.go"
