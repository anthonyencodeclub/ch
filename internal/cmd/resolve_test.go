package cmd

import (
	"testing"

	"github.com/anthonyencodeclub/ch/internal/config"
)

func TestResolveCompanyNumber_Explicit(t *testing.T) {
	got, err := resolveCompanyNumber("12345678")
	if err != nil {
		t.Fatalf("resolveCompanyNumber() error: %v", err)
	}
	if got != "12345678" {
		t.Errorf("resolveCompanyNumber() = %q, want %q", got, "12345678")
	}
}

func TestResolveCompanyNumber_FromConfig(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CH_CONFIG_DIR", tmp)

	cfg := config.File{DefaultCompany: "00445790"}
	if err := config.WriteConfig(cfg); err != nil {
		t.Fatalf("WriteConfig() error: %v", err)
	}

	got, err := resolveCompanyNumber("")
	if err != nil {
		t.Fatalf("resolveCompanyNumber() error: %v", err)
	}
	if got != "00445790" {
		t.Errorf("resolveCompanyNumber() = %q, want %q", got, "00445790")
	}
}

func TestResolveCompanyNumber_NoDefault(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CH_CONFIG_DIR", tmp)

	_, err := resolveCompanyNumber("")
	if err == nil {
		t.Fatal("resolveCompanyNumber() should return error when no default set")
	}
}

func TestResolveCompanyNumber_ExplicitOverridesDefault(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CH_CONFIG_DIR", tmp)

	cfg := config.File{DefaultCompany: "00445790"}
	if err := config.WriteConfig(cfg); err != nil {
		t.Fatalf("WriteConfig() error: %v", err)
	}

	got, err := resolveCompanyNumber("99999999")
	if err != nil {
		t.Fatalf("resolveCompanyNumber() error: %v", err)
	}
	if got != "99999999" {
		t.Errorf("resolveCompanyNumber() = %q, want %q, explicit should override default", got, "99999999")
	}
}
