package ui_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/anthonyencodeclub/ch/internal/ui"
)

func newTestUI(t *testing.T) (*ui.UI, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	u, err := ui.New(ui.Options{
		Stdout: stdout,
		Stderr: stderr,
		Color:  "never",
	})
	if err != nil {
		t.Fatalf("ui.New() error: %v", err)
	}
	return u, stdout, stderr
}

func TestNew_Defaults(t *testing.T) {
	u, err := ui.New(ui.Options{Color: "never"})
	if err != nil {
		t.Fatalf("ui.New() error: %v", err)
	}
	if u == nil {
		t.Fatal("ui.New() returned nil")
	}
}

func TestNew_NeverColor(t *testing.T) {
	u, _, _ := newTestUI(t)
	if u == nil {
		t.Fatal("ui.New() returned nil")
	}
}

func TestSuccess(t *testing.T) {
	u, _, stderr := newTestUI(t)
	u.Success("all good")
	if !bytes.Contains(stderr.Bytes(), []byte("all good")) {
		t.Errorf("Success() output = %q, should contain 'all good'", stderr.String())
	}
}

func TestWarn(t *testing.T) {
	u, _, stderr := newTestUI(t)
	u.Warn("be careful")
	if !bytes.Contains(stderr.Bytes(), []byte("be careful")) {
		t.Errorf("Warn() output = %q, should contain 'be careful'", stderr.String())
	}
}

func TestError(t *testing.T) {
	u, _, stderr := newTestUI(t)
	u.Error("something broke")
	if !bytes.Contains(stderr.Bytes(), []byte("something broke")) {
		t.Errorf("Error() output = %q, should contain 'something broke'", stderr.String())
	}
}

func TestInfo(t *testing.T) {
	u, _, stderr := newTestUI(t)
	u.Info("fyi")
	if !bytes.Contains(stderr.Bytes(), []byte("fyi")) {
		t.Errorf("Info() output = %q, should contain 'fyi'", stderr.String())
	}
}

func TestContext_RoundTrip(t *testing.T) {
	u, _, _ := newTestUI(t)
	ctx := ui.WithUI(context.Background(), u)

	got := ui.FromContext(ctx)
	if got == nil {
		t.Fatal("FromContext() returned nil")
	}
	if got != u {
		t.Error("FromContext() returned different UI instance")
	}
}

func TestFromContext_Empty(t *testing.T) {
	got := ui.FromContext(context.Background())
	if got != nil {
		t.Error("FromContext(empty) should return nil")
	}
}

func TestStdout(t *testing.T) {
	u, stdout, _ := newTestUI(t)
	w := u.Stdout()
	if w == nil {
		t.Fatal("Stdout() returned nil")
	}
	if w != stdout {
		t.Error("Stdout() should return the configured stdout writer")
	}
}

func TestStderr(t *testing.T) {
	u, _, stderr := newTestUI(t)
	w := u.Stderr()
	if w == nil {
		t.Fatal("Stderr() returned nil")
	}
	if w != stderr {
		t.Error("Stderr() should return the configured stderr writer")
	}
}

func TestOutput(t *testing.T) {
	u, _, _ := newTestUI(t)
	out := u.Output()
	if out == nil {
		t.Fatal("Output() returned nil")
	}
}
