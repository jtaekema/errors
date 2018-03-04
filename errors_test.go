package errors

import (
	"fmt"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	err, ok := New("my test error").(Error)
	if !ok {
		t.Fatal("New must return an `Error`")
	}

	if Cause(err) == nil {
		t.Errorf("New error must have a cause")
	}

	if err.Error() != "my test error" {
		t.Error("New error must preserve the error message")
	}

	trace := err.Trace()
	if len(trace) != 1 {
		t.Error("New error must have a trace")
	}
	if !strings.HasSuffix(trace[0], "my test error") {
		t.Error("New error must have error message in the trace")
	}

	details := err.Details()
	if len(details) == 0 {
		t.Error("New error must return trace details")
	}
}

func TestNewF(t *testing.T) {
	err, ok := Newf("formatted: '%s'", "message").(Error)
	if !ok {
		t.Fatal("Newf must return an `Error`")
	}

	if Cause(err) == nil {
		t.Error("Newf error must have a cause")
	}

	if err.Error() != "formatted: 'message'" {
		t.Error("Newf error must preserve the error message")
	}

	trace := err.Trace()
	if len(trace) != 1 {
		t.Error("Newf error must have a trace")
	}
	if !strings.HasSuffix(trace[0], "formatted: 'message'") {
		t.Error("Newf error must have error message in the trace")
	}

	details := err.Details()
	if len(details) == 0 {
		t.Error("Newf error must return trace details")
	}
}

func TestTrace(t *testing.T) {
	err, ok := Trace(New("first")).(Error)
	if !ok {
		t.Fatal("Trace must return an `Error`")
	}

	if Cause(err) == nil {
		t.Error("Traced error must have a cause")
	}

	if err.Error() != "first" {
		t.Error("Traced error must preserve the error message")
	}

	trace := err.Trace()
	if len(trace) != 2 {
		t.Error("Traced error must have the error message from two callers")
	}
	if !strings.HasSuffix(trace[0], "first") {
		t.Error("Traced error must have error message in the trace")
	}

	details := err.Details()
	if len(details) == 0 {
		t.Error("Traced error must return trace details")
	}
}

func TestTracef(t *testing.T) {
	err, ok := Tracef(New("first"), "my trace message").(Error)
	if !ok {
		t.Fatal("Tracef must return an `Error`")
	}

	if Cause(err) == nil {
		t.Error("Tracef error must have a cause")
	}

	if err.Error() != "first" {
		t.Error("Tracef error must preserve the error message")
	}

	trace := err.Trace()
	if len(trace) != 2 {
		t.Fatal("Tracef error must have the error message from two callers")
	}
	if !strings.HasSuffix(trace[0], "first") {
		t.Error("Tracef error must have error message in the trace")
	}
	if !strings.HasSuffix(trace[1], "my trace message") {
		t.Error("Tracef error must have trace message in the trace")
	}

	details := err.Details()
	if len(details) == 0 {
		t.Error("Tracef error must return trace details")
	}
}

func TestTracefStandardError(t *testing.T) {
	err, ok := Tracef(fmt.Errorf("standard"), "my trace message").(Error)
	if !ok {
		t.Fatal("Tracef must return an `Error`")
	}

	if Cause(err) == nil {
		t.Error("Tracef error must have a cause")
	}

	if err.Error() != "standard" {
		t.Error("Tracef error must preserve the error message")
	}

	trace := err.Trace()
	if len(trace) != 2 {
		t.Fatal("Tracef error must have the error message from two callers")
	}
	if !strings.HasSuffix(trace[0], "standard") {
		t.Error("Tracef error must have error message in the trace")
	}
	if !strings.HasSuffix(trace[1], "my trace message") {
		t.Error("Tracef error must have trace message in the trace")
	}

	details := err.Details()
	if len(details) == 0 {
		t.Error("Tracef error must return trace details")
	}
}

func TestTraceNilError(t *testing.T) {
	err := Trace(nil)
	if err != nil {
		t.Error("Tracef must return nil if the error was nil")
	}
}

func TestTracefNilError(t *testing.T) {
	err := Tracef(nil, "error")
	if err != nil {
		t.Error("Tracef must return nil if the error was nil")
	}
}

func TestCauseStandardError(t *testing.T) {
	standard := fmt.Errorf("standard")
	err := Cause(standard)
	if err == nil {
		t.Error("Cause must return a non-nil error")
	}
	if err != standard {
		t.Error("Cause must return the underlying error")
	}
	wrapped := Trace(standard)
	err = Cause(wrapped)
	if err == nil {
		t.Error("Cause must return a non-nil error")
	}
	if err != standard {
		t.Error("Cause must return the underlying error")
	}
}

func TestCauseNilError(t *testing.T) {
	err := Cause(nil)
	if err != nil {
		t.Error("Cause must return nil if the error was nil")
	}
}

func TestDetailsNilError(t *testing.T) {
	details := Details(nil)
	if details != "" {
		t.Error("Details must return nothing for a nil error")
	}
}

func TestGetTraceNilError(t *testing.T) {
	trace := GetTrace(nil)
	if len(trace) > 0 {
		t.Error("GetTrace must return nothing for a nil error")
	}
}
