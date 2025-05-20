// Package errors extends standard errors with caller information as well as the
// ability to explicitly wrap errors.
package errors

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// IncludeCaller if true will include filename tracing to the error logs. This
// is useful for debugging the program that is creating the error, while leaving
// this off is better for users of the program.
var IncludeCaller bool

// New creates a new error and adds trace information to it.
func New(text string) error {
	err := &wrappedError{
		message: text,
	}
	if IncludeCaller {
		err.addCaller(1)
	}
	return err
}

// Newf creates a new error and adds trace information to it.
func Newf(format string, args ...any) error {
	err := &wrappedError{
		message: fmt.Sprintf(format, args...),
	}
	if IncludeCaller {
		err.addCaller(1)
	}
	return err
}

// Cause returns the underlying original error.
func Cause(err error) error {
	for {
		e, ok := err.(*wrappedError)
		if ok && e.err != nil {
			err = e.err
		} else {
			break
		}
	}
	return err
}

// Trace returns a list of trace information attached to this error.
func Trace(err error) []string {
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
			buf.WriteString(e.message)
			err = e.err
		} else {
			// just the error message
			buf.WriteString(err.Error())
			err = nil
		}
		lines = append(lines, buf.String())
	}

	// reverse the list so that the first error in the slice is the original
	// underlying error.
	var result []string
	for i := len(lines); i > 0; i-- {
		result = append(result, lines[i-1])
	}
	return result
}

// Details returns a formatted trace string suitable for printing an error trace.
func Details(err error) string {
	return strings.Join(Trace(err), "\n")
}

// Is reports whether any error in err's tree matches target.
//
// For more information see: https://pkg.go.dev/errors#Is
func Is(err, target error) bool {
	// Use Is from the errors package.
	return errors.Is(err, target)
}

// As finds the first error in err's tree that matches target, and if one is found, sets
// target to that error value and returns true. Otherwise, it returns false.
//
// For more information see: https://pkg.go.dev/errors#As
func As(err error, target any) bool {
	// Use As from the errors package.
	return errors.As(err, target)
}

// Wrap returns a wrapped error that may include caller information, unless the
// provided error is nil.
func Wrap(err error, text string) error {
	if err == nil {
		return nil
	}

	new := &wrappedError{
		message: text,
		err:     err,
	}
	if IncludeCaller {
		new.addCaller(1)
	}

	return new
}

// Wrapf returns a wrapped error that may include caller information, unless the
// provided error is nil.
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	new := &wrappedError{
		message: fmt.Sprintf(format, args...),
		err:     err,
	}
	if IncludeCaller {
		new.addCaller(1)
	}

	return new
}

func Unwrap(err error) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*wrappedError); ok && e.err != nil {
		return e.err
	}

	return err
}

type wrappedError struct {
	file    string
	line    int
	message string

	err error
}

// Error returns the error message of the underlying error.
func (e *wrappedError) Error() string {
	if e.err == nil {
		return e.message
	}

	return Unwrap(e).Error()
}

func (e *wrappedError) Unwrap() error {
	return Unwrap(e)
}

func (e *wrappedError) addCaller(callDepth int) {
	_, file, line, _ := runtime.Caller(callDepth + 1)
	e.file = file
	e.line = line
}

func (e *wrappedError) getCaller() (file string, line int) {
	return e.file, e.line
}
