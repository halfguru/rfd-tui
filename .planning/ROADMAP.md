# Roadmap: RFD TUI

## Overview

Build a terminal UI for browsing RedFlagDeals.com hot deals in three phases: first, establish the API client and deliver the core deal list with vim navigation; then add thread reading and browser integration; finally, layer on search, filtering, sorting, and a help overlay to round out the v1 experience.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [ ] **Phase 1: Foundation & Deal List** - Fetch and display hot deals in a scrollable, navigable TUI
- [ ] **Phase 2: Thread Detail & Browser** - Read deal threads and open deals in the system browser
- [ ] **Phase 3: Search, Filter & Polish** - Find specific deals and discover all keybindings

## Phase Details

### Phase 1: Foundation & Deal List
**Goal**: Users launch the app into a responsive, scrollable list of current hot deals from RedFlagDeals
**Depends on**: Nothing (first phase)
**Requirements**: API-01, API-03, API-04, LIST-01, LIST-02, LIST-03, LIST-04, LIST-05, UX-02
**Success Criteria** (what must be TRUE):
  1. App fetches deals on launch and displays a scrollable list with title, score (color-coded), views, replies, dealer, category, and relative age per deal
  2. User can navigate with j/k/arrow keys, open a deal with Enter, paginate with n/p, and quit with q/Ctrl+C
  3. App shows a loading spinner during API fetches and displays user-friendly error messages for network failures or malformed responses
  4. App adapts its layout when the terminal is resized
**Plans**: TBD

### Phase 2: Thread Detail & Browser
**Goal**: Users can read full deal threads and open deals in their web browser
**Depends on**: Phase 1
**Requirements**: API-02, THRD-01, THRD-02, THRD-03, THRD-04, BROW-01, BROW-02
**Success Criteria** (what must be TRUE):
  1. User can open a deal to view the thread detail showing the original post with HTML stripped to plain text
  2. User can scroll down through a thread to lazy-load additional reply pages from the API
  3. User can collapse and expand comment threads with the tab key for easier navigation
  4. User can press 'o' to open the deal's URL in the system web browser, with graceful handling of deals with no external link
  5. User can return to the deal list from thread detail with Escape/q
**Plans**: TBD
**UI hint**: yes

### Phase 3: Search, Filter & Polish
**Goal**: Users can find specific deals through search and filtering, and discover all keybindings via help
**Depends on**: Phase 2
**Requirements**: SRCH-01, SRCH-02, SRCH-03, SRCH-04, SRCH-05, UX-01
**Success Criteria** (what must be TRUE):
  1. User can search deals by keyword or regex across titles and dealer names via an inline search input
  2. User can filter the deal list by category and set a minimum score threshold to hide low-quality deals
  3. User can sort deals by score or view count with a client-side sort toggle
  4. User can press '?' to see a help overlay listing all available keybindings for the current view
**Plans**: TBD
**UI hint**: yes

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation & Deal List | 0/? | Not started | - |
| 2. Thread Detail & Browser | 0/? | Not started | - |
| 3. Search, Filter & Polish | 0/? | Not started | - |
