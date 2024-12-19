// File: "const.go"

package xlog

// Default logger configure
const (
	FILE      = "stdout" // log file path OR stdout/stderr ("" -> stdout)
	FILE_MODE = "0640"   // log file mode (if FILE is not stdout/stderr)
	LEVEL     = LvlInfo  // log level (flood/trace/debug/info/warn/error/critical/fatal/silent)
	SLOG      = false    // use slog instead standard log (slog.TextHandler)
	JSON      = false    // use JSON log (slog.JSONHandelr)
	TINT      = false    // use tinted (colorized) log (xlog.TintHandler)
	TIME      = false    // add time stamp
	TIME_US   = false    // us time stamp (only if SLOG=false)
	TIME_TINT = ""       // tinted log time format (~time.Kitchen, "15:04:05.999")
	SRC       = false    // log file name and line number
	SRC_LONG  = false    // log long file path (directory + file name)
	NO_EXT    = false    // remove ".go" extension from file name
	NO_LEVEL  = false    // don't print log level tag to log (~level="INFO")
	NO_COLOR  = false    // don't use tinted colors (only if Tint=true)
	PREFIX    = ""       // add prefix to standard log (SLOG=false)
	ADD_KEY   = ""       // add key to structured log (SLOG=true)
	ADD_VALUE = ""       // add value to structured log (SLOG=true
)

const (
	// Add addition log level marks (TRACE/NOTICE/FATAL/PANIC)
	ADD_LEVELS = true

	// Log file mode in error configuration
	DEFAULT_FILE_MODE = 0600 // read/write only for owner for more secure

	// Set false for go > 1.21 with log/slog
	OLD_SLOG_FIX = false // runtime.Version() < go1.21.0

	// Pretty alignment time format in tinted handler (add zeros to end)
	TINT_ALIGN_TIME = true
)

// EOF: "const.go"
