# Project Research Summary

**Project:** RFD TUI — Go terminal interface for browsing RedFlagDeals.com
**Domain:** Terminal UI (TUI) consuming a public REST JSON API
**Researched:** 2026-04-17
**Confidence:** HIGH

## Executive Summary

RFD TUI is a read-only terminal client for RedFlagDeals.com, built on the Charm v2 ecosystem (Bubble Tea, Lipgloss, Bubbles) — the de-facto standard for Go TUI applications. Expert-built TUIs in this category (hackernews-TUI, rtv, nom) all follow the same proven pattern: Elm-style Model-View-Update architecture, vim navigation, list + detail views, and async I/O via command messages. Our stack is unanimous across every source: Bubble Tea v2 for the framework, Bubbles components (list, viewport, textinput, spinner) for pre-built UI, Lipgloss for styling, and Go stdlib for HTTP/JSON. No alternatives come close.

The recommended approach is a single-root model with an `activeView` enum routing between deal list and thread detail screens. All state lives in one model struct. HTTP calls go through a `tea.Cmd` wrapper — never blocking the UI. The RFD JSON API is simple (unauthenticated GET requests returning JSON), making the API client layer thin and testable. The project should be organized into `rfd/` (API client + types), `views/` (rendering helpers), `msg/` (message types), and a root `app.go`.

The biggest risk is accidentally using Bubble Tea v1 APIs — the v2 import paths (`charm.land/*`), return types (`tea.View` not `string`), and event types (`tea.KeyPressMsg` not `tea.KeyMsg`) are fundamentally different. The second risk is RFD API stability: the reference Python project was archived in 2024, and while the API works as of April 2026, there are no guarantees. Mitigate by building a resilient API client that handles unexpected responses gracefully.

## Key Findings

### Recommended Stack

The Charm v2 ecosystem is the clear winner for Go TUI development. Bubble Tea v2 (v2.0.6) provides the Elm-style architecture, Lipgloss v2 (v2.0.3) handles declarative styling and layout, and Bubbles v2 (v2.1.0) supplies battle-tested components. All HTTP/JSON needs are served by Go stdlib — the RFD API is simple enough that external libraries add no value. CLI flags use stdlib `flag`; browser opening uses `os/exec`. No Cobra, no RESTy, no JSON decoders beyond stdlib.

**Core technologies:**
- **Bubble Tea v2** (`charm.land/bubbletea/v2`): TUI framework — Elm architecture, 41.6k stars, v2 is stable with declarative views
- **Lipgloss v2** (`charm.land/lipgloss/v2`): Terminal styling — CSS-like declarative styles, responsive layout, borders, tables
- **Bubbles v2** (`charm.land/bubbles/v2`): Pre-built components — list (fuzzy filter, pagination), viewport (scrollable content), textinput, spinner, help
- **Bubblezone v2** (`github.com/lrstanley/bubblezone/v2`): Mouse zone tracking — clickable regions for mouse support
- **Go stdlib** (`net/http`, `encoding/json`, `flag`, `os/exec`): HTTP client, JSON parsing, CLI flags, browser opening — no external deps needed

**Critical version note:** All Charm libraries use `charm.land/*` vanity import paths in v2, NOT `github.com/charmbracelet/*`. Bubblezone is NOT a Charm library — it stays on GitHub. View() returns `tea.View` (not `string`), keys use `tea.KeyPressMsg` (not `tea.KeyMsg`).

### Expected Features

Feature analysis is based on four comparable TUIs (hackernews-TUI, rtv, nom, davegallant/rfd CLI) plus the PROJECT.md requirements. The feature set is well-understood with clear MVP boundaries.

**Must have (table stakes — P1):**
- Deal list view — core value prop, users launch to see deals
- Vim-style navigation (j/k/enter/q) — expected by every CLI user; also support arrow keys
- Thread/post detail view — scrollable viewport with original post + lazy-loaded replies
- Open deal in browser (`o` key) — TUI can't replace browser for images/purchase
- Search/filter — keyword search across titles and dealer names
- Sort by score/views — client-side sort, no API parameter needed
- Pagination — `n`/`p` keys, RFD API has `page` param and `pager.total_pages`
- Help overlay (`?` key) — auto-generated from key bindings
- Loading spinner — visual feedback during API calls (1-3 second fetches)
- Clean quit (`q`/`Ctrl+C`) — Bubble Tea handles this natively

