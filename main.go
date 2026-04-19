package main

import (
	"log/slog"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/halfguru/rfd-tui/internal/config"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	_       = version
	_       = commit
	_       = date
)

func initLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, opts))
	slog.SetDefault(logger)
}

func main() {
	initLogger()

	cfg, err := config.Load()
	if err != nil {
		slog.Warn("config load failed, using defaults", "error", err)
	}

	slog.Info("starting rfd-tui", "version", version, "commit", commit)

	m := NewModel(cfg)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
}
