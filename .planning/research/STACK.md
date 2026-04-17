# Technology Stack

**Project:** RFD TUI — Go terminal interface for browsing RedFlagDeals.com
**Researched:** 2026-04-17
**Go version:** 1.25+ (required by Charm v2 ecosystem)

## Recommended Stack

### Core TUI Framework

| Technology | Version | Purpose | Why | Confidence |
|------------|---------|---------|-----|------------|
| Bubble Tea | v2.0.6 | TUI framework (Elm architecture) | De-facto standard Go TUI framework. v2 is stable with declarative views, better mouse support, cleaner API. 41.6k GitHub stars. | HIGH |
| Lipgloss | v2.0.3 | Terminal styling & layout | Declarative CSS-like styling. Same Charm ecosystem as Bubble Tea. Tables, borders, colors, responsive layout built in. | HIGH |
| Bubbles | v2.1.0 | Pre-built TUI components | Battle-tested components (List, Viewport, Spinner, Paginator, Help, Key bindings). Maintained by Charm, used in production by Crush. | HIGH |
| Bubblezone | v2.0.0 | Mouse event zone tracking | Zero-width zone markers for clickable regions. Essential if we want mouse support in the deal list. | HIGH |

### Network & Data

| Technology | Version | Purpose | Why | Confidence |
|------------|---------|---------|-----|------------|
| `net/http` | stdlib | HTTP client for RFD API | RFD API is simple GET requests returning JSON. No auth, no retries, no complexity. Standard library is the right tool. | HIGH |
| `encoding/json` | stdlib | JSON response parsing | Standard library handles RFD's JSON responses fine. No performance bottleneck at TUI scale. | HIGH |

### Supporting Libraries

| Library | Version | Purpose | When to Use | Confidence |
|---------|---------|---------|-------------|------------|
| `os/exec` | stdlib | Open URLs in browser | When user presses `o` on a deal. Call `open` (macOS) or `xdg-open` (Linux). No external dependency needed. | HIGH |
| `flag` | stdlib | CLI argument parsing | For `--forum`, `--search`, `--version` flags. Simple enough that Cobra is overkill. | HIGH |
| `fmt` | stdlib | Text formatting | Standard Go formatting for rendering deal info. | HIGH |

### Optional / Future

| Library | Version | Purpose | When to Use | Confidence |
|---------|---------|---------|-------------|------------|
| Glamour | v2.0.0 | Markdown rendering in terminal | If we convert RFD HTML posts to markdown for richer thread display. Defer to post-MVP — plain text rendering is sufficient for launch. | MEDIUM |
| GoReleaser | latest | Cross-platform binary releases | When ready to distribute via `brew`, `go install`, GitHub Releases. Set up in a later phase. | HIGH |

## Charm v2 Ecosystem Notes

### Import Paths Changed (Critical)

All Charm libraries migrated from `github.com/charmbracelet/*` to `charm.land/*` vanity domains:

```go
// v1 (OLD — do not use)
import tea "github.com/charmbracelet/bubbletea"
import "github.com/charmbracelet/lipgloss"
import "github.com/charmbracelet/bubbles/list"

// v2 (CORRECT — use these)
import tea "charm.land/bubbletea/v2"
import "charm.land/lipgloss/v2"
import "charm.land/bubbles/v2/list"
```

Bubblezone is NOT a Charm library — it stays on GitHub:
```go
import zone "github.com/lrstanley/bubblezone/v2"
```

### Bubble Tea v2 Key API Changes

| Aspect | v1 | v2 |
|--------|----|----|
| View method | `View() string` | `View() tea.View` |
| Alt screen | `tea.WithAltScreen()` program option | `view.AltScreen = true` in View() |
| Mouse mode | `tea.WithMouseCellMotion()` program option | `view.MouseMode = tea.MouseModeCellMotion` in View() |
| Key events | `tea.KeyMsg` (struct) | `tea.KeyPressMsg` for key presses |
| Mouse events | `tea.MouseMsg` with `.Action` check | `tea.MouseClickMsg`, `tea.MouseWheelMsg`, etc. |
| Space key | `case " ":` | `case "space":` |
| Window title | `tea.SetWindowTitle("...")` command | `view.WindowTitle = "..."` in View() |

