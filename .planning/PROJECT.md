# RFD TUI

## What This Is

A terminal user interface (TUI) for browsing RedFlagDeals.com forums, built with Go and Bubble Tea. Users can browse hot deals, search deals by keyword or regex, read deal threads, and sort/filter results — all from the terminal. Designed as a community tool for RFD deal hunters who prefer the command line.

## Core Value

Browse and discover RedFlagDeals hot deals from the terminal with a responsive, keyboard-driven interface.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] Display hot deals list on launch with title, score, views, replies, dealer, category, and age
- [ ] Navigate deals with vim-style keys (j/k scroll, Enter open, q back)
- [ ] View deal thread detail (original post, expand replies on demand)
- [ ] Search deals by keyword or regex across titles and dealer names
- [ ] Sort deals by score or views
- [ ] Filter deals by category and minimum score threshold
- [ ] Open deal URL in web browser with a single keypress (o)
- [ ] Paginate through multiple pages of deals

### Out of Scope

- User authentication / posting — RFD API is read-only without auth, unnecessary for browsing
- Notifications / deal alerts — defer to future version
- Bookmarks / saved deals — defer to future version
- Image rendering in terminal — out of scope for TUI
- Multiple forum categories on first screen — Hot Deals is the default; other forums can be added later

## Context

- **Reference project**: davegallant/rfd (archived Aug 2024) — Python CLI that used the RFD JSON API. Same API endpoints confirmed working as of April 2026.
- **RFD JSON API** (no auth required):
  - Base URL: `https://forums.redflagdeals.com`
  - Threads: `GET /api/topics?forum_id={id}&per_page={n}&page={n}`
  - Posts: `GET /api/topics/{id}/posts?per_page={n}&page={n}`
  - Forum IDs: 9 = Hot Deals, others TBD (Freebies, Contests, etc.)
  - Response includes: title, score, votes (up/down), views, replies, dealer name, price, savings %, deal URL, cover image, user info
  - Pagination: `pager.total_pages` field in response
- **Tech stack**: Go with Bubble Tea (Elm-style TUI framework), Lipgloss (styling), Bubbles (pre-built components)
- **Target audience**: RFD community deal hunters, CLI enthusiasts
- **Distribution**: Single static binary (go install, brew, GitHub releases)

## Constraints

- **Tech stack**: Go + Bubble Tea + Lipgloss + Bubbles — chosen for Elm-style architecture, single binary output, and community traction
- **Data source**: RFD JSON API only — no HTML scraping needed, API is stable and publicly accessible
- **No authentication**: Read-only browsing, no posting or account features
- **Terminal only**: Must work in standard terminals (no GUI dependencies)
- **Cross-platform**: Should work on Linux and macOS (Windows is nice-to-have)

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Go + Bubble Tea | Elm-style architecture, single binary, strong TUI ecosystem | — Pending |
| RFD JSON API over HTML scraping | Reference project proved the API works; cleaner, faster, more reliable than parsing HTML | — Pending |
| Vim-style navigation | Target audience is CLI users who expect vim keybindings | — Pending |
| Hot deals as first screen | Most popular forum, gets users to value immediately | — Pending |
| Single binary distribution | Easy install, no runtime dependencies, community-friendly | — Pending |
| Original post + expand replies | Avoids loading entire threads upfront; better performance for long threads | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-04-17 after initialization*
