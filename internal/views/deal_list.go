package views

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/simon/rfdtui/internal/client"
	"github.com/simon/rfdtui/internal/styles"
	"github.com/simon/rfdtui/internal/types"
)

type sortMode int

const (
	sortNone sortMode = iota
	sortScore
	sortViews
)

var sortLabels = map[sortMode]string{
	sortNone:  "Default",
	sortScore: "Score",
	sortViews: "Views",
}

var minScoreOptions = []int{0, 1, 5, 10, 25, 50}

type DealListModel struct {
	topics         []types.Topic
	filteredTopics []types.Topic
	cursor         int
	page           int
	loading        bool
	err            error
	spinner        spinner.Model
	width          int
	height         int
	client         *client.Client

	searchActive bool
	searchInput  textinput.Model
	searchQuery  string
	isRegex      bool

	minScoreIdx int
	sortBy      sortMode
}

func NewDealList(c *client.Client) DealListModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.LoadingStyle

	ti := textinput.New()
	ti.Prompt = "/"
	ti.Placeholder = "search deals..."
	ti.CharLimit = 100

	return DealListModel{
		client:         c,
		page:           1,
		spinner:        s,
		searchInput:    ti,
		filteredTopics: nil,
	}
}

func (m DealListModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.client.FetchTopics(1))
}

func (m DealListModel) Update(msg tea.Msg) (DealListModel, tea.Cmd) {
	if m.searchActive {
		return m.updateSearch(msg)
	}
	return m.updateNormal(msg)
}

func (m DealListModel) updateSearch(msg tea.Msg) (DealListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			m.searchActive = false
			m.searchQuery = m.searchInput.Value()
			m.searchInput.Blur()
			m.applyFilters()
			m.cursor = 0
			return m, nil
		case "esc":
			m.searchActive = false
			m.searchInput.SetValue("")
			m.searchInput.Blur()
			m.searchQuery = ""
			m.applyFilters()
			m.cursor = 0
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	return m, cmd
}

func (m DealListModel) updateNormal(msg tea.Msg) (DealListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case client.TopicsMsg:
		m.topics = msg.Topics
		m.page = msg.Page
		m.loading = false
		m.err = nil
		m.applyFilters()
		if m.cursor >= len(m.filteredTopics) {
			m.cursor = len(m.filteredTopics) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}
		return m, nil

	case client.ErrMsg:
		m.loading = false
		m.err = msg.Err
		return m, nil

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.filteredTopics)-1 {
				m.cursor++
			}
		case "n":
			if !m.loading {
				m.loading = true
				m.page++
				return m, tea.Batch(m.spinner.Tick, m.client.FetchTopics(m.page))
			}
		case "p":
			if !m.loading && m.page > 1 {
				m.loading = true
				m.page--
				return m, tea.Batch(m.spinner.Tick, m.client.FetchTopics(m.page))
			}
		case "/":
			m.searchActive = true
			m.searchInput.SetValue(m.searchQuery)
			w := m.width - 10
			if w <= 0 {
				w = 30
			}
			m.searchInput.SetWidth(w)
			return m, m.searchInput.Focus()
		case "s":
			m.sortBy = (m.sortBy + 1) % 3
			m.applyFilters()
			m.cursor = 0
		case "f":
			m.minScoreIdx = (m.minScoreIdx + 1) % len(minScoreOptions)
			m.applyFilters()
			m.cursor = 0
		}
	}

	return m, nil
}

func (m *DealListModel) applyFilters() {
	topics := m.topics

	if m.searchQuery != "" {
		topics = m.filterBySearch(topics)
	}

	minScore := minScoreOptions[m.minScoreIdx]
	if minScore > 0 {
		filtered := make([]types.Topic, 0, len(topics))
		for _, t := range topics {
			if t.Score >= minScore {
				filtered = append(filtered, t)
			}
		}
		topics = filtered
	}

	switch m.sortBy {
	case sortScore:
		sort.SliceStable(topics, func(i, j int) bool {
			return topics[i].Score > topics[j].Score
		})
	case sortViews:
		sort.SliceStable(topics, func(i, j int) bool {
			return topics[i].TotalViews > topics[j].TotalViews
		})
	}

	m.filteredTopics = topics
}

