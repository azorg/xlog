// File: "rotator.go"

package xlog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Interface of log rotator
type Rotator interface {
	io.Writer
	Rotable() bool // check rotation possible
	Rotate() error // rotate log
	Close() error  // close file
}

// Custop io.Writer
type writer struct{ io.Writer }

// Stdout/Stderr writer
type pipe struct{ *os.File }

// Direct *os.File writer
type file struct{ *os.File }

// Rotator by *lumberjack.Logger
type rotator struct{ *lumberjack.Logger }

// Ensure file/rotator implements Rotator
var _ Rotator = writer{}
var _ Rotator = pipe{}
var _ Rotator = file{}
var _ Rotator = rotator{}

// Convert io.Writer to Rotator
func newWriter(w io.Writer) Rotator { return writer{w} }

// Convert os.Stdout/os.Stderr to Rotator
func newPipe(f *os.File) Rotator { return pipe{f} }

// Convert *os.File to Rotator
func newFile(f *os.File) Rotator { return file{f} }

// Select (open) log file, setup log rotation (return file or rotator)
func newRotator(fileName, mode string, rotate *RotateOpt) Rotator {
	switch fileName {
	case "stdout", "os.Stdout", "":
		return pipe{os.Stdout}
	case "stderr", "os.Stderr":
		return pipe{os.Stderr}
	}

	// Get oct file mode from string
	perm := FileMode(mode)

	// Make logs directory with permission
	dir := filepath.Dir(fileName)
	if dir != "" {
		dirPerm := perm | ((perm & 0044) >> 2) | 0700 // FIXME: some magic
		err := os.MkdirAll(dir, dirPerm)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"ERROR: can't create logfile directory: %v; use stdout\n", err)
			return pipe{os.Stdout}
		}
	}

	// Open log file
	out, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, perm)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"ERROR: can't open/create logfile: %v; use stdout\n", err)
		return pipe{os.Stdout}
	}

	if rotate == nil || !rotate.Enable {
		return file{out}
	}

	out.Close()
	return rotator{&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    rotate.MaxSize,
		MaxAge:     rotate.MaxAge,
		MaxBackups: rotate.MaxBackups,
		LocalTime:  rotate.LocalTime,
		Compress:   rotate.Compress,
	}}
}

// Check rotation possible
func (p writer) Rotable() bool  { return false }
func (p pipe) Rotable() bool    { return false }
func (p file) Rotable() bool    { return false }
func (p rotator) Rotable() bool { return true }

// Rotate log
func (w writer) Rotate() error { return nil } // do nothing
func (p pipe) Rotate() error   { return nil } // do nothing
func (f file) Rotate() error   { return nil } // do nothing

// Close log
func (w writer) Close() error { return nil } // do nothing
func (p pipe) Close() error   { return nil } // do nothing

// EOF: "rotator.go"
