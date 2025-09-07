// File: "ansi.go"

package xlog

// ANSI Escape последовательности для TintHandler'а
const (
	ansiReset            = "\033[0m"        // All attributes off
	ansiFaint            = "\033[2m"        // Decreased intensity
	ansiResetFaint       = "\033[22m"       // Normal color (reset faint)
	ansiRed              = "\033[31m"       // Red
	ansiGreen            = "\033[32m"       // Green
	ansiYellow           = "\033[33m"       // Yellow
	ansiBlue             = "\033[34m"       // blue
	ansiMagenta          = "\033[35m"       // Magenta
	ansiCyan             = "\033[36m"       // Cyan
	ansiWhile            = "\033[37m"       // While
	ansiBrightRed        = "\033[91m"       // Bright Red
	ansiBrightRedNoFaint = "\033[91;22m"    // Bright Red and normal intensity
	ansiBrightGreen      = "\033[92m"       // Bright Green
	ansiBrightYellow     = "\033[93m"       // Bright Yellow
	ansiBrightBlue       = "\033[94m"       // Bright Blue
	ansiBrightMagenta    = "\033[95m"       // Bright Magenta
	ansiBrightCyan       = "\033[96m"       // Bright Cyan
	ansiBrightWight      = "\033[97m"       // Bright White
	ansiBrightRedFaint   = "\033[91;2m"     // Bright Red and decreased intensity
	ansiBlackOnWhite     = "\033[30;107;1m" // Black on Bright White background
	ansiBlueOnWhite      = "\033[34;47;1m"  // Blue on Bright White background
	ansiWhiteOnMagenta   = "\033[37;45;1m"  // Bright White on Magenta background
	ansiWhiteOnRed       = "\033[37;41;1m"  // White on Red background

	ansiEsc = '\u001b' // символ Escape
)

// Подсветка меток уровней логирования
const (
	ansiFlood  = ansiGreen
	ansiTrace  = ansiBrightBlue
	ansiDebug  = ansiBrightCyan
	ansiInfo   = ansiBrightGreen
	ansiNotice = ansiBrightMagenta
	ansiWarn   = ansiBrightYellow
	ansiError  = ansiBrightRed
	ansiCrit   = ansiWhiteOnRed
	ansiAlert  = ansiWhiteOnRed
	ansiEmerg  = ansiWhiteOnRed
	ansiFatal  = ansiWhiteOnMagenta
	ansiPanic  = ansiBlackOnWhite
)

// Подсветка стандартных полей в журнале
const (
	ansiTime   = ansiYellow    // метка времени
	ansiSource = ansiMagenta   // ссылка на исходные тексты
	ansiKey    = ansiCyan      // ключ атрибута
	ansiErrKey = ansiRed       // ключ ошибки (err)
	ansiErrVal = ansiBrightRed // текст ошибки
)

// EOF: "ansi.go"
