// Package signal implement some signal bridge to channels.
// File: "signal.go"
package signal

const CHAN_SIZE = 10

type None struct{}

var (
	CtrlC   chan None // Ctrl+C -> SIGINT or SIGTERM
	CtrlZ   chan None // Ctrl+Z -> SIGTSTP
	CtrlBS  chan None // Ctrl+\ -> SIGQUIT
	SIGTERM chan None
	SIGHUP  chan None
)

func send(ch chan<- None) bool {
	select {
	case ch <- None{}:
		return true
	default:
		return false
	}
}

// EOF: "signal.go"
