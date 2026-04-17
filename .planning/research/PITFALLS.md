# Pitfalls Research

**Domain:** Go TUI application consuming RFD JSON REST API
**Researched:** 2026-04-17
**Confidence:** HIGH (verified against Bubble Tea v2 source, Context7 docs, and real-world Charmbracelet ecosystem projects)

## Critical Pitfalls

### Pitfall 1: Using Bubble Tea v1 APIs Instead of v2

**What goes wrong:**
Most tutorials, blog posts, and Stack Overflow answers use Bubble Tea v1 imports (`github.com/charmbracelet/bubbletea`) and v1 patterns (`tea.KeyMsg`, `tea.EnterAltScreen`, `tea.WithAltScreen()`). The v2 API is fundamentally different: `View()` returns a `tea.View` struct instead of a string, key handling uses `tea.KeyPressMsg` instead of `tea.KeyMsg`, and terminal features are declarative fields on the View struct instead of imperative commands. Copying v1 code causes compilation errors, confusing runtime behavior, and a frustrating migration later.

**Why it happens:**
Bubble Tea v2 changed import paths to `charm.land/bubbletea/v2` and overhauled the architecture. The v1 tutorials still rank high in search results. GitHub Copilot and similar tools still suggest v1 patterns from training data.

**How to avoid:**
- Use `charm.land/bubbletea/v2` import path exclusively
- Use `charm.land/lipgloss/v2` and `charm.land/bubbles/v2` for matching v2 ecosystem
- In `View()`, return `tea.NewView(content)` with declarative fields (`.AltScreen = true`, `.MouseMode = tea.MouseModeCellMotion`)
- Key matching: use `tea.KeyPressMsg` (not `tea.KeyMsg`), match with `msg.String()` returning `"space"` not `" "`, and `msg.String() == "ctrl+c"` not `tea.KeyCtrlC`
- Reference the official upgrade guide: `github.com/charmbracelet/bubbletea/blob/main/UPGRADE_GUIDE_V2.md`

**Warning signs:**
- Import path contains `github.com/charmbracelet/bubbletea` (v1) instead of `charm.land/bubbletea/v2`
- `View()` method returns `string` instead of `tea.View`
- Code uses `tea.WithAltScreen()`, `tea.EnterAltScreen`, or `tea.KeyMsg`
- Type switches on `tea.KeyMsg` instead of `tea.KeyPressMsg`

**Phase to address:**
Phase 1 (project scaffolding). Getting the wrong version is a showstopper that compounds through every subsequent phase.

---

### Pitfall 2: Blocking in Update() Instead of Using tea.Cmd

**What goes wrong:**
Performing HTTP requests, file I/O, or any blocking operation inside the `Update()` method freezes the entire TUI. The terminal stops rendering, spinner animations freeze, and key presses are queued but not processed until the blocking operation completes. Users see a hung application.

**Why it happens:**
Bubble Tea's Elm architecture is unfamiliar to developers coming from imperative UI frameworks. The `Update()` function runs on the main render loop. Blocking it blocks everything. New developers intuitively write `data, err := fetchDeals()` inside `Update()` because that's how normal Go code works.

**How to avoid:**
- All I/O operations must be wrapped in `tea.Cmd` functions (functions that return `tea.Msg`)
- Return the command from `Update()`: `return m, fetchDeals()`
- Use `tea.Batch()` for concurrent commands (e.g., loading deals + spinner tick simultaneously)
- Use `tea.Sequence()` when commands must run in order
- Pattern for API calls:

```go
type dealsLoadedMsg struct{ deals []Deal; err error }

func fetchDeals(page int) tea.Cmd {
    return func() tea.Msg {
        client := &http.Client{Timeout: 10 * time.Second}
        resp, err := client.Get("https://forums.redflagdeals.com/api/topics?forum_id=9&per_page=20&page=" + strconv.Itoa(page))
        if err != nil {
            return dealsLoadedMsg{err: err}
        }
        defer resp.Body.Close()
        // parse and return
        return dealsLoadedMsg{deals: deals}
    }
}

// In Update():
case tea.KeyPressMsg:
    if msg.String() == "r" {
        m.loading = true
        return m, tea.Batch(fetchDeals(m.page), spinner.Tick)
    }
case dealsLoadedMsg:
    m.loading = false
    m.deals = msg.deals
    return m, nil
```

