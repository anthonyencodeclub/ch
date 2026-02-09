package outfmt

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// Mode controls how output is rendered.
type Mode struct {
	JSON  bool
	Plain bool
}

// ParseError is returned when output flags conflict.
type ParseError struct{ msg string }

func (e *ParseError) Error() string { return e.msg }

// FromFlags validates and returns a Mode from CLI flags.
func FromFlags(jsonOut bool, plainOut bool) (Mode, error) {
	if jsonOut && plainOut {
		return Mode{}, &ParseError{msg: "invalid output mode (cannot combine --json and --plain)"}
	}
	return Mode{JSON: jsonOut, Plain: plainOut}, nil
}

// FromEnv reads output mode from environment variables.
func FromEnv() Mode {
	return Mode{
		JSON:  envBool("CH_JSON"),
		Plain: envBool("CH_PLAIN"),
	}
}

type ctxKey struct{}

// WithMode stores the output mode in the context.
func WithMode(ctx context.Context, mode Mode) context.Context {
	return context.WithValue(ctx, ctxKey{}, mode)
}

// FromContext retrieves the output mode from the context.
func FromContext(ctx context.Context) Mode {
	if v := ctx.Value(ctxKey{}); v != nil {
		if m, ok := v.(Mode); ok {
			return m
		}
	}
	return Mode{}
}

// IsJSON returns true if JSON output is enabled.
func IsJSON(ctx context.Context) bool { return FromContext(ctx).JSON }

// IsPlain returns true if plain output is enabled.
func IsPlain(ctx context.Context) bool { return FromContext(ctx).Plain }

// WriteJSON encodes a value as pretty-printed JSON.
func WriteJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func envBool(key string) bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	switch v {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}
