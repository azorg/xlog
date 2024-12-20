// File: "rotator.go"

package xlog

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Interface of *lumberjack.Logger + get io.Writer
type Rotator interface {
	io.WriteCloser
	Rotate() error
	Writer() io.Writer
}

// Any log writer
type Writer struct {
	out any // *os.File or *lumberjack.Logger or nil
}

// Ensure Writer implements Rotator
var _ Rotator = Writer{}

// Select (open) log file, setup log rotation
// (return *os.File or *lumberjack.Logger)
func NewWriter(file, mode string, rotate *RotateOpt) Writer {
	switch file {
	case "stdout", "os.Stdout", "":
		return Writer{os.Stdout}
	case "stderr", "os.Stderr":
		return Writer{os.Stderr}
	}

	if rotate != nil && rotate.Enable {
		return Writer{&lumberjack.Logger{
			Filename:   file,
			MaxSize:    rotate.MaxSize,
			MaxAge:     rotate.MaxAge,
			MaxBackups: rotate.MaxBackups,
			LocalTime:  rotate.LocalTime,
			Compress:   rotate.Compress,
		}}
	}

	perm := fileMode(mode)
	out, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: can't create logfile: %v; use os.Stdout\n", err)
		return Writer{os.Stdout}
	}

	return Writer{out}
}

// Write data to file
func (w Writer) Write(p []byte) (n int, err error) {
	switch writer := w.out.(type) {
	case *os.File:
		return writer.Write(p)
	case *lumberjack.Logger:
		return writer.Write(p)
	}
	return 0, nil
}

// Close log file
func (w Writer) Close() error {
	switch closer := w.out.(type) {
	case *os.File:
		return closer.Close()
	case *lumberjack.Logger:
		return closer.Close()
	}
	return nil
}

// Rotate log file
func (w Writer) Rotate() error {
	rotator, ok := w.out.(*lumberjack.Logger)
	Debugf("type of out is %T", w.out)  //!!!
	Debugf("value of out is %v", w.out) //!!!
	if ok {
		Debug("rotate")
		return rotator.Rotate()
	}
	return nil
}

// Get io.Writer
func (w Writer) Writer() io.Writer {
	switch writer := w.out.(type) {
	case *os.File:
		return writer
	case *lumberjack.Logger:
		return writer
	}
	return nil
}

// EOF: "rotator.go"