**Should have (differentiators — P2):**
- Color-coded scores (green=hot, yellow=warm, red=cold) — instant visual scanning
- Category filtering — impossible in web UI's linear list
- Minimum score threshold — hide low-quality deals automatically
- Regex search — power users expect this (reference CLI supported it)
- Multiple forum sections (F1=Hot Deals, F2=Freebies, F3=Contests) — same API endpoint, different `forum_id`
- Collapsible comment threads — essential for 100+ reply threads
- Deal age display — "2h ago" relative timestamps
- JSON output mode (`--json` flag) — pipe to `jq`, scripting/automation

**Defer (v2+):**
- Configurable keybindings (TOML config) — requires config parsing, validation, migration
- Color themes — cosmetic, ship one good default first
- Shell completions — only useful for CLI mode, not TUI
- Markdown rendering (Glamour) — RFD posts are HTML, not markdown; would need HTML→markdown first

**Anti-features (do NOT build):**
- Real-time auto-refresh — rate limiting risk, breaks "open, scan, close" workflow
- User authentication/posting — no official API, ToS risk, massive scope creep
- Image rendering — fragmented terminal protocols, most users' terminals won't work
- Notifications/deal alerts — separate product category (rfd-notify exists)
- Bookmarks — scope creep into "deal management"
- Offline mode/caching — stale data is low-value for hourly-changing deals

### Architecture Approach

The application follows Bubble Tea's Elm architecture with a single root model routing between views via an `activeView` enum. This pattern is verified across 7+ production Bubble Tea apps. The root model holds all state (navigation, data, UI components, window dimensions) and delegates Update/View calls based on which screen is active. Sub-components (list, viewport, textinput) are embedded fields, not separate `tea.Model` implementations — avoiding the wiring complexity of component trees.

**Major components:**
1. **Root Model (`app.go`)** — Application state, view routing, key dispatch. Holds `activeView` enum, sub-models, API client, window dimensions.
2. **API Client (`rfd/client.go`)** — HTTP requests to RFD JSON API with timeout, User-Agent header, error wrapping. `FetchTopics()`, `FetchPosts()`, `SearchTopics()`.
3. **API Types (`rfd/types.go`)** — Go structs matching RFD JSON response format (Topic, Post, User, Pager).
4. **Views (`views/deals.go`, `views/thread.go`)** — Rendering helpers for deal list (custom list delegate) and thread detail (post formatting). Pure functions accepting data, returning strings.
5. **Message Types (`msg/messages.go`)** — Custom `tea.Msg` types for async results (`dealsLoadedMsg`, `postsLoadedMsg`).
6. **Styles (`styles.go`)** — Lipgloss style constants for consistent visual presentation.

**Key patterns:**
- **Active View Routing:** `activeView` enum (dealList, threadDetail, search) determines Update/View behavior
- **Command-Based Async I/O:** All HTTP calls wrapped in `tea.Cmd` closures returning `tea.Msg` — never block Update()
- **Embedded Sub-Models:** Bubbles components (list, viewport, textinput) as model fields, forwarded messages
- **Build order:** types.go → client.go → messages.go → keys.go → styles.go → views/ → app.go → main.go

### Critical Pitfalls

1. **Using Bubble Tea v1 APIs instead of v2** — Import paths changed to `charm.land/*`, View() returns `tea.View` not `string`, keys use `tea.KeyPressMsg` not `tea.KeyMsg`. Prevention: establish correct imports in Phase 1 scaffolding, reference official upgrade guide.

2. **Blocking in Update()** — HTTP calls in Update() freeze the entire TUI. Prevention: wrap ALL I/O in `tea.Cmd` functions from the very first API call. No exceptions.

3. **Ignoring WindowSizeMsg** — Terminal resize breaks layout with stale dimensions. Prevention: handle `WindowSizeMsg` from day one, store width/height on model, propagate to sub-components.

4. **Monolithic model growing unmanageable** — Single struct accumulating all state for all views. Prevention: use ViewID enum + sub-model decomposition from project start. Retroactive decomposition is painful.

5. **No HTTP client timeout** — `http.DefaultClient` has no timeout; unresponsive API hangs TUI forever. Prevention: create shared `http.Client{Timeout: 15s}` at app init.

6. **Raw HTML/HTML entities in terminal** — RFD posts contain HTML, special characters, emoji. Prevention: strip HTML tags, decode entities (`html.UnescapeString`), use `lipgloss.Width()` for measurements, test with real API data from day one.

