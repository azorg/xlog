// File: "misc.go"

package xlog

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

// Crop function name
func CropFuncName(function string) string {
	_, f := filepath.Split(function)
	parts := strings.Split(f, ".")
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return strings.Join(parts[1:], ".") // FIXME
}

// Return function name
func GetFuncName(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	function := runtime.FuncForPC(pc).Name()
	return CropFuncName(function)
}

// Convert file mode string (oct like "0640") to fs.FileMode
func FileMode(mode string) fs.FileMode {
	if mode == "" {
		mode = FILE_MODE
	}
	perm, err := strconv.ParseInt(mode, 8, 10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: bad logfile mode='%s'; set mode=0%03o\n",
			mode, DEFAULT_FILE_MODE)
		return DEFAULT_FILE_MODE
	}
	return fs.FileMode(perm & 0777)
}

// EOF: "misc.go"
