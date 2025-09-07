// File: "signal_windows.go"
//go:build windows
// +build windows

package signal

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/azorg/xlog"
)

// Setup Ctrl+C | Ctrl+Z | Ctrl+\ | SIGTERM | SIGHUP channels
func init() {
	CtrlC = make(chan None, CHAN_SIZE)
	CtrlZ = make(chan None, CHAN_SIZE)
	CtrlBS = make(chan None, CHAN_SIZE)
	SIGTERM = make(chan None, CHAN_SIZE)
	SIGHUP = make(chan None, CHAN_SIZE)

	ch := make(chan os.Signal, 1)

	// Can't use syscall under Windows :(
	signal.Notify(ch, os.Interrupt)

	go func() {
		for sig := range ch {
			fmt.Fprint(os.Stderr, "\r\n")
			xlog.Trace("interrupt signal", "sig", sig)

			// Only Ctrl+C emulated (SIGINT)
			send(CtrlC)
		} // for
	}()
}

// EOF: "signal_windows.go"
