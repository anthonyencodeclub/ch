package errfmt_test

import (
	"errors"
	"testing"

	"github.com/anthonyencodeclub/ch/internal/errfmt"
)

func TestFormat_NilError(t *testing.T) {
	got := errfmt.Format(nil)
	if got != "" {
		t.Errorf("Format(nil) = %q, want empty string", got)
	}
}

func TestFormat_SimpleError(t *testing.T) {
	err := errors.New("something went wrong")
	got := errfmt.Format(err)
	want := "Error: something went wrong"
	if got != want {
		t.Errorf("Format(%v) = %q, want %q", err, got, want)
	}
}

func TestFormat_WrappedError(t *testing.T) {
	inner := errors.New("connection refused")
	err := errors.Join(errors.New("api call failed"), inner)
	got := errfmt.Format(err)
	if got == "" {
		t.Error("Format(wrapped) returned empty string")
	}
	if got[:6] != "Error:" {
		t.Errorf("Format(wrapped) should start with 'Error:', got %q", got)
	}
}
