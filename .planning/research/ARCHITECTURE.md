# Architecture Research

**Domain:** Bubble Tea TUI application consuming a REST JSON API
**Researched:** 2026-04-17
**Confidence:** HIGH

## Standard Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         TUI Layer (Bubble Tea)                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                   Root Model (app.go)                     │   │
│  │  ┌────────────┐  ┌────────────┐  ┌───────────────────┐  │   │
│  │  │  Deal List │  │   Thread   │  │   Search/Filter   │  │   │
│  │  │  (bubbles  │  │  Detail    │  │    (textinput +   │  │   │
│  │  │   /list)   │  │ (viewport) │  │    custom logic)  │  │   │
│  │  └─────┬──────┘  └─────┬──────┘  └────────┬──────────┘  │   │
│  └────────┴────────────────┴──────────────────┴─────────────┘   │
│         │                 │                    │                  │
│    ┌────┴─────────────────┴────────────────────┴──────┐         │
│    │              tea.Cmd (async messages)             │         │
│    └────┬─────────────────┬────────────────────┬──────┘         │
├─────────┼─────────────────┼────────────────────┼────────────────┤
│         │           API Client Layer            │                │
│  ┌──────┴──────────────────────────────────────┴───────────┐    │
│  │                    rfd.Client                            │    │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │    │
│  │  │ FetchTopics  │  │ FetchPosts   │  │  SearchTopics │  │    │
│  │  │ (deals list) │  │ (thread)     │  │  (keyword)    │  │    │
│  │  └──────────────┘  └──────────────┘  └──────────────┘  │    │
│  └─────────────────────────┬───────────────────────────────┘    │
│                            │                                     │
├────────────────────────────┼─────────────────────────────────────┤
│                       HTTP Layer                                  │
│  ┌─────────────────────────┴───────────────────────────────┐    │
│  │                   net/http.Client                         │   │
│  │          https://forums.redflagdeals.com/api/             │   │
│  └──────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
```

### Component Responsibilities

| Component | Responsibility | Implementation |
|-----------|----------------|----------------|
| **Root Model** | Application state, view routing, key dispatch | Single `model` struct implementing `tea.Model` |
| **Deal List** | Browse, scroll, select deals | `bubbles/list.Model` with custom delegate |
| **Thread Detail** | Display posts, scrollable content | `bubbles/viewport.Model` |
| **Search/Filter** | Keyword/regex input, filter results | `bubbles/textinput.Model` + custom filter logic |
| **rfd.Client** | HTTP requests, JSON parsing, error wrapping | Plain Go struct with `http.Client` |
| **tea.Cmd** | Bridge between UI and API — async I/O | Functions returning `tea.Msg` types |
| **Styles** | Consistent visual presentation | `lipgloss.Style` definitions in dedicated file |

## Recommended Project Structure

```
rfd/
├── main.go                    # Entry point, wire dependencies, run program
├── app.go                     # Root model: state, Init/Update/View, view routing
├── keys.go                    # Key bindings (key.Map definition)
├── styles.go                  # Lipgloss style constants
├── views/
│   ├── deals.go               # Deal list model: list setup, custom delegate, rendering
│   ├── thread.go              # Thread detail model: viewport, post rendering
│   └── help.go                # Help overlay rendering
├── rfd/
│   ├── client.go              # HTTP client: FetchTopics, FetchPosts, SearchTopics
│   ├── types.go               # API response structs (Topic, Post, User, Pager, etc.)
│   └── client_test.go         # Client tests with httptest server
├── msg/
│   └── messages.go            # Custom tea.Msg types (dealsLoadedMsg, postsLoadedMsg, etc.)
├── go.mod
├── go.sum
└── .goreleaser.yml            # Release config for single binary
```

### Structure Rationale

- **`app.go` as root model**: Standard Bubble Tea convention. The root model holds `activeView`, window dimensions, sub-models (list, viewport, textinput), and the API client. It routes `Update` and `View` calls based on `activeView`.

- **`views/` package for UI sub-models**: Each view (deals list, thread detail) has its own file with setup helpers, rendering functions, and view-specific key handling. These aren't separate `tea.Model` implementations — they're helper structs/functions called by the root model. This keeps the root model from becoming a 1000-line file.

- **`rfd/` package for API concerns**: Complete isolation of HTTP logic. The TUI layer never sees `net/http` — it calls `client.FetchTopics()` and receives Go structs or errors. Testable with `httptest.NewServer`.

- **`msg/` package for message types**: Custom `tea.Msg` types that the root model switches on. Separating them makes the Update function's type switch clearer and avoids circular imports.

## Architectural Patterns

### Pattern 1: Elm Architecture (Model-View-Update)

**What:** The core Bubble Tea pattern. All application state lives in a single model struct. `Update` receives messages and returns a new model + optional commands. `View` renders the model as a string. No mutation outside Update.

**When to use:** Always — this IS Bubble Tea. Every Bubble Tea app uses this.

**Trade-offs:** Takes discipline to keep all state in the model (no globals). Advantage is complete testability — you can test Update by creating a model, passing a message, and asserting the returned model.

```go
type model struct {
    activeView  viewMode       // which view is active
    dealsList   list.Model     // bubbles list for deals
    threadVP    viewport.Model // viewport for thread detail
    filterInput textinput.Model // search input
    client      *rfd.Client    // API client (no tea.Model)
    deals       []rfd.Topic    // raw deal data
    page        int            // current pagination page
    totalPages  int            // from API pager
    width       int
    height      int
    loading     bool
    err         error
    quitting    bool
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
        m.dealsList.SetSize(msg.Width, msg.Height-4)
    case dealsLoadedMsg:
        m.loading = false
        if msg.err != nil {
            m.err = msg.err
            return m, nil
        }
        m.deals = msg.deals
        m.totalPages = msg.totalPages
        // Convert deals to list items
        items := make([]list.Item, len(m.deals))
        for i, d := range m.deals {
            items[i] = dealItem{deal: d}
        }
        m.dealsList.SetItems(items)
    case tea.KeyPressMsg:
        // route to view-specific handler
        return m.handleKey(msg)
    }
    // Forward to active sub-component
    var cmd tea.Cmd
    m.dealsList, cmd = m.dealsList.Update(msg)
    return m, cmd
}
```

### Pattern 2: Active View Routing

**What:** The root model holds an `activeView` enum. `Update` and `View` delegate to different logic based on which view is active. Views share the same root model but render and handle keys differently.

**When to use:** Any TUI with more than one screen. Verified across 7+ production Bubble Tea apps (yt-browse, gocovsh, wasteland, hey-cli, MastodonCLI, pproftui, ez-monitor).

**Trade-offs:** Root model can grow large. Mitigate by extracting view-specific logic into helper methods or the `views/` package. Do NOT create separate `tea.Model` implementations for each view — the wiring complexity isn't worth it for 2-3 views.

```go
type viewMode int

