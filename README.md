<div align="center">

# 🏷️ rfd-tui

A beautiful terminal UI for browsing [RedFlagDeals.com](https://forums.redflagdeals.com) hot deals.

Built with Go + [Bubble Tea v2](https://charm.land)

[![CI](https://github.com/halfguru/rfd-tui/actions/workflows/ci.yml/badge.svg)](https://github.com/halfguru/rfd-tui/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/halfguru/rfd-tui)](https://goreportcard.com/report/github.com/halfguru/rfd-tui)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://go.dev)

![demo](docs/demo.gif)

</div>

---

## ✨ Features

- **Browse Hot Deals** — scores, prices, dealer info, and savings at a glance
- **Hot Badge** — trending deals (score ≥ 25) highlighted with a `HOT` badge
- **Read Threads** — full post content with vote counts and user avatars
- **Search** — filter deals by text or regex
- **Sort & Filter** — by score, views, or minimum score threshold
- **Open in Browser** — jump to any deal with `o` (works in deal list and thread views)
- **Copy to Clipboard** — copy deal URL with `c`
- **Mouse Support** — click to navigate, scroll to move
- **Config File** — customize behavior via `~/.config/rfd-tui/config.yaml`
- **Scroll Indicator** — visual position marker in deal list
- **Loading Shimmer** — skeleton placeholders while loading
- **Alt Screen** — fullscreen terminal mode (configurable)
- **Vim-style Keybindings** — `j`/`k`, `Enter`, `/`, `q`, and more

## 📦 Install

### From Source

```bash
git clone https://github.com/halfguru/rfd-tui.git
cd rfd-tui
task build
```

### From Release

Download the latest binary from [Releases](https://github.com/halfguru/rfd-tui/releases).

## 🚀 Usage

```bash
task run
# or
./rfdtui
```

### Keybindings

| Key          | Action                  |
| ------------ | ----------------------- |
| `j` / `k`    | Navigate up / down      |
| `Enter`      | Open thread             |
| `o`          | Open deal in browser    |
| `c`          | Copy deal URL           |
| `/`          | Search                  |
| `s`          | Cycle sort mode         |
| `f`          | Cycle min score filter  |
| `n` / `p`    | Next / previous page    |
| `Space`      | Load more posts         |
| `Esc` / `q`  | Back / Quit             |
| `?`          | Help                    |

### Configuration

Create `~/.config/rfdtui/config.yaml` (see `config.example.yaml`):

```yaml
mouse: true
alt_screen: true
theme: default
```

## 🛠 Tech Stack

| Library | Purpose |
| ------- | ------- |
| [Bubble Tea v2](https://charm.land/bubbletea/v2) | TUI framework |
| [Lipgloss v2](https://charm.land/lipgloss/v2) | Terminal styling |
| [Bubbles v2](https://charm.land/bubbles/v2) | UI components |

## 📄 License

[MIT](LICENSE)
