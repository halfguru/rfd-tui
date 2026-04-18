<div align="center">

# 🏷️ rfd-tui

A beautiful terminal UI for browsing [RedFlagDeals.com](https://forums.redflagdeals.com) hot deals.

Built with Go + [Bubble Tea v2](https://charm.land)

[![Go Report Card](https://goreportcard.com/badge/github.com/simon/rfd)](https://goreportcard.com/report/github.com/simon/rfd)
[![Go Reference](https://pkg.go.dev/badge/github.com/simon/rfd.svg)](https://pkg.go.dev/github.com/simon/rfd)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev)

![demo](docs/demo.gif)

</div>

---

## ✨ Features

- **Browse Hot Deals** — scores, prices, dealer info, and savings at a glance
- **Read Threads** — full post content with vote counts
- **Search** — filter deals by text or regex
- **Sort & Filter** — by score, views, or minimum score threshold
- **Open in Browser** — jump to any deal with `o`
- **Vim-style Keybindings** — `j`/`k`, `Enter`, `/`, `q`, and more

## 📦 Install

```bash
git clone https://github.com/simon/rfd.git
cd rfd
go build -o rfd .
```

## 🚀 Usage

```bash
./rfd
```

### Keybindings

| Key        | Action                  |
| ---------- | ----------------------- |
| `j` / `k`  | Navigate up / down      |
| `Enter`    | Open thread             |
| `o`        | Open deal in browser    |
| `/`        | Search                  |
| `s`        | Cycle sort mode         |
| `f`        | Cycle min score filter  |
| `n` / `p`  | Next / previous page    |
| `Space`    | Load more posts         |
| `Esc` / `q`| Back / Quit             |
| `?`        | Help                    |

## 🛠 Tech Stack

| Library | Purpose |
| ------- | ------- |
| [Bubble Tea v2](https://charm.land/bubbletea/v2) | TUI framework |
| [Lipgloss v2](https://charm.land/lipgloss/v2) | Terminal styling |
| [Bubbles v2](https://charm.land/bubbles/v2) | UI components |

## 📄 License

[MIT](LICENSE)
