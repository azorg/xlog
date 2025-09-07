// File: "format.go"

package xlog

import "strings"

// Форматы журнала (4 типа)
const (
	// Формат журнала JSON, эксплуатируется slog.JSONHandler
	LogFormatJSON = "json"
	LogFormatProd = "prod"

	// Формат журнала Logfmt, эксплуатируется slog.TextHandler
	LogFormatText   = "text"
	LogFormatSlog   = "slog"
	LogFormatLogfmt = "logfmt"

	// Формат журнала HumanText, эксплуатируется кастомный TintHandler
	LogFormatTint   = "tint"
	LogFormatTinted = "tinted"
	LogFormatHuman  = "human"

	// Формат журнала, предоставляемый библиотекой Go по умолчанию
	LogFormatStd     = "std"
	LogFormatDefault = "default"
)

// Формат журнала для внутреннего использования
type logFmt byte

const (
	logFmtDefault logFmt = iota // default/std slog.Logger
	logFmtJSON                  // JSON   (slog.JSONHandler)
	logFmtText                  // Logfmt (slog.TextHandler)
	logFmtTint                  // Tinted (TintHandler)

	// Формат журнала используемый по умолчанию (если передана пустая строка)
	defaultLogFmt = logFmtTint
)

// Преобразовать строку ("json", "text", "tint", "std") к типу logFmt
func logFormat(format string) logFmt {
	switch strings.ToLower(format) {
	case LogFormatDefault, LogFormatStd:
		return logFmtDefault

	case LogFormatJSON, LogFormatProd:
		return logFmtJSON

	case LogFormatText, LogFormatLogfmt, LogFormatSlog:
		return logFmtText

	case LogFormatTinted, LogFormatTint, LogFormatHuman:
		return logFmtTint

	default:
		return defaultLogFmt
	}
}

// EOF: "format.go"
