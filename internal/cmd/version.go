package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/openclaw/ch/internal/outfmt"
)

var (
	version = "0.1.0"
	commit  = ""
	date    = ""
)

// VersionString returns a human-readable version string.
func VersionString() string {
	v := strings.TrimSpace(version)
	if v == "" {
		v = "dev"
	}
	if strings.TrimSpace(commit) == "" && strings.TrimSpace(date) == "" {
		return v
	}
	if strings.TrimSpace(commit) == "" {
		return fmt.Sprintf("%s (%s)", v, strings.TrimSpace(date))
	}
	if strings.TrimSpace(date) == "" {
		return fmt.Sprintf("%s (%s)", v, strings.TrimSpace(commit))
	}
	return fmt.Sprintf("%s (%s %s)", v, strings.TrimSpace(commit), strings.TrimSpace(date))
}

// VersionCmd prints the version.
type VersionCmd struct{}

func (c *VersionCmd) Run(ctx context.Context) error {
	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]any{
			"version": strings.TrimSpace(version),
			"commit":  strings.TrimSpace(commit),
			"date":    strings.TrimSpace(date),
		})
	}
	fmt.Fprintln(os.Stdout, VersionString())
	return nil
}
