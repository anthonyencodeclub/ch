package cmd_test

import (
	"testing"

	"github.com/anthonyencodeclub/ch/internal/cmd"
	"github.com/anthonyencodeclub/ch/internal/config"
)

func TestExecute_SetupHelp(t *testing.T) {
	err := cmd.Execute([]string{"setup", "--help"})
	if err != nil {
		t.Fatalf("Execute(setup --help) error: %v", err)
	}
}

func TestSetup_ConfigPersistsDefaultCompany(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CH_CONFIG_DIR", tmp)

	// Manually set a default company to verify the config round-trip
	cfg := config.File{
		APIKey:         "test-key",
		DefaultCompany: "00445790",
		CompanyName:    "TESCO PLC",
	}
	if err := config.WriteConfig(cfg); err != nil {
		t.Fatalf("WriteConfig() error: %v", err)
	}

	got, err := config.ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig() error: %v", err)
	}
	if got.DefaultCompany != "00445790" {
		t.Errorf("DefaultCompany = %q, want %q", got.DefaultCompany, "00445790")
	}
	if got.CompanyName != "TESCO PLC" {
		t.Errorf("CompanyName = %q, want %q", got.CompanyName, "TESCO PLC")
	}
}
