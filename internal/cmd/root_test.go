package cmd_test

import (
	"testing"

	"github.com/anthonyencodeclub/ch/internal/cmd"
)

func TestExecute_Help(t *testing.T) {
	err := cmd.Execute([]string{"--help"})
	// --help triggers exit(0) which becomes nil
	if err != nil {
		t.Fatalf("Execute(--help) error: %v", err)
	}
}

func TestExecute_Version(t *testing.T) {
	err := cmd.Execute([]string{"version"})
	if err != nil {
		t.Fatalf("Execute(version) error: %v", err)
	}
}

func TestExecute_UnknownCommand(t *testing.T) {
	err := cmd.Execute([]string{"nonexistent-command"})
	if err == nil {
		t.Fatal("Execute(nonexistent) should return error")
	}
}

func TestExecute_AuthStatusNoKey(t *testing.T) {
	t.Setenv("CH_API_KEY", "")
	t.Setenv("CH_CONFIG_DIR", t.TempDir())

	// Should not error — just shows "not authenticated"
	err := cmd.Execute([]string{"auth", "status"})
	if err != nil {
		t.Fatalf("Execute(auth status) error: %v", err)
	}
}

func TestExecute_AuthSetKey(t *testing.T) {
	t.Setenv("CH_CONFIG_DIR", t.TempDir())

	err := cmd.Execute([]string{"auth", "set-key", "test-key-123"})
	if err != nil {
		t.Fatalf("Execute(auth set-key) error: %v", err)
	}
}

func TestExecute_VersionJSON(t *testing.T) {
	err := cmd.Execute([]string{"version", "--json"})
	if err != nil {
		t.Fatalf("Execute(version --json) error: %v", err)
	}
}

func TestExecute_ConflictingOutputFlags(t *testing.T) {
	err := cmd.Execute([]string{"version", "--json", "--plain"})
	if err == nil {
		t.Fatal("Execute(--json --plain) should return error")
	}
}