**Warning signs:**
- Any `http.Get`, `json.Decoder`, `os.Read`, or `time.Sleep` inside `Update()`
- View contains spinner but spinner never visually rotates
- Variable named `loading` or `fetching` set to true but UI never shows the loading state

**Phase to address:**
Phase 1 (first API call). This must be correct from the very first network request.

---

### Pitfall 3: Ignoring WindowSizeMsg (Terminal Resize)

**What goes wrong:**
Deal list renders with wrong column widths after terminal resize. Content overflows or underflows. Viewport doesn't adjust. Tables break. The TUI looks fine at launch but becomes garbled when the user resizes their terminal or uses a split pane.

**Why it happens:**
Developers test at a fixed terminal size during development and never handle `tea.WindowSizeMsg`. Bubble Tea sends this message on startup and every resize event, but if you don't handle it, your layout uses stale dimensions.

**How to avoid:**
- Always store `width` and `height` from `tea.WindowSizeMsg` in your model
- Recalculate all layout dimensions when size changes
- Use Lipgloss `Width()` and `MaxWidth()` styles that respond to stored terminal width
- Re-set viewport dimensions on resize: `m.viewport.Width = msg.Width; m.viewport.Height = msg.Height - headerHeight - footerHeight`
- Use `lipgloss.Width(renderedString)` to measure actual rendered width of styled content
- Every Charmbracelet example and production project handles `WindowSizeMsg` — it's not optional

**Warning signs:**
- No `case tea.WindowSizeMsg:` in the Update switch
- Model has no `width`, `height` fields
- Hard-coded column widths or layout dimensions
- Testing only at one terminal size

**Phase to address:**
Phase 1 (deal list rendering). WindowSizeMsg handling should be built into the initial deal list view.

---

### Pitfall 4: Monolithic Model That Becomes Unmanageable

**What goes wrong:**
A single `model` struct accumulates all state (deals, threads, search, filters, pagination, viewport, cursor, loading states, errors) and the `Update()` method grows to hundreds of lines with nested switches. Adding a new feature requires touching the same giant file. State management becomes error-prone — clearing loading state for one feature accidentally clears it for another.

**Why it happens:**
Bubble Tea tutorials start with a single struct for simplicity. There's no framework-enforced pattern for decomposition. Developers keep adding fields to the same struct because it's the path of least resistance.

**How to avoid:**
- Use a ViewID enum pattern for multi-screen navigation:

```go
type viewID int
const (
    dealListID viewID = iota
    threadDetailID
    searchID
)

type model struct {
    currentView viewID
    width, height int
    // Each view gets its own sub-model
    list   listModel
    thread threadModel
    search searchModel
}
```

- Delegate `Update` and `View` based on `currentView`:

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
        // propagate to sub-models
    case tea.KeyPressMsg:
        // global keys (quit, etc.)
    }
    // delegate to current view
    switch m.currentView {
    case dealListID:
        return m.list.Update(msg, m.width, m.height)
    case threadDetailID:
        return m.thread.Update(msg, m.width, m.height)
    }
}
```

- Real-world precedent: Grafana Loki's bench TUI (`pkg/logql/bench/cmd/bench/`) uses exactly this pattern with `ViewID`, `ListView`, `RunView` sub-models, and `SwitchViewMsg` for transitions.
- Each sub-model owns its own state. Transitions between views use custom messages.

**Warning signs:**
- Model struct has 20+ fields
- `Update()` method exceeds 100 lines
- Multiple unrelated `loading` boolean fields (`dealsLoading`, `threadLoading`, `searchLoading`)
- Adding a new screen requires modifying the same Update switch statement

**Phase to address:**
Phase 1 (project scaffolding). The ViewID + sub-model pattern should be established in the initial structure before any features are added. Retroactively decomposing a monolithic model is painful.

---

### Pitfall 5: Not Setting HTTP Client Timeout

**What goes wrong:**
Using `http.DefaultClient` (which has no timeout) means an unresponsive RFD API server will hang the TUI indefinitely. The loading spinner runs forever. The user has no choice but to force-kill the application. Even worse, the goroutine inside the `tea.Cmd` leaks, holding resources.

**Why it happens:**
Go's `http.Get()` uses `http.DefaultClient` which has no timeout. Tutorials often skip timeout configuration. During development, requests always succeed quickly, masking the issue.

**How to avoid:**
- Always create an `http.Client` with a timeout for API calls
- Charmbracelet's own projects use 30 seconds (`client := &http.Client{Timeout: 30 * time.Second}`)
- For a deal-browsing CLI, 10-15 seconds is appropriate (users expect snappy responses)
- Create a shared client at app initialization, reuse across all requests
- Handle timeout errors gracefully in the UI (show "request timed out" message, offer retry)

```go
type apiClient struct {
    http    *http.Client
    baseURL string
}