7. **Loading entire thread on open** — Popular threads have 1000+ replies. Prevention: load first page only (original post), lazy-load replies on demand with API pagination.

## Implications for Roadmap

Based on research, suggested phase structure:

### Phase 1: Foundation & Deal Browsing

**Rationale:** The API client is the foundation for everything. Building it first with proper HTTP client configuration, type definitions, and Bubble Tea v2 scaffolding prevents compounding mistakes. The deal list view is the core value proposition — users should see deals within the first phase.

**Delivers:** Working TUI that fetches and displays hot deals from RFD API with vim navigation, loading states, error handling, and clean terminal resize behavior.

**Addresses:** RFD API Client, Deal List View, Vim Navigation, Sort by Score/Views, Pagination, Help Overlay, Loading Spinner, Clean Quit (all P1 features except Thread Detail and Search)

**Avoids:** v1/v2 API confusion (correct imports from start), blocking in Update (tea.Cmd pattern established), WindowSizeMsg (handled in deal list), monolithic model (ViewID + sub-models from day one), HTTP timeout (configured in API client), content sanitization (real API data from start)

**Build order within phase:**
1. `rfd/types.go` — API response structs
2. `rfd/client.go` — HTTP client with timeout, User-Agent, error handling
3. `msg/messages.go` — tea.Msg types
4. `keys.go` — Key binding definitions
5. `styles.go` — Lipgloss style constants
6. `views/deals.go` — Deal list delegate + rendering
7. `app.go` — Root model with ViewID routing
8. `main.go` — Entry point, wiring, program start

### Phase 2: Thread Detail & Browser Integration

**Rationale:** Thread reading is the #2 use case after deal browsing. It requires a new view (viewport) and a new API endpoint (posts). Opening deals in the browser is the bridge between TUI and the full web experience.

**Delivers:** Thread detail view with original post display, lazy-loaded paginated replies, and browser URL opening.

**Addresses:** Thread Detail View, Open in Browser, Lazy Reply Loading

**Avoids:** Loading entire thread on open (lazy pagination from the start), no URL error handling (graceful behavior for discussion-only deals)

**Uses:** Bubbles viewport component, `tea.ExecProcess` for browser, RFD posts API endpoint

### Phase 3: Search, Filter & Sort Enhancements

**Rationale:** Search and filtering are enhancements to the deal list — they operate on the same view with the same data. Building them after the core list and thread views means the infrastructure is stable and these are additive features.

**Delivers:** Keyword search, regex search, category filtering, minimum score threshold, deal age display, color-coded scores.

**Addresses:** Search/Filter, Regex Search, Category Filtering, Score Threshold, Deal Age Display, Color-coded Scores

**Avoids:** Search without debouncing (trigger on Enter, not every keystroke), rate limiting (no API calls on every keypress)

### Phase 4: Multi-Forum & Polish

**Rationale:** Multiple forum sections (Freebies, Contests) reuse the entire deal list infrastructure with a different `forum_id` parameter — very low marginal cost. Polish items (JSON output mode, better error messages) round out the experience.

**Delivers:** F-key forum switching (F1=Hot Deals, F2=Freebies, F3=Contests), JSON output mode (`--json` flag), collapsible comment threads.

**Addresses:** Multiple Forum Sections, JSON Output Mode, Collapsible Comments

### Phase 5: Distribution & Configuration (Post-MVP)

**Rationale:** Configurable keybindings, color themes, and shell completions are customization features that only matter after users validate the core product. GoReleaser setup enables distribution via brew/go install/GitHub Releases.

**Delivers:** TOML config file for keybindings, multiple color themes, shell completions, cross-platform binary releases.

**Addresses:** Configurable Keybindings, Color Themes, Shell Completions, GoReleaser setup

### Phase Ordering Rationale

- **Phase 1 first** because the API client + deal list is the foundation for every other feature — all features depend on fetching and displaying deals
- **Phase 2 second** because thread detail requires a new API endpoint and new view, but builds on the established deal list navigation
- **Phase 3 third** because search/filter/sort are purely additive to the deal list — same view, same data, client-side operations
- **Phase 4 fourth** because multi-forum is a trivial extension (different `forum_id`), collapsible comments enhance the thread view, and JSON output is a separate mode
- **Phase 5 last** because customization and distribution only matter after core product validation
- Architecture build order (types → client → messages → keys → styles → views → app → main) matches dependency graph — each layer depends only on previous layers

