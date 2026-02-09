package cmd_test

import (
	"errors"
	"testing"

	"github.com/anthonyencodeclub/ch/internal/cmd"
)

func TestExitCode_Nil(t *testing.T) {
	if got := cmd.ExitCode(nil); got != 0 {
		t.Errorf("ExitCode(nil) = %d, want 0", got)
	}
}

func TestExitCode_ExitError(t *testing.T) {
	err := &cmd.ExitError{Code: 2, Err: errors.New("usage error")}
	if got := cmd.ExitCode(err); got != 2 {
		t.Errorf("ExitCode(ExitError{2}) = %d, want 2", got)
	}
}

func TestExitCode_NegativeCode(t *testing.T) {
	err := &cmd.ExitError{Code: -1, Err: errors.New("negative")}
	if got := cmd.ExitCode(err); got != 1 {
		t.Errorf("ExitCode(ExitError{-1}) = %d, want 1", got)
	}
}

func TestExitCode_GenericError(t *testing.T) {
	err := errors.New("generic error")
	if got := cmd.ExitCode(err); got != 1 {
		t.Errorf("ExitCode(generic) = %d, want 1", got)
	}
}

func TestExitError_Error(t *testing.T) {
	err := &cmd.ExitError{Code: 1, Err: errors.New("something failed")}
	if got := err.Error(); got != "something failed" {
		t.Errorf("Error() = %q, want %q", got, "something failed")
	}
}

func TestExitError_NilError(t *testing.T) {
	err := &cmd.ExitError{Code: 0, Err: nil}
	if got := err.Error(); got != "" {
		t.Errorf("Error() = %q, want empty", got)
	}
}

func TestExitError_NilReceiver(t *testing.T) {
	var err *cmd.ExitError
	if got := err.Error(); got != "" {
		t.Errorf("nil.Error() = %q, want empty", got)
	}
}

func TestExitError_Unwrap(t *testing.T) {
	inner := errors.New("inner")
	err := &cmd.ExitError{Code: 1, Err: inner}
	if got := err.Unwrap(); got != inner {
		t.Errorf("Unwrap() = %v, want %v", got, inner)
	}
}

func TestExitError_UnwrapNil(t *testing.T) {
	var err *cmd.ExitError
	if got := err.Unwrap(); got != nil {
		t.Errorf("nil.Unwrap() = %v, want nil", got)
	}
}
