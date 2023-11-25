// File: "color.go"

package xlog

import "fmt"

// ANSI modes
const (
	AnsiReset            = "\033[0m"       // All attributes off
	AnsiFaint            = "\033[2m"       // Decreased intensity
	AnsiResetFaint       = "\033[22m"      // Normal color (reset faint)
	AnsiRed              = "\033[31m"      // Red
	AnsiGreen            = "\033[32m"      // Green
	AnsiYellow           = "\033[33m"      // Yellow
	AnsiBlue             = "\033[34m"      // blue
	AnsiMagenta          = "\033[35m"      // Magenta
	AnsiCyan             = "\033[36m"      // Cyan
	AnsiWhile            = "\033[37m"      // While
	AnsiBrightRed        = "\033[91m"      // Bright Red
	AnsiBrightRedNoFaint = "\033[91;22m"   // Bright Red and normal intensity
	AnsiBrightGreen      = "\033[92m"      // Bright Green
	AnsiBrightYellow     = "\033[93m"      // Bright Yellow
	AnsiBrightBlue       = "\033[94m"      // Bright Blue
	AnsiBrightMagenta    = "\033[95m"      // Bright Magenta
	AnsiBrightCyan       = "\033[96m"      // Bright Cyan
	AnsiBrightWight      = "\033[97m"      // Bright White
	AnsiBrightRedFaint   = "\033[91;2m"    // Bright Red and decreased intensity
	AnsiWhiteOnMagenta   = "\033[37;45;1m" // Bright White on Magenta background
	AnsiWhiteOnRed       = "\033[37;41;1m" // Bright White on Red backgroun
)

// Level keys ANSI colors
const (
	AnsiTrace  = AnsiBrightBlue
	AnsiDebug  = AnsiBrightCyan
	AnsiInfo   = AnsiBrightGreen
	AnsiNotice = AnsiBrightMagenta
	AnsiWarn   = AnsiBrightYellow
	AnsiError  = AnsiBrightRed
	AnsiFatal  = AnsiWhiteOnMagenta
	AnsiPanic  = AnsiWhiteOnRed
)

// Log part colors
const (
	AnsiTime   = AnsiFaint
	AnsiSource = AnsiFaint
	AnsiKey    = AnsiFaint
	AnsiErrKey = AnsiBrightRedFaint
	AnsiErrVal = AnsiBrightRedNoFaint
)

// ColorString() returns a color label for the level
func (l Level) ColorString() string {
	str := func(ansi, base string, delta Level) string {
		if delta == 0 {
			return ansi + base + AnsiReset
		}
		return fmt.Sprintf("%s%s%+d"+AnsiReset, ansi, base, delta)
	}

	switch {
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
	case l < LevelFatal:
		return str(AnsiError, LabelError, l-LevelError)
	case l < LevelPanic:
		return str(AnsiFatal, LabelFatal, l-LevelFatal)
	case l < LevelSilent:
		return str(AnsiPanic, LabelPanic, l-LevelPanic)
	default: // l >= LevelSilent
		return str(AnsiPanic, LabelSilent, l-LevelSilent)
	}
}

// EOF: "color.go"
