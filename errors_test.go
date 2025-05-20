package errors

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	err := New("my test error")
	if err.Error() != "my test error" {
		t.Error("New error must preserve the error message")
	}
}

func TestNewf(t *testing.T) {
	err := Newf("formatted: '%s'", "message")
	if err.Error() != "formatted: 'message'" {
		t.Error("Newf error must preserve the error message")
	}
}

func TestCause(t *testing.T) {
	err := New("error")
	wrapped := Wrapf(err, "first wrap")
	wrapped2 := Wrapf(wrapped, "second wrap")

	fmtError := fmt.Errorf("error")
	fmtWrapped := Wrap(fmtError, "first wrap")
	fmtWrapped2 := Wrap(fmtWrapped, "second wrap")

	stdError := errors.New("error")
	stdWrapped := Wrap(stdError, "first wrap")
	stdWrapped2 := Wrap(stdWrapped, "second wrap")

	tests := []struct {
		err      error
		expected error
		match    bool
	}{
		// self
		{nil, nil, true},
		{err, err, true},
		{fmtError, fmtError, true},
		{stdError, stdError, true},
		{fmtWrapped, fmtWrapped, false},
		{stdWrapped, stdWrapped, false},
		{wrapped, wrapped, false},
		{wrapped2, wrapped2, false},

		// new
		{err, New("error"), false},
		{fmtError, fmt.Errorf("error"), false},
		{stdError, errors.New("error"), false},
		{fmtWrapped, Wrap(fmtError, "wrapped"), false},
		{fmtWrapped2, err, false},
		{stdWrapped, Wrap(stdError, "wrapped"), false},
		{stdWrapped2, err, false},

		// original
		{wrapped, err, true},
		{wrapped2, err, true},
		{fmtWrapped, fmtError, true},
		{fmtWrapped2, fmtError, true},
		{stdWrapped, stdError, true},
		{stdWrapped2, stdError, true},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			cause := Cause(tc.err)
			if (cause == tc.expected) != tc.match {
				t.Errorf("Cause(%q) => %v, expected %q match %v", tc.err, cause, tc.expected, tc.match)
			}
		})
	}
}

func TestTrace(t *testing.T) {
	trace := Trace(nil)
	if len(trace) > 0 {
		t.Error("GetTrace must return nothing for a nil error")
	}
}

func TestDetails(t *testing.T) {
	details := Details(nil)
	if details != "" {
		t.Error("Details must return nothing for a nil error")
	}
}

func TestWrap(t *testing.T) {
	err := Wrap(nil, "nothing")
	if err != nil {
		t.Error("Wrap must return nil if the error was nil")
	}

	original := New("error")
	wrapped := Wrap(original, "first wrap")

	if wrapped.Error() != "error" {
		t.Error("Wrapped error must preserve the error message")
	}
}

func TestWrapf(t *testing.T) {
	err := Wrapf(nil, "nil error")
	if err != nil {
		t.Error("Wrapf must return nil if the error was nil")
	}

	tests := []struct {
		err      error
		traceMsg string
		errMsg   string
	}{
		{New("first"), "wrapping message", "first"},
		{fmt.Errorf("standard"), "extra context", "standard"},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {

			e := Wrap(tc.err, tc.traceMsg)
			if e.Error() != tc.errMsg {
				t.Error("Wrapf must preserve the original underlying error message")
			}

			errors := Trace(e)
			if len(errors) != 2 {
				t.Fatal("Wrapf error must have trace information from original and wrapped errors")
			}
			if !strings.HasSuffix(errors[0], tc.errMsg) {
				t.Error("Wrapf error must have error message in the trace")
			}
			if !strings.HasSuffix(errors[1], tc.traceMsg) {
				t.Errorf("Wrapf error must have trace message in the trace: got %q expected %q", errors[1], tc.traceMsg)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	err := New("error")
	wrapped := Wrap(err, "test")

	if unwrapped := Unwrap(nil); unwrapped != nil {
		t.Errorf("Unwrap(%v) = %v, want 'nil'", err, unwrapped)
	}

	if unwrapped := Unwrap(err); unwrapped != err {
		t.Errorf("Unwrap(%v) = %v, want %q", err, unwrapped, err)
	}

	if unwrapped := Unwrap(wrapped); unwrapped != err {
		t.Errorf("Unwrap(%v) = %v, want %q", wrapped, unwrapped, err)
	}
}
