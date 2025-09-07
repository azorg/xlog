// File: "signal_unix.go"
//go:build linux || aix
// +build linux aix

package signal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/azorg/xlog"
)

// Setup Ctrl+C | Ctrl+Z | Ctrl+\ | SIGTERM | SIGHUP channels
func init() {
	CtrlC = make(chan None, CHAN_SIZE)
	CtrlZ = make(chan None, CHAN_SIZE)
	CtrlBS = make(chan None, CHAN_SIZE)
	SIGTERM = make(chan None, CHAN_SIZE)
	SIGHUP = make(chan None, CHAN_SIZE)

	ch := make(chan os.Signal, CHAN_SIZE)

	sigList := []os.Signal{
		syscall.SIGINT,  // Ctrl-C
		syscall.SIGTSTP, // Ctrl-Z
		syscall.SIGQUIT, // Ctrl-\
		//os.Interrupt, // syscall.SIGINT (Ctrl-C)
		syscall.SIGTERM,
		syscall.SIGHUP,
	}

	//signal.Ignore(sigList...)
	signal.Notify(ch, sigList...)

	go func() {
		for sig := range ch {
			switch sig {
			case syscall.SIGINT:
				fmt.Fprint(os.Stderr, "\r\n")
				xlog.Trace("SIGINT received (or Ctrl+C pressed)")
				send(CtrlC)

			case syscall.SIGTSTP:
				fmt.Fprint(os.Stderr, "\r\n")
				xlog.Trace("SIGTSTP received (or Ctrl+Z pressed)")
				send(CtrlZ)

			case syscall.SIGQUIT:
				fmt.Fprint(os.Stderr, "\r\n")
				xlog.Trace("SIGQUIT received (or Ctrl+\\ pressed)")
				send(CtrlBS)

			case syscall.SIGTERM:
				xlog.Trace("SIGTERM received")
				send(SIGTERM)

			case syscall.SIGHUP:
				xlog.Trace("SIGHUP received")
				send(SIGHUP)

			default:
				xlog.Warn("unknown signal received", "sig", sig)
			} // switch
		} // for
	}()
}

// EOF: "signal_unix.go"
