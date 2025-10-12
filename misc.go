// File: "misc.go"

package xlog

import (
	"fmt"
	"io/fs"
	"log/slog" // go>=1.21
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	// FIXME: "golang.org/x/exp/slog" // экспериментальный пакет для go=1.20 только
)

// removeGoExt обрезает расширение ".go" из имени файла
func removeGoExt(file string) string {
	if n := len(file); n > 3 && file[n-3:] == ".go" {
		return file[:n-3]
	}
	return file
}

// cropFuncName укорачивает специальным образом имя функции
func cropFuncName(function string) string {
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

// getSource формирует slog.Source на основе поля PC в slog.Record
func getSource(pc uintptr) *slog.Source {
	fs := runtime.CallersFrames([]uintptr{pc})
	f, _ := fs.Next()
	if f.Func == nil {
		return nil
	}
	return &slog.Source{
		Function: f.Function,
		File:     f.File,
		Line:     f.Line,
	}
}

// getFuncName возвращает имя текущей функции
//func getFuncName(skip int) string {
//	pc, _, _, _ := runtime.Caller(skip)
//	function := runtime.FuncForPC(pc).Name()
//	return cropFuncName(function)
//}

// fileMode преобразует права доступа к файлу в восьмеричной Unix нотации
// (например "0644") к типу fs.FileMode
func fileMode(mode string) fs.FileMode {
	if mode == "" {
		return FileModeDefault
	}
	perm, err := strconv.ParseInt(mode, 8, 10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: bad logfile mode='%s'; set mode=0%03o\n",
			mode, FileModeOnError)
		return FileModeOnError
	}
	return fs.FileMode(perm & 0777)
}

// EOF: "misc.go"
