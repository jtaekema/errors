package errors

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

// Error interface.
type Error interface {
	Cause() error
	Error() string
	Trace() []string
	Details() string
}

// New creates a new raw error and trace it.
func New(text string) error {
	err := &wrappedError{
		traceMessage: text,
	}
	err.addCaller(1)
	return err
}

// Newf creates a new raw error and trace it.
func Newf(format string, args ...interface{}) error {
	err := &wrappedError{
		traceMessage: fmt.Sprintf(format, args...),
	}
	err.addCaller(1)
	return err
}

// Trace attaches a location to the error, unless the error is nil.
func Trace(other error) error {
	if other == nil {
		return nil
	}
	err := &wrappedError{
		previous: other,
	}
	err.addCaller(1)
	return err
}

// Tracef attaches a location and a message to the error, unless the error is
// nil.
func Tracef(other error, format string, args ...interface{}) error {
	if other == nil {
		return nil
	}
	err := &wrappedError{
		traceMessage: fmt.Sprintf(format, args...),
		previous:     other,
	}
	err.addCaller(1)
	return err
}

// Cause returns the underlying raw error.
func Cause(err error) error {
	if err, ok := err.(*wrappedError); ok {
		return err.Cause()
	}
	return err
}

// GetTrace returns a list of trace information attached to this error.
func GetTrace(err error) []string {
	if err == nil {
		return []string{}
	}

	var lines []string
	for err != nil {
		buf := bytes.Buffer{}
		if e, ok := err.(*wrappedError); ok {
			// add the trace info to this line
			file, line := e.getCaller()
			if file != "" {
				buf.WriteString(fmt.Sprintf("%s:%d", file, line))
			}
			buf.WriteString(" [error] ")
			buf.WriteString(e.traceMessage)
			err = e.previous
		} else {
			// just the error message
			buf.WriteString(err.Error())
			err = nil
		}
		lines = append(lines, buf.String())
	}

	// reverse the list
	var result []string
	for i := len(lines); i > 0; i-- {
		result = append(result, lines[i-1])
	}
	return result
}

// Details returns a formatted trace string.
func Details(err error) string {
	return strings.Join(GetTrace(err), "\n")
}

type wrappedError struct {
	traceFile    string
	traceLine    int
	traceMessage string

	previous error
}

// Cause returns the underlying error.
func (e *wrappedError) Cause() error {
	if e.previous == nil {
		return e
	}
	switch err := e.previous.(type) {
	case *wrappedError:
		return err.Cause()
	default:
		return err
	}
}

// Error returns the error message of the underlying error.
func (e *wrappedError) Error() string {
	if e.previous == nil {
		return e.traceMessage
	}
	return e.Cause().Error()
}

// Trace returns a list of trace information attached to this error.
func (e *wrappedError) Trace() []string {
	return GetTrace(e)
}

// Details returns a formatted trace string.
func (e *wrappedError) Details() string {
	return Details(e)
}

func (e *wrappedError) addCaller(callDepth int) {
	_, file, line, _ := runtime.Caller(callDepth + 1)
	e.traceFile = file
	e.traceLine = line
}

func (e *wrappedError) getCaller() (filename string, line int) {
	return e.traceFile, e.traceLine
}

// enforce wrappedError follows the interface
var _ Error = Error(&wrappedError{})
