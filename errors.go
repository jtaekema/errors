package errors

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

// IncludeCaller if true will include filename tracing to the error logs. This
// is useful for debugging the program that is creating the error, while leaving
// this off is better for users of the program.
var IncludeCaller bool

// Error interface.
type Error interface {
	// Cause returns the underlying error.
	Cause() error

	// Error returns the error message of the underlying error.
	Error() string

	// Trace returns a list of trace information attached to this error.
	Trace() []string

	// Details returns a formatted trace string.
	Details() string
}

// New creates a new raw error and trace it.
func New(text string) error {
	err := &wrappedError{
		traceMessage: text,
	}
	if IncludeCaller {
		err.addCaller(1)
	}
	return err
}

// Newf creates a new raw error and trace it.
func Newf(format string, args ...interface{}) error {
	err := &wrappedError{
		traceMessage: fmt.Sprintf(format, args...),
	}
	if IncludeCaller {
		err.addCaller(1)
	}
	return err
}

// Trace returns a wrapped error that may include caller information, unless the
// previous error is nil.
func Trace(previous error) error {
	if previous == nil {
		return nil
	}
	err := &wrappedError{
		previous: previous,
	}
	if IncludeCaller {
		err.addCaller(1)
	}
	return err
}

// Tracef returns a wrapped error that may include caller information, unless
// the previous error is nil.
func Tracef(previous error, format string, args ...interface{}) error {
	if previous == nil {
		return nil
	}
	err := &wrappedError{
		traceMessage: fmt.Sprintf(format, args...),
		previous:     previous,
	}
	if IncludeCaller {
		err.addCaller(1)
	}
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
				buf.WriteString(fmt.Sprintf("%s:%d ", file, line))
			}
			buf.WriteString("[error] ")
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

// ensure wrappedError follows the interface
var _ Error = Error(&wrappedError{})