const (
    viewDeals viewMode = iota
    viewThread
    viewSearch
)

func (m model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
    switch m.activeView {
    case viewDeals:
        return m.handleDealsKey(msg)
    case viewThread:
        return m.handleThreadKey(msg)
    case viewSearch:
        return m.handleSearchKey(msg)
    }
    return m, nil
}

func (m model) View() string {
    switch m.activeView {
    case viewDeals:
        return m.viewDeals()
    case viewThread:
        return m.viewThread()
    case viewSearch:
        return m.viewSearch()
    }
    return ""
}
```

### Pattern 3: Command-Based Async I/O

**What:** All HTTP requests are wrapped in `tea.Cmd` functions. A command is a function `func() tea.Msg` that runs in a goroutine. When it completes, the returned message is fed back into `Update`. The UI never blocks.

**When to use:** Every network call, file read, or slow operation. Bubble Tea provides no other async mechanism.

**Trade-offs:** Commands run in goroutines — they can't directly modify the model. They must return a message. This is correct by design: state changes only happen in Update.

```go
// In msg/messages.go
type dealsLoadedMsg struct {
    deals       []rfd.Topic
    totalPages  int
    page        int
    err         error
}

// Command factory — captures parameters in a closure
func fetchDeals(client *rfd.Client, forumID, page int) tea.Cmd {
    return func() tea.Msg {
        topics, totalPages, err := client.FetchTopics(forumID, 40, page)
        return dealsLoadedMsg{
            deals:      topics,
            totalPages: totalPages,
            page:       page,
            err:        err,
        }
    }
}

// In Update — kick off the command
func (m model) Init() tea.Cmd {
    return fetchDeals(m.client, 9, 1) // forum_id=9 = Hot Deals
}

// In Update — handle the result
case dealsLoadedMsg:
    m.loading = false
    if msg.err != nil {
        m.err = msg.err
        return m, nil
    }
    m.deals = msg.deals
    m.totalPages = msg.totalPages
    m.page = msg.page
    // ... update list items
