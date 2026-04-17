# Feature Research

**Domain:** Terminal UI (TUI) for browsing RedFlagDeals.com forums
**Researched:** 2026-04-17
**Confidence:** HIGH

## Feature Landscape

### Table Stakes (Users Expect These)

Features users assume exist in any TUI forum/deal browser. Missing these = product feels incomplete. Derived from analysis of hackernews-TUI (696 stars), michael-lazar/rtv (Reddit Terminal Viewer), guyfedwards/nom (RSS reader), and davegallant/rfd (reference RFD CLI).

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Deal list view | Core value prop — users launch to see deals | MEDIUM | Use Bubble Tea `list` component. Display title, score, views, replies, dealer, category, age per row. RFD API returns all fields in one call. |
| Vim-style navigation (j/k/enter/q) | Target audience is CLI users; every comparable TUI (hn-tui, rtv, nom, newsboat) defaults to vim keys | LOW | Bubble Tea's `tea.KeyPressMsg` makes this trivial. Also support arrow keys for accessibility. |
| Thread/post detail view | Users need to read original post and replies before deciding on a deal | MEDIUM | Use Bubble Tea `viewport` component for scrollable content. RFD API: `GET /api/topics/{id}/posts`. Lazy-load replies on demand (not all at once). |
| Open deal in browser | Users must see full deal page, images, purchase links — TUI can't replace this | LOW | `exec.Command("open"/"xdg-open", url)`. Single keypress `o`. Every comparable TUI has this. |
| Search/filter | Finding specific deals is the #2 use case after browsing (proven by rfd CLI's search command) | MEDIUM | Bubble Tea `textinput` component. Search across titles and dealer names. RFD API supports regex via client-side filtering. |
| Sort by score/views | Deal quality ranking is fundamental; rfd CLI has `--sort-by` as primary flag | LOW | Client-side sort of already-fetched deals. No API parameter needed. |
| Pagination | Hot Deals has 50+ pages; users need to browse beyond first 20 results | MEDIUM | RFD API has `page` param and `pager.total_pages`. Use Bubble Tea `paginator` or custom. `n`/`p` keys like hn-tui. |
| Help overlay (? key) | Users can't memorize all shortcuts; hn-tui uses `?` to show help in every view | LOW | Modal overlay with keybinding list. Render with Lipgloss. Dismissible with `esc`. |
| Quit cleanly (q/Ctrl+C) | Every TUI must exit cleanly, restoring terminal state | LOW | Bubble Tea handles this natively with `tea.Quit`. |
| Loading spinner | API calls take 1-3 seconds; users need feedback during fetch | LOW | Bubble Tea `spinner` component. Show during API calls. |

### Differentiators (Competitive Advantage)

