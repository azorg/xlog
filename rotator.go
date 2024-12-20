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
	Writer() io.Writer // get io.Writer for write
	Rotate() error     // rotate log
	Close() error      // close file
}

// Stdout/Stderr writer
type pipe struct{ *os.File }

// Direct *os.File writer
type file struct{ *os.File }

// Rotator by *lumberjack.Logger
type rotator struct{ *lumberjack.Logger }

// Ensure file/rotator implements Rotator
var _ Rotator = pipe{}
var _ Rotator = file{}
var _ Rotator = rotator{}

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

	if rotate != nil && rotate.Enable {
		return rotator{&lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    rotate.MaxSize,
			MaxAge:     rotate.MaxAge,
			MaxBackups: rotate.MaxBackups,
			LocalTime:  rotate.LocalTime,
			Compress:   rotate.Compress,
		}}
	}

	perm := FileMode(mode)

	// Make directory
	dir := filepath.Dir(fileName)
	if dir != "" {
		err := os.MkdirAll(dir, perm|0711)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"ERROR: can't create logfile directory: %v; use os.Stdout\n", err)
			return pipe{os.Stdout}
		}
	}

	// Open log file
	out, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"ERROR: can't create logfile: %v; use os.Stdout\n", err)
		return pipe{os.Stdout}
	}

	return file{out}
}

// Get io.Writer
func (p pipe) Writer() io.Writer    { return p.File }
func (f file) Writer() io.Writer    { return f.File }
func (r rotator) Writer() io.Writer { return r.Logger }

// Rotate log
func (p pipe) Rotate() error    { return nil } // do nothing
func (f file) Rotate() error    { return nil } // do nothing
func (r rotator) Rotate() error { return r.Logger.Rotate() }

// Close log
func (p pipe) Close() error    { return nil } // do nothing
func (f file) Close() error    { return f.File.Close() }
func (r rotator) Close() error { return r.Logger.Close() }

// EOF: "rotator.go"
