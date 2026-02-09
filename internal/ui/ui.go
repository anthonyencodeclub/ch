package ui

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/muesli/termenv"
)

type ctxKey struct{}

// Options configures the UI.
type Options struct {
	Stdout io.Writer
	Stderr io.Writer
	Color  string
}

// UI provides styled terminal output.
type UI struct {
	stdout io.Writer
	stderr io.Writer
	out    *termenv.Output
}

// New creates a new UI instance.
func New(opts Options) (*UI, error) {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	profile := termenv.ColorProfile()
	if opts.Color == "never" {
		profile = termenv.Ascii
	}

	out := termenv.NewOutput(stdout, termenv.WithProfile(profile))

	return &UI{
		stdout: stdout,
		stderr: stderr,
		out:    out,
	}, nil
}

// WithUI stores the UI in the context.
func WithUI(ctx context.Context, u *UI) context.Context {
	return context.WithValue(ctx, ctxKey{}, u)
}

// FromContext retrieves the UI from the context.
func FromContext(ctx context.Context) *UI {
	if v := ctx.Value(ctxKey{}); v != nil {
		if u, ok := v.(*UI); ok {
			return u
		}
	}
	return nil
}

// Stdout returns the stdout writer.
func (u *UI) Stdout() io.Writer { return u.stdout }

// Stderr returns the stderr writer.
func (u *UI) Stderr() io.Writer { return u.stderr }

// Output returns the termenv output for styled printing.
func (u *UI) Output() *termenv.Output { return u.out }

// Success prints a success message to stderr.
func (u *UI) Success(msg string) {
	fmt.Fprintln(u.stderr, u.out.String(msg).Foreground(u.out.Color("2")))
}

// Warn prints a warning message to stderr.
func (u *UI) Warn(msg string) {
	fmt.Fprintln(u.stderr, u.out.String(msg).Foreground(u.out.Color("3")))
}

// Error prints an error message to stderr.
func (u *UI) Error(msg string) {
	fmt.Fprintln(u.stderr, u.out.String(msg).Foreground(u.out.Color("1")))
}

// Info prints an info message to stderr.
func (u *UI) Info(msg string) {
	fmt.Fprintln(u.stderr, u.out.String(msg).Foreground(u.out.Color("6")))
}
