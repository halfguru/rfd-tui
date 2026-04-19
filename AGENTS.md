# AGENTS.md — Project Instructions for AI Agents

## Project: RFD TUI

A terminal UI for browsing RedFlagDeals.com hot deals, built with Go + Bubble Tea v2.

## Tech Stack

- **Language**: Go 1.26+
- **TUI Framework**: Bubble Tea v2 (`charm.land/bubbletea/v2`)
- **Styling**: Lipgloss v2 (`charm.land/lipgloss/v2`)
- **Components**: Bubbles v2 (`charm.land/bubbles/v2`)
- **Config**: YAML via `gopkg.in/yaml.v3`
- **Task Runner**: [Taskfile](https://taskfile.dev) — use `task` commands for all build/lint/test
- **CRITICAL**: These are v2 import paths. Do NOT use `github.com/charmbracelet/*` v1 paths.
  - `View()` method returns `tea.View` (not `string`)
  - Key messages are `tea.KeyPressMsg` (not `tea.KeyMsg`)
  - Program creation uses `tea.NewProgram()` (same API but different import)

## Build & Run

```bash
go build -o rfdtui .
./rfdtui
```

## Lint & Typecheck

```bash
go vet ./...
golangci-lint run
```

## Test

```bash
go test ./...
```

## RFD JSON API

- Base URL: `https://forums.redflagdeals.com`
- Threads: `GET /api/topics?forum_id=9&per_page=40&page=N` (forum_id 9 = Hot Deals)
- Posts: `GET /api/topics/{id}/posts?per_page=40&page=N`
- No auth required
- Response includes: title, score, votes (up/down), views, replies, dealer, price, savings %, deal URL, user info
- Pagination: `pager.total_pages` field in response

## Architecture

- Single root `Model` with `activeView` enum routing between views (deal list, thread detail, help)
- API client in `internal/client/` with `tea.Cmd` wrappers, retry/backoff, per-request context
- Types in `internal/types/` — Go structs matching RFD API JSON shapes
- Views in `internal/views/` — each view is a sub-model with Init/Update/View methods
- Styles in `internal/styles/` — Lipgloss style definitions
- Config in `internal/config/` — YAML config file loading with defaults

## Code Style

- No comments unless explicitly requested
- Follow existing Go conventions in the codebase
- Use `tea.Cmd` for all async operations (HTTP, browser open, clipboard)
- Keep view rendering pure — no side effects in `View()` methods
- Use built-in `min()` — do not reimplement

## Project Structure

```
main.go                    # Entry point (config loading, program options)
app.go                     # Root model (Bubble Tea Model interface)
internal/
  client/                  # RFD API client + tea.Cmd wrappers (retry, context)
  config/                  # YAML config file support
  types/                   # Go structs for API responses
  styles/                  # Lipgloss styles + theme helpers
  views/                   # View sub-models (deal_list, thread, help)
.github/
  workflows/ci.yml         # CI (build, vet, test, lint)
  workflows/release.yml    # GoReleaser on tag push
  ISSUE_TEMPLATE/          # Bug report + feature request
  PULL_REQUEST_TEMPLATE.md
  CODEOWNERS
.goreleaser.yml            # Cross-compile + homebrew release
config.example.yaml        # Example config file
```

## Key Constraints

- Read-only API access (no posting/authentication)
- Cross-platform terminal support (Linux + macOS primary)
- Single static binary output
- 15-second HTTP timeout per request, 3 retries with exponential backoff
- Descriptive User-Agent header
- Strip HTML to plain text for thread content display
- Config file at `~/.config/rfdtui/config.yaml` (optional, defaults used if missing)
