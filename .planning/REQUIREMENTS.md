# Requirements: RFD TUI

**Defined:** 2026-04-17
**Core Value:** Browse and discover RedFlagDeals hot deals from the terminal with a responsive, keyboard-driven interface.

## v1 Requirements

### API Client

- [ ] **API-01**: App fetches hot deal topics from `GET /api/topics?forum_id=9&per_page=40&page=N` and parses JSON response into typed Go structs
- [ ] **API-02**: App fetches thread posts from `GET /api/topics/{id}/posts?per_page=40&page=N` with pagination support
- [ ] **API-03**: HTTP client has a 15-second timeout and a descriptive User-Agent header
- [ ] **API-04**: API errors (network failure, non-200 status, malformed JSON) are caught and displayed as user-friendly error messages in the TUI

### Deal List View

- [ ] **LIST-01**: App launches directly into a scrollable deal list showing title, score (color-coded green/yellow/red), views, replies, dealer name, and relative age ("2h ago") per deal
- [ ] **LIST-02**: Deal list supports vim-style navigation (j/k to scroll, enter to open, q to quit) and arrow keys
- [ ] **LIST-03**: Deal list displays a loading spinner while fetching deals from the API
- [ ] **LIST-04**: Deal list handles terminal resize gracefully, updating layout to match new dimensions
- [ ] **LIST-05**: Users can paginate through deal pages with n/p keys, fetching the next/previous page from the API

### Thread Detail

- [ ] **THRD-01**: User can open a deal to view the thread detail showing the original post with body text (HTML stripped to plain text)
- [ ] **THRD-02**: Thread detail lazy-loads replies on demand, fetching additional pages from the API as the user scrolls down
- [ ] **THRD-03**: User can collapse and expand comment threads with tab key for easier navigation of long threads
- [ ] **THRD-04**: User can return to the deal list from thread detail with escape/q key

### Browser Integration

- [ ] **BROW-01**: User can press `o` to open the selected deal's URL in the system web browser (`open` on macOS, `xdg-open` on Linux)
- [ ] **BROW-02**: Browser opening handles missing URL gracefully (discussion-only deals with no external link)

### Search & Filter

- [ ] **SRCH-01**: User can search deals by keyword across titles and dealer names via an inline search input
- [ ] **SRCH-02**: User can search using regular expressions (regex errors fall back to plain text search)
- [ ] **SRCH-03**: User can filter the deal list by category (e.g., Electronics, Home, Food)
- [ ] **SRCH-04**: User can set a minimum score threshold to hide low-quality deals
- [ ] **SRCH-05**: User can sort deals by score or view count (client-side sort toggle)

### UX Polish

- [ ] **UX-01**: User can press `?` to show a help overlay listing all available keybindings for the current view
- [ ] **UX-02**: App exits cleanly on `q` or `Ctrl+C`, restoring terminal state

## v2 Requirements

### Distribution

- **DIST-01**: Single static binary installable via `go install`, brew tap, or GitHub Releases
- **DIST-02**: GoReleaser CI pipeline for automated cross-platform binary builds

### Customization

- **CUST-01**: Configurable keybindings via TOML config file at `~/.config/rfd-tui/config.toml`
- **CUST-02**: Multiple color themes (dark, light, minimal) selectable via config
- **CUST-03**: Shell completions for bash and zsh

### Multi-Forum

- **FORUM-01**: Quick-switch between forum sections with F-keys (F1=Hot Deals, F2=Freebies, F3=Contests)

### Output

- **OUT-01**: JSON output mode via `--json` flag for scripting and piping to jq

## Out of Scope

| Feature | Reason |
|---------|--------|
| User authentication / posting | No official API, ToS risk, read-only browsing is the product |
| Real-time auto-refresh | Rate limiting risk, breaks scan-and-close workflow |
| Image rendering in terminal | Fragmented terminal protocols, most users won't benefit |
| Notifications / deal alerts | Separate product category (rfd-notify exists) |
| Bookmarks / saved deals | Scope creep into deal management |
| Offline mode / caching | Deals change hourly, stale data is low-value |
| HTML content rendering | Strip to plain text is sufficient for deal browsing |
| Windows support | Nice-to-have, not blocking; Linux + macOS are primary targets |

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| API-01 | Phase 1 | Pending |
| API-02 | Phase 2 | Pending |
| API-03 | Phase 1 | Pending |
| API-04 | Phase 1 | Pending |
| LIST-01 | Phase 1 | Pending |
| LIST-02 | Phase 1 | Pending |
| LIST-03 | Phase 1 | Pending |
| LIST-04 | Phase 1 | Pending |
| LIST-05 | Phase 1 | Pending |
| THRD-01 | Phase 2 | Pending |
| THRD-02 | Phase 2 | Pending |
| THRD-03 | Phase 2 | Pending |
| THRD-04 | Phase 2 | Pending |
| BROW-01 | Phase 2 | Pending |
| BROW-02 | Phase 2 | Pending |
| SRCH-01 | Phase 3 | Pending |
| SRCH-02 | Phase 3 | Pending |
| SRCH-03 | Phase 3 | Pending |
| SRCH-04 | Phase 3 | Pending |
| SRCH-05 | Phase 3 | Pending |
| UX-01 | Phase 3 | Pending |
| UX-02 | Phase 1 | Pending |

**Coverage:**
- v1 requirements: 22 total
- Mapped to phases: 22
- Unmapped: 0 ✓

---
*Requirements defined: 2026-04-17*
*Last updated: 2026-04-17 after roadmap creation*