### Research Flags

Phases likely needing deeper research during planning:
- **Phase 2:** Thread detail requires understanding RFD posts API response format in detail — need to verify `parent_id` nesting structure for future collapsible comments
- **Phase 4:** Multiple forum sections need enumeration of all RFD forum IDs (currently only Hot Deals=9 is confirmed; Freebies=15 and Contests=12 need verification against live API)

Phases with standard patterns (skip research-phase):
- **Phase 1:** Well-documented Bubble Tea v2 patterns — Elm architecture, tea.Cmd async, list component, activeView routing all verified across 7+ production apps and official docs
- **Phase 3:** Standard client-side operations — Go regexp, sort.Slice, lipgloss styling — no external APIs or complex integrations
- **Phase 5:** Standard GoReleaser config and TOML parsing — well-documented tooling

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH | All versions verified via Go module proxy. Bubble Tea v2 stable (v2.0.6, released 2026-04-16). Ecosystem well-established with 41.6k stars. Import paths and API changes confirmed via official upgrade guide. |
| Features | HIGH | Feature set derived from 4 comparable TUIs and the reference Python CLI. MVP boundaries clear. Prioritization matrix validated against real-world usage patterns. |
| Architecture | HIGH | Active View Routing pattern verified across 7+ production Bubble Tea apps (yt-browse, gocovsh, wasteland, hey-cli, Grafana Loki bench TUI). Build order follows dependency graph. Anti-patterns documented with real-world examples. |
| Pitfalls | HIGH | All critical pitfalls verified against Bubble Tea v2 source, official docs, and production code. Recovery strategies are mechanical (not architectural). Pitfall-to-phase mapping ensures prevention in the right phase. |

**Overall confidence:** HIGH

### Gaps to Address

- **RFD API response format for posts:** The posts endpoint (`GET /api/topics/{id}/posts`) response structure needs validation against live data. The reference Python project was archived in 2024. Verify during Phase 2 planning with a live API call.
- **RFD forum ID enumeration:** Only Hot Deals (forum_id=9) is confirmed. Freebies (15) and Contests (12) are from the reference project and need verification. Verify during Phase 4 planning.
- **RFD API rate limiting:** No official rate limit documentation. The API is public and appears unthrottled, but aggressive polling could trigger blocks. Conservative approach: manual refresh only, no auto-refresh. Monitor during development.
- **RFD HTML content complexity:** We know posts contain HTML, but the full range of HTML elements (embedded videos, tables, complex formatting) is unknown. May need adjustment to the HTML stripping approach when real thread data is rendered in Phase 2.

## Sources

### Primary (HIGH confidence)
- Bubble Tea v2 upgrade guide: `github.com/charmbracelet/bubbletea/blob/main/UPGRADE_GUIDE_V2.md` — v1→v2 API changes, import paths, View/Key types
- Context7: `/charmbracelet/bubbletea` — Elm architecture, tea.Cmd patterns, multiple views
- Context7: `/charmbracelet/lipgloss` — Layout, styling, measurement utilities
- Context7: `/charmbracelet/bubbles` — list, viewport, textinput, spinner, paginator, help components
- Go module proxy: Verified versions for all Charm v2 libraries (bubbletea v2.0.6, lipgloss v2.0.3, bubbles v2.1.0, bubblezone v2.0.0)

### Secondary (MEDIUM confidence)
- **aome510/hackernews-TUI** (696 stars) — Feature patterns for TUI forum clients, vim navigation, F-key sections, configurable keys
- **davegallant/rfd** (archived Aug 2024) — Reference Python CLI proving RFD JSON API works, feature validation
- **yt-browse, gocovsh, wasteland** — Production Bubble Tea apps confirming activeView routing pattern
- **Grafana Loki bench TUI** — ViewID + sub-model decomposition pattern

### Tertiary (LOW confidence)
- RFD JSON API stability — API works as of April 2026 but no official documentation or stability guarantees; reference project archived
- RFD forum IDs beyond Hot Deals (9) — values 15 (Freebies) and 12 (Contests) from reference project, not verified against current API
- RFD rate limiting behavior — no documentation; conservative approach adopted

---
*Research completed: 2026-04-17*
*Ready for roadmap: yes*
