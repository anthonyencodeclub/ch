package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/anthonyencodeclub/ch/internal/config"
)

func TestReadConfig_MissingFile(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CH_CONFIG_DIR", tmp)

	cfg, err := config.ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig() error: %v", err)
	}
	if cfg.APIKey != "" {
		t.Errorf("expected empty APIKey, got %q", cfg.APIKey)
	}
}

func TestWriteAndReadConfig(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CH_CONFIG_DIR", tmp)

	want := config.File{APIKey: "test-key-123"}
	if err := config.WriteConfig(want); err != nil {
		t.Fatalf("WriteConfig() error: %v", err)
	}

	got, err := config.ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig() error: %v", err)
	}
	if got.APIKey != want.APIKey {
		t.Errorf("APIKey = %q, want %q", got.APIKey, want.APIKey)
	}
}

func TestWriteConfig_AtomicWrite(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CH_CONFIG_DIR", tmp)

	cfg := config.File{APIKey: "atomic-test"}
	if err := config.WriteConfig(cfg); err != nil {
		t.Fatalf("WriteConfig() error: %v", err)
	}

	// No .tmp file should remain
	tmpFile := filepath.Join(tmp, "config.json.tmp")
	if _, err := os.Stat(tmpFile); !os.IsNotExist(err) {
		t.Error("temp file should not exist after WriteConfig")
	}

	// config.json should exist
	configFile := filepath.Join(tmp, "config.json")
	if _, err := os.Stat(configFile); err != nil {
		t.Errorf("config.json should exist: %v", err)
	}
}

func TestConfigPath(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CH_CONFIG_DIR", tmp)

	path, err := config.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error: %v", err)
	}
	want := filepath.Join(tmp, "config.json")
	if path != want {
		t.Errorf("ConfigPath() = %q, want %q", path, want)
	}
}

func TestDir_FromEnv(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CH_CONFIG_DIR", tmp)

	dir, err := config.Dir()
	if err != nil {
		t.Fatalf("Dir() error: %v", err)
	}
	if dir != tmp {
		t.Errorf("Dir() = %q, want %q", dir, tmp)
	}
}

func TestEnsureDir_CreatesDirectory(t *testing.T) {
	tmp := t.TempDir()
	newDir := filepath.Join(tmp, "subdir", "config")
	t.Setenv("CH_CONFIG_DIR", newDir)

	got, err := config.EnsureDir()
	if err != nil {
		t.Fatalf("EnsureDir() error: %v", err)
	}
	if got != newDir {
		t.Errorf("EnsureDir() = %q, want %q", got, newDir)
	}

	info, err := os.Stat(newDir)
	if err != nil {
		t.Fatalf("directory should exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("should be a directory")
	}
}

func TestAPIKey_FromEnv(t *testing.T) {
	t.Setenv("CH_API_KEY", "env-key-456")

	key, err := config.APIKey()
	if err != nil {
		t.Fatalf("APIKey() error: %v", err)
	}
	if key != "env-key-456" {
		t.Errorf("APIKey() = %q, want %q", key, "env-key-456")
	}
}

func TestAPIKey_FromConfig(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CH_CONFIG_DIR", tmp)
	t.Setenv("CH_API_KEY", "")

	cfg := config.File{APIKey: "config-key-789"}
	if err := config.WriteConfig(cfg); err != nil {
		t.Fatalf("WriteConfig() error: %v", err)
	}

	key, err := config.APIKey()
	if err != nil {
		t.Fatalf("APIKey() error: %v", err)
	}
	if key != "config-key-789" {
		t.Errorf("APIKey() = %q, want %q", key, "config-key-789")
	}
}

func TestAPIKey_NotConfigured(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CH_CONFIG_DIR", tmp)
	t.Setenv("CH_API_KEY", "")

	_, err := config.APIKey()
	if err == nil {
		t.Fatal("APIKey() should return error when not configured")
	}
}

func TestReadConfig_InvalidJSON(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CH_CONFIG_DIR", tmp)

	configFile := filepath.Join(tmp, "config.json")
	if err := os.WriteFile(configFile, []byte("{invalid json}"), 0o600); err != nil {
		t.Fatalf("write invalid config: %v", err)
	}

	_, err := config.ReadConfig()
	if err == nil {
		t.Fatal("ReadConfig() should return error for invalid JSON")
	}
}

func TestWriteConfig_OverwritesExisting(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CH_CONFIG_DIR", tmp)

	cfg1 := config.File{APIKey: "first-key"}
	if err := config.WriteConfig(cfg1); err != nil {
		t.Fatalf("WriteConfig(1) error: %v", err)
	}

	cfg2 := config.File{APIKey: "second-key"}
	if err := config.WriteConfig(cfg2); err != nil {
		t.Fatalf("WriteConfig(2) error: %v", err)
	}

	got, err := config.ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig() error: %v", err)
	}
	if got.APIKey != "second-key" {
		t.Errorf("APIKey = %q, want %q", got.APIKey, "second-key")
	}
}