```

### Pattern 4: Embedded Sub-Models (NOT Composition)

**What:** Bubble Tea sub-components (`list.Model`, `viewport.Model`, `textinput.Model`) are embedded as fields in the root model. The root model's `Update` forwards messages to the active sub-component and captures its returned command.

**When to use:** Always for Bubble Tea built-in components.

**Trade-offs:** The root model's Update becomes a router. Keep it clean with helper methods. DO NOT try to make each view a separate `tea.Model` — the wiring (Init/Update/View delegation, command forwarding, state sharing) creates more complexity than it solves for 2-3 views.

```go
// Forward messages to the active list
var cmd tea.Cmd
switch m.activeView {
case viewDeals:
    m.dealsList, cmd = m.dealsList.Update(msg)
case viewThread:
    m.threadVP, cmd = m.threadVP.Update(msg)
}
return m, cmd
```

## Data Flow

### Deal Browsing Flow

```
User launches app
    ↓
Init() → fetchDeals(client, forumID=9, page=1)     [tea.Cmd — goroutine]
    ↓
API Client → HTTP GET /api/topics?forum_id=9&per_page=40&page=1
    ↓
Parse JSON → []rfd.Topic, totalPages                [in goroutine]
    ↓
Return dealsLoadedMsg{deals, totalPages}            [tea.Msg]
    ↓
Update() → m.deals = msg.deals                      [main goroutine]
         → Convert to []list.Item, set on dealsList
         → m.loading = false
    ↓
View() → Render dealsList.View() + status bar + help
    ↓
User presses j/k → dealsList handles scroll internally
User presses Enter → m.openSelectedDeal()
    ↓
openSelectedDeal() → fetchPosts(client, topicID)
    ↓
API Client → HTTP GET /api/topics/{id}/posts?per_page=40&page=1
    ↓
Return postsLoadedMsg{posts}
    ↓
Update() → m.activeView = viewThread
         → m.threadVP.SetContent(renderPosts(msg.posts))
    ↓
View() → Render threadVP.View()
```

### Search Flow

```
User presses / → m.activeView = viewSearch, focus textinput
    ↓
User types query → textinput.Update(msg)
User presses Enter → applyFilter()
    ↓
applyFilter() → fetchDeals(client, forumID, page) with search params
              OR filter m.deals in-memory (for already-loaded pages)
    ↓
dealsLoadedMsg → Update list items with filtered results
```

### State Management

```
┌─────────────────────────────────────────────┐
│                  model (single source of truth)
│
│  Navigation State:   activeView, page, totalPages
│  Deal Data:          deals []rfd.Topic, filteredDeals
│  Thread Data:        posts []rfd.Post
│  UI Components:      dealsList, threadVP, filterInput
│  Window:             width, height
│  Status:             loading, err, quitting
│  External:           client *rfd.Client
│
│  No goroutines share state. All mutation in Update().
└─────────────────────────────────────────────┘
```

### Key Data Flows

1. **Deals Fetching:** `tea.Cmd` (goroutine) → `rfd.Client.FetchTopics()` → `dealsLoadedMsg` → `Update()` sets model fields → `View()` renders list
2. **Thread Viewing:** User selects deal → `tea.Cmd` → `rfd.Client.FetchPosts()` → `postsLoadedMsg` → `Update()` switches view + sets viewport content → `View()` renders viewport
3. **Pagination:** User presses n/p → increment/decrement page → `tea.Cmd` → `rfd.Client.FetchTopics(page)` → same flow as #1
4. **Open in Browser:** User presses o → `tea.ExecProcess(exec.Command("open", dealURL))` → browser opens (no state change)
5. **Search/Filter:** User types → filter loaded deals in-memory (or new API call for regex/keyword across pages) → update list items

## Scaling Considerations

This is a single-user CLI tool. Scaling here means code organization, not load.

| Scale | Architecture Adjustments |
|-------|--------------------------|
| 2-3 views (MVP) | Root model with activeView routing. All state in one struct. Views as helper methods. |
| 5+ views | Extract views into separate files in `views/`. Consider if any view is complex enough to warrant its own sub-model (but probably not). |
| Many API endpoints | Client grows — keep methods on one struct, group by domain (topics, posts, search). |

### Scaling Priorities

1. **First bottleneck: Root model size.** Mitigate early by extracting rendering logic into `views/` package functions that accept the model (or relevant subset) as a parameter.
2. **Second bottleneck: List item rendering.** Custom delegates can get complex. Keep the delegate simple — format the string in a helper, don't compute in Render().

## Anti-Patterns

### Anti-Pattern 1: Separate tea.Model per View

**What people do:** Create `DealsModel`, `ThreadModel`, etc. each implementing `tea.Model`. Try to wire them together with a parent model.

**Why it's wrong:** Bubble Tea doesn't natively support component trees. You end up manually routing Init/Update/View calls, forwarding commands, sharing state via pointers, and creating fragile wiring. The Charm team's own examples don't do this.

**Do this instead:** One root model with `activeView` routing. Extract complex view logic into helper functions or a `views/` package, but keep state centralized.

### Anti-Pattern 2: Blocking HTTP in Update

**What people do:** Call `http.Get()` directly in the Update function.

**Why it's wrong:** Update runs on the main goroutine. Blocking it freezes the entire TUI — no key handling, no rendering, no spinner animation.

**Do this instead:** Always wrap HTTP calls in a `tea.Cmd` (closure that returns `tea.Msg`). Commands run in goroutines automatically.

### Anti-Pattern 3: Global State or Package-Level Variables

**What people do:** Use package-level variables for configuration, API client, or shared state.

**Why it's wrong:** Breaks testability. Makes it impossible to run multiple instances. Bubble Tea's Elm architecture works because all state flows through the model.

**Do this instead:** Pass dependencies (API client, config) into the model constructor. Store them as model fields.

### Anti-Pattern 4: Ignoring WindowSizeMsg

**What people do:** Don't handle `tea.WindowSizeMsg`, use hardcoded widths.

**Why it's wrong:** Terminals vary widely. Your TUI will break on different terminal sizes, tmux panes, or when the user resizes.

**Do this instead:** Always handle `WindowSizeMsg`. Store width/height on the model. Pass them to sub-components. Use `lipgloss.NewStyle().Width(m.width)` for responsive layouts.

## Integration Points

### External Services

| Service | Integration Pattern | Notes |
|---------|---------------------|-------|
| RFD JSON API | `rfd.Client` → `http.Client` → GET requests | Base URL: `https://forums.redflagdeals.com`. No auth required. Rate-limit conservatively (sleep between paginated requests). |
| System Browser | `tea.ExecProcess` with `exec.Command("open" or "xdg-open")` | Cross-platform: `open` (macOS), `xdg-open` (Linux), `start` (Windows). Use `exec.LookPath` to find the right one. |