func (m *DealListModel) filterBySearch(topics []types.Topic) []types.Topic {
	query := m.searchQuery
	m.isRegex = false

	re, err := regexp.Compile("(?i)" + query)
	if err == nil {
		m.isRegex = true
		return m.filterWithRegex(topics, re)
	}

	lowerQ := strings.ToLower(query)
	filtered := make([]types.Topic, 0, len(topics))
	for _, t := range topics {
		if strings.Contains(strings.ToLower(t.Title), lowerQ) ||
			strings.Contains(strings.ToLower(t.DealerName()), lowerQ) {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

func (m *DealListModel) filterWithRegex(topics []types.Topic, re *regexp.Regexp) []types.Topic {
	filtered := make([]types.Topic, 0, len(topics))
	for _, t := range topics {
		if re.MatchString(t.Title) || re.MatchString(t.DealerName()) {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

func (m DealListModel) View() tea.View {
	if m.loading {
		var b strings.Builder
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("  %s Loading deals...\n\n", m.spinner.View()))
		for i := 0; i < 5; i++ {
			b.WriteString("  ")
			b.WriteString(renderShimmer(m.width - 4))
			b.WriteString("\n")
		}
		return tea.NewView(b.String())
	}

	if m.err != nil {
		return tea.NewView(fmt.Sprintf("\n  %s Error: %v\n\n  Press q to quit.", styles.ErrorStyle.Render("✗"), m.err))
	}

	topics := m.filteredTopics
	if topics == nil {
		topics = m.topics
	}

	if len(topics) == 0 {
		return tea.NewView("\n  No deals found.")
	}

	var b strings.Builder

	header := styles.TitleStyle.Render("RedFlagDeals — Hot Deals")
	b.WriteString(header)
	b.WriteString("\n\n")

	linesPerCard := 3
	visibleHeight := (m.height - 6) / linesPerCard
	if visibleHeight < 1 {
		visibleHeight = 10
	}

	start := 0
	end := len(topics)

	if end > visibleHeight {
		if m.cursor > visibleHeight/2 {
			start = m.cursor - visibleHeight/2
		}
		if start+visibleHeight < end {
			end = start + visibleHeight
		} else {
			end = len(topics)
			start = end - visibleHeight
			if start < 0 {
				start = 0
			}
		}
	}

	for i := start; i < end; i++ {
		t := topics[i]
		titleLine, metaLine := formatDealCard(t, m.width)

		num := styles.PostHeaderStyle.Render(fmt.Sprintf("%2d.", i+1))
		titleLine = num + " " + titleLine

		if i == m.cursor {
			titleLine = styles.SelectedStyle.Render("▶ " + titleLine)
			metaLine = styles.SelectedStyle.Render("    " + metaLine)
		} else {
			titleLine = "  " + titleLine
			metaLine = "    " + metaLine
		}

		b.WriteString(titleLine)
		b.WriteString("\n")
		b.WriteString(metaLine)
		b.WriteString("\n")

		if i < end-1 {
			b.WriteString(styles.DividerStyle.Render("  · · · · · · · · · · · · · · · · · · · · · · · · · · · · ·"))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")

	if m.searchActive {
		b.WriteString("  ")
		b.WriteString(m.searchInput.View())
		b.WriteString("\n")
	} else {
		b.WriteString(m.renderStatusBar(topics))
		b.WriteString("\n")

		scroll := renderScrollBar(m.cursor, len(topics), m.height)
		help := " j/k:nav  Enter:open  /:search  s:sort  f:filter  n/p:page  o:browser  c:copy  ?:help  q:quit"
		b.WriteString(styles.HelpStyle.Render(help))
		b.WriteString(scroll)
	}

	return tea.NewView(b.String())
}

func (m DealListModel) renderStatusBar(topics []types.Topic) string {
	var filterParts []string
	if m.searchQuery != "" {
		mode := "text"
		if m.isRegex {
			mode = "regex"
		}
		filterParts = append(filterParts, fmt.Sprintf("search(%s): \"%s\"", mode, m.searchQuery))
	}
	if minScoreOptions[m.minScoreIdx] > 0 {
		filterParts = append(filterParts, fmt.Sprintf("min score: %d", minScoreOptions[m.minScoreIdx]))
	}
	if m.sortBy != sortNone {
		filterParts = append(filterParts, fmt.Sprintf("sort: %s", sortLabels[m.sortBy]))
	}

	var filterStr string
	if len(filterParts) > 0 {
		filterStr = " │ " + strings.Join(filterParts, " · ")
	}

	status := fmt.Sprintf("Page %d · %d deals%s", m.page, len(topics), filterStr)
	return styles.StatusBarStyle.Render(status)
}

func (m DealListModel) SearchActive() bool {
	return m.searchActive
}

func (m DealListModel) SelectedTopic() *types.Topic {
	topics := m.filteredTopics
	if topics == nil {
		topics = m.topics
	}
	if m.cursor >= 0 && m.cursor < len(topics) {
		return &topics[m.cursor]
	}
	return nil
}

func formatDealCard(t types.Topic, width int) (string, string) {
	title := t.Title
	maxTitleWidth := width - 6
	if maxTitleWidth > 0 && len(title) > maxTitleWidth {
		title = title[:maxTitleWidth-1] + "…"
	}

	var titleParts []string
	if t.Score >= 25 {
		titleParts = append(titleParts, styles.HotBadgeStyle.Render("HOT"))
	}
	titleParts = append(titleParts, styles.TitleLineStyle.Render(title))
	titleLine := strings.Join(titleParts, " ")

	scoreStr := fmt.Sprintf("%+d", t.Score)
	scoreBadge := styles.ScoreBadge(t.Score).Render(scoreStr)

	var metaParts []string
	metaParts = append(metaParts, scoreBadge)

	if catName := t.CategoryName(); catName != "" {
		catStyle, ok := styles.CategoryTagStyles[t.CategoryID()]
		if ok {
			metaParts = append(metaParts, catStyle.Render(catName))
		}
	}

	if t.DealerName() != "" {
		metaParts = append(metaParts, styles.DealerStyle.Render(t.DealerName()))
	}

	if t.Price() != "" {
		priceText := t.Price()
		if t.Savings() != "" {
			priceText = fmt.Sprintf("%s (%s off)", t.Price(), t.Savings())
		}
		metaParts = append(metaParts, styles.PriceBadgeStyle.Render(priceText))
	}

	views := styles.ViewsStyle.Render(fmt.Sprintf("%d views", t.TotalViews))
	replies := styles.RepliesStyle.Render(fmt.Sprintf("%d replies", t.TotalReplies))
	age := styles.AgeStyle.Render(relativeAge(t.PostTime))

	metaParts = append(metaParts, views, replies, age)

	sep := " " + styles.SeparatorStyle.Render("·") + " "
	metaLine := strings.Join(metaParts, sep)

	return titleLine, metaLine
}

func relativeAge(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 48*time.Hour:
		return "1d ago"
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}

func renderScrollBar(cursor, total, height int) string {
	if total <= 0 {
		return ""
	}
	trackLen := 10
	pos := (cursor * trackLen) / total
	if pos >= trackLen {
		pos = trackLen - 1
	}

	track := make([]rune, trackLen)
	for i := range track {
		track[i] = '·'
	}
	track[pos] = '▌'

	return "  " + styles.ScrollTrackStyle.Render(string(track))
}

func renderShimmer(width int) string {
	if width <= 0 {
		width = 60
	}
	line1 := strings.Repeat("━", width*2/3) + strings.Repeat(" ", width-width*2/3)
	line2 := strings.Repeat("━", width/3) + strings.Repeat(" ", width-width/3)
	line3 := strings.Repeat("━", width/2) + strings.Repeat(" ", width-width/2)
	return styles.ShimmerStyle.Render(line1) + "\n" +
		styles.ShimmerHighlightStyle.Render(line2) + "\n" +
		styles.ShimmerStyle.Render(line3)
}