func newAPIClient() *apiClient {
    return &apiClient{
        http:    &http.Client{Timeout: 15 * time.Second},
        baseURL: "https://forums.redflagdeals.com",
    }
}
```

**Warning signs:**
- `http.Get()` or `http.DefaultClient` used directly
- No timeout configured anywhere in the HTTP client
- No error handling for network failures in the UI
- No way for users to cancel a stuck request

**Phase to address:**
Phase 1 (first API integration). The HTTP client should be properly configured from the first request.

---

### Pitfall 6: Not Sanitizing API Content for Terminal Rendering

**What goes wrong:**
RFD forum posts contain HTML markup, special characters, emoji, and ANSI escape sequences. Rendering these raw in the terminal causes garbled output, broken layouts, and potential terminal injection. Deal titles with HTML entities (`&amp;`, `&#39;`) display literally instead of as the intended characters. Unicode emoji from forum posts may not render in all terminals, causing misaligned columns.

**Why it happens:**
The RFD JSON API returns semi-HTML content in post bodies and sometimes in titles. Developers test with simple ASCII deal titles and don't encounter the issue until real-world data with special characters flows through.

**How to avoid:**
- Strip HTML tags from post content before rendering (use a simple regex or `html.UnescapeString()` for entities)
- Truncate long titles to fit column widths using `lipgloss.NewStyle().MaxWidth(n).Render(title)`
- Handle emoji and wide characters: use `lipgloss.Width()` (which accounts for ANSI escape codes and wide characters) instead of `len(string)` for measuring rendered width
- Test with real RFD API data early — not mock data with clean ASCII strings
- Consider replacing or removing emoji that don't render well in common terminals

