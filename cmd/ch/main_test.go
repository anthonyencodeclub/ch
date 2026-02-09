package main_test

import (
	"os"
	"os/exec"
	"testing"
)

func TestMain_Build(t *testing.T) {
	// Verify the binary builds successfully
	cmd := exec.Command("go", "build", "-o", os.DevNull, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
}
