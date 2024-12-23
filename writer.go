// File: "writer.go"

package xlog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Interface of log writer
type Writer interface {
	Rotable() bool // check rotation possible
	Rotate() error // rotate log
	Close() error  // close file
	io.Writer
}

// Null log writer
type null struct{}

// stdout/stderr pipe writer only
type pipe struct{ *os.File }

// Custom io.Writer only
type writer struct{ io.Writer }

// Custom io.Writer + pipe writer
type pipeAndWriter struct {
	pipe      *os.File // os.Stdout/os.Stderr
	io.Writer          // custom writer
}

// File writer only
type file struct{ *os.File }

// File writer + pipe writer
type pipeAndFile struct {
	pipe     *os.File // os.Stdout/os.Stderr
	*os.File          // regular file
}

// Log rotator by *lumberjack.Logger
type rotator struct{ *lumberjack.Logger }

// Log rotator by *lumberjack.Logger + pipe writer
type pipeAndRotator struct {
	pipe               *os.File // os.Stdout/os.Stderr
	*lumberjack.Logger          // log rotator
}

// Ensure we fully implements writers
var _ Writer = null{}
var _ Writer = pipe{}
var _ Writer = writer{}
var _ Writer = pipeAndWriter{}
var _ Writer = file{}
var _ Writer = pipeAndFile{}
var _ Writer = rotator{}
var _ Writer = pipeAndRotator{}

// Select pipe (os.Stdout, os.Stderr or nil)
func getPipe(pipeName string, noFile bool) *os.File {
	if pipeName == "" && noFile {
		return os.Stdout // use stdout on zero configuration
	}
	switch strings.ToLower(pipeName) {
	case "stdout", "os.stdout", "standard":
		return os.Stdout
	case "stderr", "os.stderr", "error":
		return os.Stderr
	default: // ~ "", "null", "nil", "none", "off"
		return nil
	}
}

// Create null{} or pipe{} writer
func newPipe(p *os.File) Writer {
	if p == nil {
		return null{}
	}
	return pipe{p}
}

// Create null{}, pipe{}, file{}, pipeAndFile{}, rortator{} or pipeAndRotator{}
// select log file, setup log rotation, return Writer interface object
func newRotator(pipeName, fileName, mode string, rotate *RotateOpt) Writer {
	noFile := fileName == ""
	p := getPipe(pipeName, noFile) // os.Stdout, os.Stderr or nil
	if noFile {
		return newPipe(p)
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
			return newPipe(p)
		}
	}

	// Open log file
	out, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, perm)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"ERROR: can't open/create logfile: %v; use stdout\n", err)
		return newPipe(p)
	}

	if rotate == nil || !rotate.Enable {
		if p == nil {
			return file{out}
		}
		return pipeAndFile{
			pipe: p,
			File: out,
		}
	}

	out.Close()
	logger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    rotate.MaxSize,
		MaxAge:     rotate.MaxAge,
		MaxBackups: rotate.MaxBackups,
		LocalTime:  rotate.LocalTime,
		Compress:   rotate.Compress,
	}

	if p == nil {
		return rotator{logger}
	}

	return pipeAndRotator{
		pipe:   p,
		Logger: logger,
	}
}

// Create null{}, writer{} or pipeAndWriter{}
func newWriter(pipeName string, w io.Writer) Writer {
	p := getPipe(pipeName, w == nil) // os.Stdout, os.Stderr or nil
	if p == nil && w == nil {
		return null{}
	}
	if p == nil {
		return writer{w}
	}
	if w == nil {
		return pipe{p}
	}
	return pipeAndWriter{
		pipe:   p,
		Writer: w,
	}
}

// Check rotation possible
func (_ null) Rotable() bool           { return false }
func (_ pipe) Rotable() bool           { return false }
func (_ writer) Rotable() bool         { return false }
func (_ pipeAndWriter) Rotable() bool  { return false }
func (_ file) Rotable() bool           { return false }
func (_ pipeAndFile) Rotable() bool    { return false }
func (_ rotator) Rotable() bool        { return true }
func (_ pipeAndRotator) Rotable() bool { return true }

// Rotate log
func (_ null) Rotate() error          { return nil } // do nothing
func (_ pipe) Rotate() error          { return nil } // do nothing
func (_ writer) Rotate() error        { return nil } // do nothing
func (_ pipeAndWriter) Rotate() error { return nil } // do nothing
func (_ file) Rotate() error          { return nil } // do nothing
func (_ pipeAndFile) Rotate() error   { return nil } // do nothing

// Close log
func (_ null) Close() error           { return nil } // do nothing
func (_ pipe) Close() error           { return nil } // do nothing
func (_ writer) Close() error         { return nil } // do nothing
func (_ pipeAndWriter) Close() error  { return nil } // do nothing
func (f pipeAndFile) Close() error    { return f.File.Close() }
func (r pipeAndRotator) Close() error { return r.Logger.Close() }

// Write to /dev/null
func (_ null) Write(_ []byte) (int, error) { return 0, nil }

// Write to pipe and custom writer
func (w pipeAndWriter) Write(b []byte) (n int, err error) {
	w.pipe.Write(b)
	return w.Writer.Write(b)
}

// Write to pipe and file
func (f pipeAndFile) Write(b []byte) (n int, err error) {
	f.pipe.Write(b)
	return f.File.Write(b)
}

// Write to pipe and rotator
func (r pipeAndRotator) Write(b []byte) (n int, err error) {
	r.pipe.Write(b)
	return r.Logger.Write(b)
}

// EOF: "writer.go"
