package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Mouse     bool   `yaml:"mouse"`
	AltScreen bool   `yaml:"alt_screen"`
	Theme     string `yaml:"theme"`
}

func Default() Config {
	return Config{
		Mouse:     true,
		AltScreen: true,
		Theme:     "default",
	}
}

func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "rfdtui", "config.yaml"), nil
}

func Load() (Config, error) {
	cfg := Default()

	p, err := Path()
	if err != nil {
		return cfg, nil
	}

	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Debug("no config file found, using defaults", "path", p)
			return cfg, nil
		}
		return cfg, fmt.Errorf("reading config: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing config: %w", err)
	}

	slog.Info("loaded config", "path", p, "mouse", cfg.Mouse, "alt_screen", cfg.AltScreen, "theme", cfg.Theme)

	return cfg, nil
}

func (c Config) Save() error {
	p, err := Path()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(p, data, 0o644)
}
