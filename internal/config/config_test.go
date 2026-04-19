package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if !cfg.Mouse {
		t.Error("expected Mouse true")
	}
	if !cfg.AltScreen {
		t.Error("expected AltScreen true")
	}
	if cfg.Theme != "default" {
		t.Errorf("expected Theme 'default', got %q", cfg.Theme)
	}
}

func TestLoadMissingFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cfg, err := Load()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !cfg.Mouse {
		t.Error("expected defaults from missing file")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	cfg := Config{Mouse: false, AltScreen: false, Theme: "dracula"}
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	expectedPath := filepath.Join(dir, ".config", "rfdtui", "config.yaml")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Error("config file not created")
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if loaded.Mouse != false {
		t.Error("expected Mouse false")
	}
	if loaded.Theme != "dracula" {
		t.Errorf("expected Theme 'dracula', got %q", loaded.Theme)
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	cfgDir := filepath.Join(dir, ".config", "rfdtui")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, "config.yaml"), []byte(":\tinvalid"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Load()
	if err == nil {
		t.Error("expected error from invalid YAML")
	}
}
