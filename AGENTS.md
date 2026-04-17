# AGENTS.md — Project Instructions for AI Agents

## Project: RFD TUI

A terminal UI for browsing RedFlagDeals.com hot deals, built with Go + Bubble Tea v2.

## Tech Stack

- **Language**: Go 1.22+
- **TUI Framework**: Bubble Tea v2 (`charm.land/bubbletea/v2`)
- **Styling**: Lipgloss v2 (`charm.land/lipgloss/v2`)
- **Components**: Bubbles v2 (`charm.land/bubbles/v2`)
- **CRITICAL**: These are v2 import paths. Do NOT use `github.com/charmbracelet/*` v1 paths.
  - `View()` method returns `tea.View` (not `string`)
  - Key messages are `tea.KeyPressMsg` (not `tea.KeyMsg`)
  - Program creation uses `tea.NewProgram()` (same API but different import)

## Build & Run

```bash
go build -o rfd ./cmd/rfd
./rfd
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

- Single root `Model` with `activeView` enum routing between views (deal list, thread detail, search, help)
- API client in `internal/client/` with `tea.Cmd` wrappers for all HTTP calls
- Types in `internal/types/` — Go structs matching RFD API JSON shapes
- Views in `internal/views/` — each view is a sub-model with Init/Update/View methods
- Key bindings in `internal/keys/` — vim-style (j/k, Enter, q, n/p, o, ?, Escape)

## Code Style

- No comments unless explicitly requested
- Follow existing Go conventions in the codebase
- Use `tea.Cmd` for all async operations (HTTP, browser open)
- Keep view rendering pure — no side effects in `View()` methods

## Project Structure

```
cmd/rfd/main.go          # Entry point
internal/
  client/                 # RFD API client + tea.Cmd wrappers
  types/                  # Go structs for API responses
  keys/                   # Key bindings
  styles/                 # Lipgloss styles
  views/                  # View sub-models (deal_list, thread, search, help)
app.go                    # Root model (Bubble Tea Model interface)
```

## Workflow

- GSD workflow is configured in `.planning/config.json`
- Mode: YOLO (auto-approve), coarse granularity, parallel execution
- Phase plans live in `.planning/phases/`
- All docs are committed to git
- Current phase: Phase 1 — Foundation & Deal List

## Key Constraints

- Read-only API access (no posting/authentication)
- Cross-platform terminal support (Linux + macOS primary)
- Single static binary output
- 15-second HTTP timeout, descriptive User-Agent header
- Strip HTML to plain text for thread content display
