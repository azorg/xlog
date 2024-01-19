// File: "color.go"

package xlog

import (
	"fmt"
	"log/slog" // go>=1.21
	//"golang.org/x/exp/slog" // depricated for go>=1.21
)

// ANSI modes
const (
	AnsiReset            = "\033[0m"        // All attributes off
	AnsiFaint            = "\033[2m"        // Decreased intensity
	AnsiResetFaint       = "\033[22m"       // Normal color (reset faint)
	AnsiRed              = "\033[31m"       // Red
	AnsiGreen            = "\033[32m"       // Green
	AnsiYellow           = "\033[33m"       // Yellow
	AnsiBlue             = "\033[34m"       // blue
	AnsiMagenta          = "\033[35m"       // Magenta
	AnsiCyan             = "\033[36m"       // Cyan
	AnsiWhile            = "\033[37m"       // While
	AnsiBrightRed        = "\033[91m"       // Bright Red
	AnsiBrightRedNoFaint = "\033[91;22m"    // Bright Red and normal intensity
	AnsiBrightGreen      = "\033[92m"       // Bright Green
	AnsiBrightYellow     = "\033[93m"       // Bright Yellow
	AnsiBrightBlue       = "\033[94m"       // Bright Blue
	AnsiBrightMagenta    = "\033[95m"       // Bright Magenta
	AnsiBrightCyan       = "\033[96m"       // Bright Cyan
	AnsiBrightWight      = "\033[97m"       // Bright White
	AnsiBrightRedFaint   = "\033[91;2m"     // Bright Red and decreased intensity
	AnsiBlackOnWhite     = "\033[30;107;1m" // Black on Bright White background
	AnsiBlueOnWhite      = "\033[34;47;1m"  // Blue on Bright White background
	AnsiWhiteOnMagenta   = "\033[37;45;1m"  // Bright White on Magenta background
	AnsiWhiteOnRed       = "\033[37;41;1m"  // White on Red background
)

// Level keys ANSI colors
const (
	AnsiFlood    = AnsiGreen
	AnsiTrace    = AnsiBrightBlue
	AnsiDebug    = AnsiBrightCyan
	AnsiInfo     = AnsiBrightGreen
	AnsiNotice   = AnsiBrightMagenta
	AnsiWarn     = AnsiBrightYellow
	AnsiError    = AnsiBrightRed
	AnsiCritical = AnsiWhiteOnRed
	AnsiFatal    = AnsiWhiteOnMagenta
	AnsiPanic    = AnsiBlackOnWhite
)

// Log part colors
const (
	AnsiTime   = AnsiYellow
	AnsiSource = AnsiMagenta
	AnsiKey    = AnsiCyan
	AnsiErrKey = AnsiRed
	AnsiErrVal = AnsiBrightRed
)

// ColorString() returns a color label for the level
func (lp *Level) ColorString() string {
	str := func(ansi, base string, delta slog.Level) string {
		if delta == 0 {
			return ansi + base + AnsiReset
		}
		return fmt.Sprintf("%s%s%+d"+AnsiReset, ansi, base, delta)
	}

	l := slog.Level(*lp)
	switch {
	case l < LevelTrace:
		return str(AnsiFlood, LabelFlood, l-LevelFlood)
	case l < LevelDebug:
		return str(AnsiTrace, LabelTrace, l-LevelTrace)
	case l < LevelInfo:
		return str(AnsiDebug, LabelDebug, l-LevelDebug)
	case l < LevelNotice:
		return str(AnsiInfo, LabelInfo, l-LevelInfo)
	case l < LevelWarn:
		return str(AnsiNotice, LabelNotice, l-LevelNotice)
	case l < LevelError:
		return str(AnsiWarn, LabelWarn, l-LevelWarn)
	case l < LevelCritical:
		return str(AnsiError, LabelError, l-LevelError)
	case l < LevelFatal:
		return str(AnsiCritical, LabelCritical, l-LevelCritical)
	case l < LevelPanic:
		return str(AnsiFatal, LabelFatal, l-LevelFatal)
	case l < LevelSilent:
		return str(AnsiPanic, LabelPanic, l-LevelPanic)
	default: // l >= LevelSilent
		return str(AnsiPanic, LabelSilent, l-LevelSilent)
	}
}

// EOF: "color.go"
