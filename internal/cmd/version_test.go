package cmd_test

import (
	"testing"

	"github.com/anthonyencodeclub/ch/internal/cmd"
)

func TestVersionString(t *testing.T) {
	v := cmd.VersionString()
	if v == "" {
		t.Fatal("VersionString() returned empty string")
	}
	// Should contain a version number
	if v == "dev" {
		// That's fine for a dev build
		return
	}
	// Should start with a digit or 'v'
	if v[0] != 'v' && (v[0] < '0' || v[0] > '9') {
		t.Errorf("VersionString() = %q, unexpected format", v)
	}
}
