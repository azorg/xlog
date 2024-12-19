// File: "misc.go"

package xlog

// Remove ".go" extensin from source file name
func RemoveGoExt(file string) string {
	if n := len(file); n > 3 && file[n-3:] == ".go" {
		return file[:n-3]
	}
	return file
}

// EOF: "misc.go"
