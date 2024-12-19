// File: "misc.go"

package xlog

import (
	"runtime"
	"strings"
)

// Remove ".go" extension from source file name
func RemoveGoExt(file string) string {
	if n := len(file); n > 3 && file[n-3:] == ".go" {
		return file[:n-3]
	}
	return file
}

// Return function name
func GetFuncName(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

// EOF: "misc.go"