### Internal Boundaries

| Boundary | Communication | Notes |
|----------|---------------|-------|
| Root Model ↔ rfd.Client | Client methods called in `tea.Cmd` closures | Client is injected into model. Never call client methods directly in Update — always wrap in Cmd. |
| Root Model ↔ Sub-components (list, viewport) | Forward `tea.Msg` via `.Update(msg)` | Root model's Update must forward messages to the active sub-component. Non-active components don't get messages. |
| views/ helpers ↔ Root model | Functions accept data, return strings | View helpers are pure functions. They don't hold state. They take the model (or relevant data) and return rendered strings. |

## Build Order Implications

The component dependencies suggest this build order:

```
1. rfd/types.go         — API response structs (zero deps, needed by everything)
2. rfd/client.go        — HTTP client (depends on types.go)
3. msg/messages.go      — tea.Msg types (depends on rfd types)
4. keys.go              — Key bindings (standalone)
5. styles.go            — Lipgloss styles (standalone)
6. views/deals.go       — Deal list delegate + rendering (depends on types, styles)
7. views/thread.go      — Thread post rendering (depends on types, styles)
8. app.go               — Root model (depends on everything above)
9. main.go              — Wiring + program start (depends on app.go)
```

**Phase 1 can deliver working value after step 8** — the app can fetch and display deals. Search, sort, filter, and thread view are additive layers on the same structure.

## Sources

- **Bubble Tea official documentation** (Context7: `/charmbracelet/bubbletea`) — Model-View-Update, tea.Cmd patterns, multiple views example — HIGH confidence
- **Bubbles component library** (Context7: `/charmbracelet/bubbles`) — list, viewport, paginator, textinput, table components — HIGH confidence
- **Lipgloss styling** (Context7: `/charmbracelet/lipgloss`) — layout, borders, alignment — HIGH confidence
- **yt-browse** (github.com/nolenroyalty/yt-browse) — Production Bubble Tea app with API client + activeView routing + list/viewport/filter pattern. Most similar architecture to what we're building — HIGH confidence
- **gocovsh** (github.com/orlangure/gocovsh) — list + viewport + activeView pattern — HIGH confidence
- **wasteland** (github.com/gastownhall/wasteland) — browse + detail + activeView + settings — HIGH confidence
- **basecamp/hey-cli** — section-based navigation with sub-views — MEDIUM confidence
- **Bubble Tea app template** (Context7: `/charmbracelet/bubbletea-app-template`) — Official project scaffolding — HIGH confidence

---
*Architecture research for: Bubble Tea TUI consuming RFD JSON API*
*Researched: 2026-04-17*
