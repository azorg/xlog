package signal // import "clog/signal"

Package signal implement some signal bridge to channels. File: "signal.go"

CONSTANTS

const CHAN_SIZE = 10

VARIABLES

var (
	CtrlC   chan None // Ctrl+C -> SIGINT or SIGTERM
	CtrlZ   chan None // Ctrl+Z -> SIGTSTP
	CtrlBS  chan None // Ctrl+\ -> SIGQUIT
	SIGTERM chan None
	SIGHUP  chan None
)

FUNCTIONS

func Wait() bool
    Debug wait


TYPES

type None struct{}

