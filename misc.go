// File: "misc.go"

package xlog

import (
	"io/fs"
	"runtime"
	"strconv"
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

// Convert file mode string (oct like "0640") to fs.FileMode
func FileMode(mode string) fs.FileMode {
	if mode == "" {
		mode = FILE_MODE
	}
	perm, err := strconv.ParseInt(mode, 8, 10)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "ERROR: bad logfile mode='%s'; set mode=0%03o\n",
		//	mode, DEFAULT_FILE_MODE)
		return DEFAULT_FILE_MODE
	}
	return fs.FileMode(perm & 0777)
}

// EOF: "misc.go"