Features that set RFD TUI apart from the reference Python CLI and other TUI tools. Not required for MVP, but valuable.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Deal-specific rich display | RFD has unique fields (dealer, price, savings %, votes) that no generic TUI shows. Color-coded scores (green=hot, red=cold) give instant visual scanning. | MEDIUM | Lipgloss styling. RFD API returns `votes.up`, `votes.down`, `dealer`, `price`. Score coloring: >50 green, 20-50 yellow, <20 red. |
| Category filtering | Hot Deals has sub-categories (Electronics, Home, etc.). Filtering by category is impossible in the web UI's linear list. | MEDIUM | RFD API response includes `topic.category_id`. Client-side filter. Tab/shift-tab to cycle categories, or `/` filter prompt. |
| Minimum score threshold | Hide low-quality deals automatically. Only show deals above configurable score. Unique to TUI — web UI doesn't offer this. | LOW | Simple client-side filter. Default threshold configurable. `f` key to adjust. |
| Regex search | davegallant/rfd supported regex (`rfd search '(coffee\|starbucks)'`). Power users expect this. | LOW | Go `regexp` package. Compile user input as regex, match against title and dealer fields. Fallback to plain text if regex fails. |
| Multiple forum sections | RFD has Hot Deals (9), Freebies (15), Contests (12), etc. Quick-switch with F-keys (like hn-tui's F1-F5). | MEDIUM | RFD API uses `forum_id` param. Hardcode known IDs. F1=Hot Deals, F2=Freebies, F3=Contests. Requires minimal extra API logic — same endpoint, different ID. |
| Configurable keybindings | Power users customize keys. hn-tui and gh-dash both support this. Differentiates from static CLI. | MEDIUM | TOML config file at `~/.config/rfd-tui/config.toml`. Parse on startup. Map actions to key strings. `key.NewBinding()` in Bubble Tea. |
| Collapsible comment threads | hn-tui supports tab to collapse. Essential for long RFD threads (100+ replies). Navigate top-level only. | MEDIUM | Track collapse state per comment. RFD API returns nested posts with `parent_id`. Render only visible comments. `tab` to toggle, `u` for parent. |
| Deal age display | "2h ago", "3d ago" — relative timestamps help users spot fresh deals. Web UI shows absolute time. | LOW | Calculate from RFD API's `created_at` field. Go `time.Since()`. Simple helper function. |
| JSON output mode | davegallant/rfd had `--output json`. Useful for piping to `jq`, scripting, automation. | LOW | Struct output to stdout when flag passed. Skip TUI entirely. Go `encoding/json`. |
| Color theme support | gh-dash and hn-tui both support themes. RFD users have strong preferences (some want minimal, some want colorful). | MEDIUM | Lipgloss color profiles. 2-3 built-in themes (dark, light, minimal). Configurable via TOML. |

### Anti-Features (Commonly Requested, Often Problematic)

| Anti-Feature | Why Requested | Why Problematic | Alternative |
|--------------|---------------|-----------------|-------------|
| Real-time auto-refresh | Users want to see new deals appear live | Constant API polling. Rate limiting risk (RFD is a third-party API, no official rate limit docs). Unnecessary complexity for a browsing tool. Breaks the "open, scan, close" workflow. | Manual refresh with `r` key. Users control when to fetch. Simple, predictable, no rate limit issues. |
| User authentication / posting | "Why can't I upvote or reply?" | RFD has no official API for auth. Reverse-engineering session cookies is fragile, breaks on site changes, and may violate ToS. Massive scope expansion. Reference project (davegallant/rfd) explicitly stayed read-only. | Open deal in browser (`o` key) for any interaction. TUI is for browsing/discovery, browser is for action. |
| Image rendering in terminal | Deal images help evaluate products | Terminal image rendering (sixel, kitty, sixel) is fragmented — different terminals support different protocols. Most users' terminals won't render images. Massive compatibility headache. | Open deal in browser for images. Some terminals can use `chafa` for ASCII art, but don't make it a core feature. |
| Notifications / deal alerts | "Notify me when a PS5 deal appears" | davegallant/rfd-notify already exists as a separate tool. Requires background daemon, push notification infrastructure, persistent state. Completely different product category. | Point users to rfd-notify for alert-based workflows. TUI is for active browsing. |
| Bookmarks / saved deals | Users want to save deals for later | Requires local storage, persistence layer, state management. Adds CRUD complexity. Offline-first design needed. Scope creep into "deal management" territory. | Open in browser and use browser bookmarks. Or pipe JSON output to a file. |
| Offline mode / caching | "Cache deals so I can read on the subway" | Requires local SQLite/file storage, cache invalidation, staleness logic. Complex for a browse-first tool. RFD deals expire quickly — stale cache is misleading. | Fresh fetch on every launch. Deals change hourly; cached data is low-value. |
| HTML content rendering | RFD posts have rich HTML (images, embeds, formatted text) | Full HTML rendering in terminal is unsolved. Even w3m/lynx struggle with modern HTML. Glamour (markdown renderer) can't parse HTML. Would need HTML→markdown conversion. | Strip HTML, render plain text with basic formatting (bold, links). Use `html.Parse` to extract text content. Good enough for deal browsing. |

## Feature Dependencies

```
[Deal List View]
    └──requires──> [RFD API Client]
                       └──requires──> [HTTP Client + JSON Parsing]

[Thread Detail View]
    └──requires──> [Deal List View] (must select a deal first)
    └──requires──> [RFD API Client - Posts endpoint]

[Search/Filter]
    └──enhances──> [Deal List View]
    └──requires──> [Text Input Component]

[Sort by Score/Views]
    └──enhances──> [Deal List View]
    └──independent (client-side only)

[Category Filtering]
    └──enhances──> [Deal List View]
    └──requires──> [Deal List with category data from API]

[Score Threshold]
    └──enhances──> [Deal List View]
    └──independent (client-side filter)

[Pagination]
    └──requires──> [RFD API Client - page parameter]
    └──requires──> [Deal List View]

[Open in Browser]
    └──requires──> [Deal List View] or [Thread Detail View]
    └──independent (just needs a URL)

[Collapsible Comments]
    └──requires──> [Thread Detail View]
    └──requires──> [Comment tree parsing from API]

[Configurable Keybindings]
    └──independent (config file parsing)
    └──enhances──> [All Views]

[Multiple Forum Sections]
    └──requires──> [RFD API Client - forum_id parameter]
    └──requires──> [Deal List View] (reuse same list component)

[Help Overlay]
    └──requires──> [All Views] (must know available keybindings)
    └──independent (render logic)

[Regex Search]
    └──requires──> [Search/Filter]
    └──enhances──> [Search/Filter]

[JSON Output Mode]
    └──conflicts──> [TUI Mode] (mutually exclusive — either TUI or JSON output)
    └──requires──> [RFD API Client]
```

### Dependency Notes

- **Deal List View requires RFD API Client:** The API client is the foundation. Everything depends on fetching deals from `GET /api/topics?forum_id=9`. Build this first.
- **Thread Detail View requires Deal List View:** Users select a deal from the list, then open thread detail. Navigation flows list → detail → back.
- **Search/Sort/Filter enhance Deal List View:** These all operate on the fetched deals. They don't require additional API calls — all client-side. Can be added incrementally.
- **Pagination requires RFD API Client:** The API returns `pager.total_pages`. Pagination is an API feature, not just UI. Must track current page and make new API calls for next/prev.
- **JSON Output conflicts with TUI Mode:** When `--json` flag is passed, skip TUI entirely, output JSON to stdout, and exit. This is a CLI mode, not a TUI mode.
- **Multiple Forum Sections requires same API endpoint:** Just changes the `forum_id` parameter. Reuses entire Deal List infrastructure. Very low marginal cost.
- **Collapsible Comments requires Thread Detail View and comment tree parsing:** RFD API returns flat posts with `parent_id`. Must build tree structure in memory before rendering with collapse state.

## MVP Definition

### Launch With (v1)

Minimum viable product — what's needed to validate the concept. Users should be able to browse, search, and read hot deals.

- [ ] **RFD API Client** — HTTP client that fetches deals and threads from RFD JSON API. Foundation for everything.
- [ ] **Deal List View** — Display hot deals with title, score, views, replies, dealer, category, age. Scrollable list.
- [ ] **Vim Navigation** — j/k scroll, enter open, q/esc back, Ctrl+C quit. Also support arrow keys.
- [ ] **Thread Detail View** — Original post + expandable replies. Scrollable viewport.
- [ ] **Open in Browser** — `o` key opens selected deal URL in system browser.
- [ ] **Sort by Score/Views** — Client-side sort of displayed deals.
- [ ] **Search** — Keyword search across titles and dealer names.
- [ ] **Pagination** — Next/prev page through RFD's paginated results.
- [ ] **Help Overlay** — `?` key shows available shortcuts.
- [ ] **Loading Spinner** — Visual feedback during API fetches.

### Add After Validation (v1.x)

Features to add once core browsing works and users validate the concept.

- [ ] **Regex Search** — Support regex in search input. Trigger: users ask for power search.
- [ ] **Category Filtering** — Filter by deal category. Trigger: users want to focus on specific categories.
- [ ] **Score Threshold** — Hide deals below configurable score. Trigger: users complain about low-quality deals.
- [ ] **Deal Age Display** — "2h ago" relative timestamps. Trigger: users want to spot fresh deals fast.
- [ ] **Collapsible Comments** — Tab to collapse reply threads. Trigger: users struggle with long threads.
- [ ] **Color-coded Scores** — Green/yellow/red based on deal score. Trigger: visual scanning improvement.
- [ ] **Multiple Forum Sections** — F1=Hot Deals, F2=Freebies, F3=Contests. Trigger: users request other forums.
- [ ] **JSON Output Mode** — `--json` flag for scripting/automation. Trigger: power users want to pipe output.

### Future Consideration (v2+)

Features to defer until product-market fit is established. These require significant additional work.

- [ ] **Configurable Keybindings** — TOML config file. Why defer: requires config parsing, validation, migration. Add after users request specific key changes.
- [ ] **Color Themes** — Multiple built-in themes (dark, light, minimal). Why defer: cosmetic. Ship one good default first.
- [ ] **Shell Completions** — bash/zsh completion for CLI flags. Why defer: only useful for CLI mode, not TUI.
- [ ] **Markdown Rendering** — Use glamour for rich post formatting. Why defer: RFD posts are HTML, not markdown. Would need HTML→markdown conversion first.

## Feature Prioritization Matrix

| Feature | User Value | Implementation Cost | Priority |
|---------|------------|---------------------|----------|
| RFD API Client | HIGH | MEDIUM | P1 |
| Deal List View | HIGH | MEDIUM | P1 |
| Vim Navigation | HIGH | LOW | P1 |
| Thread Detail View | HIGH | MEDIUM | P1 |
| Open in Browser | HIGH | LOW | P1 |
| Sort by Score/Views | HIGH | LOW | P1 |
| Search (keyword) | HIGH | MEDIUM | P1 |
| Pagination | HIGH | MEDIUM | P1 |
| Help Overlay | MEDIUM | LOW | P1 |
| Loading Spinner | MEDIUM | LOW | P1 |
| Color-coded Scores | MEDIUM | LOW | P2 |
| Deal Age Display | MEDIUM | LOW | P2 |
| Category Filtering | MEDIUM | MEDIUM | P2 |
| Score Threshold | MEDIUM | LOW | P2 |
| Regex Search | MEDIUM | LOW | P2 |
| Collapsible Comments | MEDIUM | MEDIUM | P2 |
| Multiple Forum Sections | MEDIUM | MEDIUM | P2 |
| JSON Output Mode | LOW | LOW | P2 |
| Configurable Keybindings | LOW | MEDIUM | P3 |
| Color Themes | LOW | MEDIUM | P3 |
| Shell Completions | LOW | LOW | P3 |

**Priority key:**
- P1: Must have for launch — core browsing loop
- P2: Should have, add when possible — enhanced browsing
- P3: Nice to have, future consideration — customization

## Competitor Feature Analysis

| Feature | davegallant/rfd (Python CLI) | aome510/hackernews-TUI (Rust TUI) | michael-lazar/rtv (Reddit TUI) | Our Approach |
|---------|------------------------------|-----------------------------------|-------------------------------|--------------|
| Interface | CLI (pager output) | Full TUI (Cursive) | Full TUI (curses) | Full TUI (Bubble Tea) |
| Navigation | CLI args only | Vim keys + F-keys | Vim keys | Vim keys + F-keys |
| Deal/Story List | `rfd threads` → pager | Persistent scrollable list | Persistent scrollable list | Persistent scrollable list (Bubble Tea `list`) |
| Thread View | `rfd posts <url>` → pager | Scrollable with collapse | Scrollable with collapse | Scrollable viewport with lazy-load replies |
| Search | `rfd search 'regex'` | Search view with mode toggle | Subreddit search | Inline search/filter (like hn-tui) |
| Sort | `--sort-by score\|views` | Cycle sort with `d` key | Hot/top/new sort | Sort toggle key, client-side |
| Open in Browser | No (CLI) | `o` key | `o` key | `o` key (essential for deals) |
| Pagination | `--pages N` | `n`/`p` keys | Infinite scroll | `n`/`p` keys |
| Configurable Keys | No | Yes (TOML) | Yes (config) | Yes (TOML), deferred to v2 |
| Categories/Forums | No | Multiple HN types (F1-F5) | Subreddit navigation | Multiple RFD forums (F1-F3), deferred |
| Auth/Posting | No | Yes (voting) | Yes (full) | No — read-only browsing |
| JSON Output | Yes (`--output json`) | No | No | Yes, deferred to v1.x |
| Help | `--help` flag | `?` overlay | `?` overlay | `?` overlay |
| Deal-specific Fields | Yes (score, views, dealer) | N/A (HN has different fields) | N/A | Yes, with color coding |
| Collapsible Comments | No (flat pager) | Yes (tab to collapse) | Yes (space to collapse) | Yes, deferred to v1.x |

## Sources

- **aome510/hackernews-TUI** (696 stars, active) — Feature-rich HN TUI in Rust. Gold standard for TUI forum clients. Features: multi-view, search modes, configurable keys, auth, article view, comment navigation. GitHub: `aome510/hackernews-TUI`
- **davegallant/rfd** (archived Aug 2024) — Reference Python CLI for RFD. Proved RFD JSON API works. Features: threads, search with regex, sort, posts view, JSON output, shell completion. GitHub: `davegallant/rfd`
- **michael-lazar/rtv** (archived) — Classic Reddit Terminal Viewer in Python. Set the standard for Reddit TUI features: vim navigation, comment collapsing, themes, auth, subreddit navigation. GitHub: `michael-lazar/rtv`
- **guyfedwards/nom** (active) — RSS reader TUI built with Bubble Tea. Closest tech stack match. Features: vim keys, config file, mark read/unread, filtering, open in browser. GitHub: `guyfedwards/nom`
- **dlvhdr/gh-dash** (active, popular) — GitHub TUI with Bubble Tea. Features: custom keybindings, YAML config, custom actions, multiple sections. GitHub: `dlvhdr/gh-dash`
- **davegallant/rfd-notify** — Separate deal notification tool. Validates that deal alerts are a separate product category, not a TUI feature. GitHub: `davegallant/rfd-notify`
- **davegallant/rfd-fyi** — Alternative web frontend for RFD. Vue 3 + Go backend. Validates RFD API usage patterns. GitHub: `davegallant/rfd-fyi`
- **Charm Bracelet Bubble Tea docs** — Verified via Context7: `list`, `viewport`, `paginator`, `spinner`, `table`, `textinput` components available in Bubbles v2 for all P1 features.
- **PROJECT.md** — Project context: Go + Bubble Tea, RFD JSON API, vim-style nav, single binary distribution.

---
*Feature research for: RFD TUI (Go terminal client for RedFlagDeals.com)*
*Researched: 2026-04-17*