### Bubbles v2 Components to Use

| Component | Use In RFD TUI | Purpose |
|-----------|---------------|---------|
| **list** | Deal list screen | Fuzzy filtering, pagination, status bar, spinner. Primary browsing component. |
| **viewport** | Thread detail screen | Scrollable content for deal posts and replies. Supports mouse wheel. |
| **spinner** | Loading states | Visual feedback while fetching from RFD API. |
| **paginator** | Deal pagination | Page through multiple pages of deals from API (`pager.total_pages`). |
| **help** | Keybinding help | Auto-generated help view from key bindings. Toggle with `?`. |
| **key** | Key binding management | Define vim-style keys (`j/k/enter/q/o`) with help text. Enables key remapping. |
| **textinput** | Search bar | Keyword/regex search input for filtering deals. |

### Bubbles Components NOT Needed

| Component | Why Not |
|-----------|---------|
| table | List is better for card-style deal items. Table is for strict tabular data. |
| textarea | No user text input in RFD (read-only browsing). |
| filepicker | No file system interaction needed. |
| timer / stopwatch | No timing features. |
| progress | No progress tracking. |

## Architecture Pattern

```
┌─────────────────────────────────────────────┐
│                  main.go                     │
│              (flag parsing,                  │
│               tea.Program init)              │
├─────────────────────────────────────────────┤
│                                             │
│  ┌───────────┐  ┌────────────────────────┐  │
│  │  Model     │  │  API Client            │  │
│  │  (state)   │  │  (net/http + json)     │  │
│  │            │  │                        │  │
│  │  - deals   │◄─┤  - FetchDeals()       │  │
│  │  - thread  │  │  - FetchThread()      │  │
│  │  - cursor  │  │  - SearchDeals()      │  │
│  │  - view    │  │                        │  │
│  └─────┬──────┘  └────────────────────────┘  │
│        │                                     │
│  ┌─────▼──────────────────────────────────┐  │
│  │  Update (tea.Msg routing)              │  │
│  │  - tea.KeyPressMsg → navigation        │  │
│  │  - tea.WindowSizeMsg → responsive      │  │
│  │  - custom msgs → API results           │  │
│  └─────┬──────────────────────────────────┘  │
│        │                                     │
│  ┌─────▼──────────────────────────────────┐  │
│  │  View (render with Lipgloss)           │  │
│  │  - deals list (bubbles/list)           │  │
│  │  - thread view (bubbles/viewport)      │  │
│  │  - search bar (bubbles/textinput)      │  │
│  │  - help overlay (bubbles/help)         │  │
│  └────────────────────────────────────────┘  │
│                                             │
└─────────────────────────────────────────────┘
```

## Alternatives Considered

| Category | Recommended | Alternative | Why Not |
|----------|-------------|-------------|---------|
| TUI Framework | Bubble Tea v2 | tview | tview uses widget-based approach, less composable. Bubble Tea's Elm architecture scales better for complex state. Bubble Tea has much larger community. |
| TUI Framework | Bubble Tea v2 | termui | Stale development (last significant update 2020). Uses buffer-based rendering. Not maintained. |
| TUI Framework | Bubble Tea v2 | tcell directly | Too low-level. Bubble Tea handles rendering, focus, layout. No reason to build from scratch. |
| HTTP Client | net/http | RESTy, req | Overkill for simple GET requests to a public API. More dependencies, no benefit. |
| JSON | encoding/json | jsoniter, sonic | Performance not a concern at TUI scale. Dealing with ~40 deals per page. Stdlib is fine. |
| CLI | flag | Cobra, pflag | App is a TUI launched with optional flags. No subcommands, no complex CLI structure. Cobra adds complexity for no gain. |
| Styling | Lipgloss v2 | tcell styles, manual ANSI | Lipgloss gives declarative styling, layout (JoinHorizontal/Vertical), borders, tables. Much more productive. |
| Components | Bubbles list | bubble-table (evertras) | bubble-table is good for strict tabular data, but RFD deals are better as styled cards in a list. Bubbles list has fuzzy filtering built in. |
| Mouse | Bubblezone | Raw tea.MouseMsg | Bubblezone makes click regions trivial with zone markers. Without it, you'd manually track element positions and compare mouse coordinates. |
| Markdown | Glamour | Skip | Defer to post-MVP. RFD posts are HTML, not markdown. Would need HTML→markdown conversion first. Plain text rendering is sufficient. |