**Warning signs:**
- Deal titles displaying `&amp;` or `&#39;` literally
- Column alignment breaks on certain deals
- Using `utf8.RuneCountInString()` or `len()` for width calculation (wrong — doesn't account for ANSI codes or wide runes)
- No HTML stripping or entity decoding before rendering

**Phase to address:**
Phase 1 (deal list rendering). Fetch real API data from day one and test sanitization against actual RFD content.

---

### Pitfall 7: Loading Entire Thread on Open (No Lazy Loading)

**What goes wrong:**
Popular RFD deal threads can have hundreds or thousands of replies. Loading the entire thread (all posts) when the user presses Enter causes a long wait, consumes excess memory, and provides a poor experience when the user just wants to read the original deal post.

**Why it happens:**
It's simpler to load all posts in one request. The RFD API supports `per_page` pagination on the posts endpoint, but developers may not discover this until threads load slowly.

**How to avoid:**
- Load only the original post (first post, page 1, `per_page=1`) when opening a thread
- Show "Load more replies" or expand-on-demand for subsequent posts
- Use the RFD API's built-in pagination: `GET /api/topics/{id}/posts?per_page={n}&page={n}`
- The `pager.total_pages` field in the API response tells you how many pages of replies exist
- Display reply count as metadata so users know what they're getting into

**Warning signs:**
- API request for posts uses a very high `per_page` value
- No pagination logic for the posts endpoint
- Thread view loads all posts at once before rendering anything
- Users report slow thread opening

**Phase to address:**
Phase 2 (thread detail view). Design thread loading with lazy pagination from the start.

---

## Technical Debt Patterns

Shortcuts that seem reasonable but create long-term problems.

| Shortcut | Immediate Benefit | Long-term Cost | When Acceptable |
|----------|-------------------|----------------|-----------------|
| Monolithic model (one big struct) | Faster to start coding | Update() becomes unmaintainable, merge conflicts, state bleed between views | Never — decompose from the start with ViewID pattern |
| Hard-coded terminal widths | Layout works on your machine | Breaks on resize, different terminals, split panes | Never — use WindowSizeMsg from day one |
| `http.Get()` without timeout | Less code, works in dev | Hangs in production, no cancellation, leaked goroutines | Never — always set timeout |
| Fetching all posts at once | Simpler thread view logic | Slow for popular threads, excessive memory, poor UX | MVP only — add pagination in same phase |
| Skipping error states in UI | Cleaner View() code | Users see blank screens or stale data with no feedback | Never — show loading/error/empty states for every async operation |
| Raw API responses in model | No data transformation layer | Tight coupling to API format, hard to change display format | MVP only — add a thin domain layer when model complexity grows |
| No caching of API responses | Simpler code | Unnecessary network requests on back navigation | MVP — acceptable initially, add before pagination |
| Ignoring Lipgloss v2 import path | Copy-paste from v1 tutorials | Must rewrite all styling when upgrading | Never — use `charm.land/lipgloss/v2` from the start |

## Integration Gotchas

Common mistakes when connecting to the RFD API.

| Integration | Common Mistake | Correct Approach |
|-------------|----------------|------------------|
| RFD Threads API | Assuming `forum_id` is always 9 (Hot Deals) | Support multiple forum IDs from the start — define as constants, not hard-coded values |
| RFD Posts API | Loading all pages of posts at once | Load page 1 (original post) first, expand replies on demand with pagination |
| Pagination | Trusting `page` count without checking `total_pages` | Always read `pager.total_pages` from response, disable "next page" when on last page |
| Rate Limiting | Assuming no rate limits (public API) | Add reasonable request throttling; don't fire requests on every keypress during search |
| HTTP Headers | Sending no User-Agent header | Set a User-Agent header (e.g., `rfd-tui/1.0`) — APIs may block default Go HTTP client User-Agent |
| Error Responses | Only checking `err` from `json.Decode` | Check HTTP status code first; RFD may return HTML error pages instead of JSON on errors |
| Deal Score | Assuming score is always a positive integer | RFD deal scores can be negative (downvoted deals). Handle negative display and sorting correctly |
| Deal URLs | Assuming `deal_url` is always present | Some deals are discussion-only with no external URL. Handle the "open URL" action when URL is empty |
| Thread Content | Assuming posts are plain text | RFD posts contain HTML (bold, links, images, quotes). Must strip HTML and decode entities for terminal display |
| API Stability | Assuming the API won't change | The reference Python project (davegallant/rfd) was archived — API may have undocumented changes. Add response parsing resilience |

## Performance Traps

Patterns that work at small scale but fail as usage grows.

| Trap | Symptoms | Prevention | When It Breaks |
|------|----------|------------|----------------|
| Re-rendering entire deal list on every message | Sluggish key response, flickering | Only re-render when data changes; use viewport for virtual scrolling (Bubbles viewport handles this) | 50+ items in list |
| Storing full post HTML in memory for all loaded threads | High memory usage | Store only what's needed for display; strip HTML on fetch, not on render | 10+ cached threads |
| No request deduplication | Duplicate API calls on rapid keypresses (e.g., pressing 'j' quickly) | Debounce or track in-flight requests with a `loading` guard | Any time user navigates rapidly |
| Synchronous JSON decoding in tea.Cmd | Brief but noticeable freeze on large responses | JSON decode already happens in Cmd goroutine (correct), but avoid doing heavy transformation there — do it in Update | 1000+ post threads |
| Growing deal list without bounds | Memory grows linearly with pages fetched | Cap cached pages; evict old pages when exceeding limit | 10+ pages of deals loaded |
| Full View() re-render on every tick (spinner) | Constant re-rendering of unchanged content | Bubbles spinner handles this correctly via messages — just don't do extra work in View() | Immediately if View() is expensive |

## Security Mistakes

Domain-specific security issues beyond general web security.

| Mistake | Risk | Prevention |
|---------|------|------------|
| Rendering raw API HTML with ANSI escapes | Terminal injection, garbled display | Strip all HTML tags and ANSI escape sequences before rendering |
| Logging full API responses | PII in logs (usernames, user IDs from posts) | Log only metadata (status, timing), not response bodies |
| No request timeout | Denial of service if RFD API is unresponsive | Always set `http.Client.Timeout` (10-15s) |
| Ignoring TLS certificate errors | MITM attacks in coffee shops, etc. | Never use `InsecureSkipVerify: true` — the RFD API uses proper HTTPS |
| Embedding credentials in binary | If auth is ever added, credentials in source | Use environment variables or config files for any future auth tokens |

## UX Pitfalls

Common user experience mistakes in Go TUIs and deal browsers.

| Pitfall | User Impact | Better Approach |
|---------|-------------|-----------------|
| No loading indicator during API calls | User thinks app is frozen, force-kills it | Show spinner + "Loading deals..." text for every async operation |
| No error message when API fails | User sees blank screen with no explanation | Show styled error message with retry instructions ("Press r to retry") |
| Vim keys without fallback for non-vim users | Arrow keys don't work, users are stuck | Support both vim keys (j/k) AND arrow keys (up/down) |
| No way to open deal URL | User must manually copy URL from title | Press 'o' opens in browser using `browser.OpenURL()` — show this in help |
| No help/legend visible | Users don't know available keybindings | Show a persistent footer with common keys, or '?' for full help |
| Deals list doesn't show enough context | User must open each deal to see if it's interesting | Show title, score, dealer, category, age, reply count in the list row |
| Sort/filter requires restart | Users can't refine their view | Allow runtime sort/filter toggles with immediate re-fetch |
| No feedback on keypress | User presses a key, nothing visible happens | Visual feedback: cursor movement, status bar updates, highlights |
| Q to quit without confirmation when loading | Data loss if user accidentally presses q during load | Only require confirmation if there's unsaved state (bookmarks in future); for read-only, instant quit is fine |
| Not restoring terminal state on crash | Terminal left in raw mode, garbled output | Bubble Tea handles this automatically via its runtime — but ensure you use `tea.Quit` properly and don't panic in View() |

## "Looks Done But Isn't" Checklist

Things that appear complete but are missing critical pieces.

- [ ] **Deal list:** Often missing keyboard scrolling for lists longer than terminal height — verify viewport works with 50+ deals
- [ ] **Thread detail:** Often missing pagination (only loads first page of replies) — verify "load more" for 5+ page threads
- [ ] **Search:** Often missing debouncing (fires API on every keystroke) — verify search triggers on Enter or with delay
- [ ] **Sort/filter:** Often missing re-fetch after sort order changes — verify new API call after sort toggle
- [ ] **Open in browser:** Often missing error handling when no URL exists — verify behavior on discussion-only deals
- [ ] **Window resize:** Often missing re-render on resize — verify at 80x24, 120x40, and narrow widths
- [ ] **Empty states:** Often missing "no deals found" or "no search results" messages — verify behavior with zero results
- [ ] **Error recovery:** Often missing retry after network failure — verify user can continue after transient error
- [ ] **Pagination:** Often missing "last page" detection — verify "next page" disabled on final page
- [ ] **HTML in titles/posts:** Often missing entity decoding — verify with real RFD data containing `&amp;`, `&#39;`, `<b>` tags
- [ ] **Negative scores:** Often missing handling for downvoted deals — verify sorting and display with negative values
- [ ] **Quit behavior:** Often missing terminal cleanup on Ctrl+C — verify terminal is usable after exit

## Recovery Strategies

When pitfalls occur despite prevention, how to recover.

| Pitfall | Recovery Cost | Recovery Steps |
|---------|---------------|----------------|
| Built on v1 APIs | MEDIUM | Migrate imports to `charm.land/bubbletea/v2`, refactor View() to return `tea.View`, update key handling to `KeyPressMsg` — the upgrade guide is comprehensive |
| Monolithic model | HIGH | Decompose into sub-models with ViewID pattern; requires restructuring Update/View and all state management |
| No window resize handling | LOW | Add `width, height int` to model, handle `WindowSizeMsg`, recompute layout in View — incremental fix |
| Blocking in Update() | LOW | Wrap blocking calls in `tea.Cmd` functions, return from Update() — mechanical fix per call site |
| No HTTP timeout | LOW | Create shared `http.Client{Timeout: 15s}`, replace all `http.Get()` calls — quick fix |
| HTML in content not stripped | LOW | Add `html.UnescapeString()` and tag stripping to render path — can be fixed per view |
| Loading entire threads | MEDIUM | Add pagination to posts endpoint, lazy-load replies — requires UX changes to thread view |
| No error/loading states | MEDIUM | Add loading/error fields per sub-model, update View() to render them — systematic but mechanical |

## Pitfall-to-Phase Mapping

How roadmap phases should address these pitfalls.

| Pitfall | Prevention Phase | Verification |
|---------|------------------|--------------|
| v1 vs v2 API confusion | Phase 1: Scaffolding | `go.mod` contains `charm.land/bubbletea/v2`; all imports use `charm.land/` vanity domains |
| Blocking in Update() | Phase 1: First API call | No `http.Get()` or blocking I/O in Update(); all async via `tea.Cmd` |
| WindowSizeMsg handling | Phase 1: Deal list render | Resize terminal while app runs — layout adjusts without garbling |
| Monolithic model decomposition | Phase 1: Scaffolding | ViewID enum exists; sub-models for each screen; Update() delegates by view |
| HTTP client timeout | Phase 1: API client setup | Shared `http.Client` with timeout; no use of `http.DefaultClient` |
| Content sanitization | Phase 1: Deal list render | Real RFD data renders cleanly (no `&amp;`, no HTML tags, proper column alignment) |
| Thread lazy loading | Phase 2: Thread detail | Opening a 100+ reply thread loads first page only; "load more" available |
| Search debouncing | Phase 3: Search/filter | Rapid typing doesn't fire multiple API calls; search triggers on Enter |
| Pagination bounds | Phase 1: Deal list render | "Next page" disabled on last page; page counter shows current/total |
| Error/loading states | Phase 1: First API call | Disconnect network while app runs — spinner shown, error displayed, retry available |
| Browser URL opening | Phase 2: Thread detail | Press 'o' opens deal URL; shows message if no URL; doesn't crash in headless environments |
| Empty state rendering | Every feature phase | Every list/search result handles zero results gracefully |

## Sources

- Bubble Tea v2 source code and upgrade guide: `github.com/charmbracelet/bubbletea/blob/main/UPGRADE_GUIDE_V2.md` — HIGH confidence (official)
- Charmbracelet ecosystem HTTP client patterns (Crush, Catwalk, Gum): 30s timeout standard across all projects — HIGH confidence (official projects)
- Grafana Loki bench TUI (`pkg/logql/bench/cmd/bench/`): ViewID + sub-model decomposition pattern — HIGH confidence (production code)
- GitHub CLI (`cli/cli`): Uses `charm.land/bubbletea/v2` with `huh/v2`, `bubbles/v2`, `lipgloss/v2` — HIGH confidence (production code)
- Charmbracelet Bubble Tea examples: WindowSizeMsg handling, tea.Cmd patterns — HIGH confidence (official examples)
- RFD JSON API behavior: Based on PROJECT.md reference project (davegallant/rfd, archived Aug 2024) — MEDIUM confidence (API observed working April 2026 but not guaranteed stable)
- Browser URL opening: `pkg/browser` standard library — HIGH confidence (stdlib)
- Lipgloss v2 layout and measurement: Context7 docs for `charm.land/lipgloss/v2` — HIGH confidence (official docs)

---
*Pitfalls research for: Go TUI for RedFlagDeals*
*Researched: 2026-04-17*