## Installation

```bash
# Initialize Go module (requires Go 1.25+)
go mod init github.com/user/rfd

# Core TUI framework
go get charm.land/bubbletea/v2@latest
go get charm.land/lipgloss/v2@latest
go get charm.land/bubbles/v2@latest

# Mouse support
go get github.com/lrstanley/bubblezone/v2@latest

# Optional: markdown rendering (defer to post-MVP)
# go get charm.land/glamour/v2@latest
```

No dev dependencies beyond standard Go tooling (`go test`, `go vet`).

## Key Technical Decisions

### 1. Bubble Tea v2 over v1
**Decision:** Use v2 exclusively.
**Rationale:** v2 is stable (v2.0.6, released 2026-04-16). Declarative views eliminate an entire class of state bugs (alt screen, mouse mode live in View, not scattered across program options). The API is cleaner. No reason to start a new project on v1.

### 2. Bubbles List over custom rendering
**Decision:** Use Bubbles' `list` component for the deal browser.
**Rationale:** It provides fuzzy filtering, pagination, spinner, status messages, and help generation out of the box. Custom rendering via `itemDelegate` gives full control over deal card appearance. Building this from scratch would take days; with Bubbles it's configuration.

### 3. Viewport for thread reading
**Decision:** Use Bubbles' `viewport` for scrollable thread content.
**Rationale:** Handles mouse wheel, keyboard scrolling, soft wrapping, and gutter/line numbers. Essential for reading long deal threads.

### 4. Standard library for HTTP/JSON
**Decision:** Use `net/http` + `encoding/json`.
**Rationale:** RFD API is simple unauthenticated GET requests with JSON responses. Adding a REST client library would be a dependency for no benefit.

### 5. `os/exec` for browser opening
**Decision:** Use `os/exec` to call `open`/`xdg-open` directly.
**Rationale:** Two lines of platform-specific code vs. an external dependency. The `pkg/browser` library hasn't been updated since 2024 and adds no value over the straightforward approach.

### 6. `flag` over CLI frameworks
**Decision:** Use standard library `flag` package.
**Rationale:** The app is a TUI. CLI flags are just `--forum`, `--search`, `--version`. No subcommands, no complex parsing. Cobra is designed for CLI tools with many subcommands, not TUI apps.

## Sources

- Bubble Tea v2 upgrade guide: https://github.com/charmbracelet/bubbletea/blob/main/UPGRADE_GUIDE_V2.md (HIGH confidence — official source)
- Bubble Tea latest version: Go module proxy `charm.land/bubbletea/v2` → v2.0.6 (2026-04-16) (HIGH)
- Lipgloss latest version: Go module proxy `charm.land/lipgloss/v2` → v2.0.3 (2026-04-13) (HIGH)
- Bubbles latest version: Go module proxy `charm.land/bubbles/v2` → v2.1.0 (2026-03-25) (HIGH)
- Bubblezone latest version: Go module proxy → v2.0.0 (2026-02-28) (HIGH)
- Glamour latest version: Go module proxy `charm.land/glamour/v2` → v2.0.0 (2026-03-09) (HIGH)
- Bubbles README (component list): https://github.com/charmbracelet/bubbles/blob/main/README.md (HIGH)
- Charm ecosystem overview: https://charm.land/ (HIGH)
- Context7 Bubble Tea docs: /charmbracelet/bubbletea (HIGH)
- Context7 Lipgloss docs: /charmbracelet/lipgloss (HIGH)
- Context7 Bubbles docs: /charmbracelet/bubbles (HIGH)
- Context7 Bubblezone docs: /lrstanley/bubblezone (HIGH)
