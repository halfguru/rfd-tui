package views

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/simon/rfd/internal/client"
	"github.com/simon/rfd/internal/styles"
	"github.com/simon/rfd/internal/types"
)

type ThreadModel struct {
	topic      *types.Topic
	posts      []types.Post
	page       int
	totalPages int
	loading    bool
	err        error
	spinner    spinner.Model
	viewport   viewport.Model
	width      int
	height     int
	client     *client.Client
}

func NewThread(topic *types.Topic, c *client.Client) ThreadModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.LoadingStyle

	vp := viewport.New()

	return ThreadModel{
		topic:    topic,
		client:   c,
		page:     1,
		spinner:  s,
		viewport: vp,
		loading:  true,
	}
}

func (m ThreadModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.client.FetchPosts(m.topic.TopicID, 1))
}

func (m ThreadModel) Update(msg tea.Msg) (ThreadModel, tea.Cmd) {
	switch msg := msg.(type) {
	case client.PostsMsg:
		m.posts = append(m.posts, msg.Posts...)
		m.page = msg.Page
		m.totalPages = msg.Pager.TotalPages
		m.loading = false
		m.err = nil
		m.viewport.SetContent(m.renderPosts())
		m.viewport.GotoTop()
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
		m.viewport.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height - 3)
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "escape":
			return m, nil
		case " ":
			if !m.loading && m.page < m.totalPages {
				m.loading = true
				return m, tea.Batch(m.spinner.Tick, m.client.FetchPosts(m.topic.TopicID, m.page+1))
			}
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m ThreadModel) View() tea.View {
	if m.loading && len(m.posts) == 0 {
		return tea.NewView(fmt.Sprintf("\n  %s Loading thread...", m.spinner.View()))
	}

	if m.err != nil && len(m.posts) == 0 {
		return tea.NewView(fmt.Sprintf("\n  %s Error: %v", styles.ErrorStyle.Render("✗"), m.err))
	}

	var b strings.Builder

	header := styles.TitleStyle.Render(m.topic.Title)
	b.WriteString(header)
	b.WriteString("\n")

	info := fmt.Sprintf("  %s  %s  %s",
		styles.DealerStyle.Render(m.topic.DealerName()),
		styles.ViewsStyle.Render(fmt.Sprintf("%d views", m.topic.TotalViews)),
		styles.RepliesStyle.Render(fmt.Sprintf("%d replies", m.topic.TotalViews)),
	)
	b.WriteString(info)
	b.WriteString("\n\n")

	b.WriteString(m.viewport.View())

	if m.page < m.totalPages {
		b.WriteString(styles.HelpStyle.Render(fmt.Sprintf("\n  Space: load more (page %d/%d)", m.page, m.totalPages)))
	}

	b.WriteString(styles.HelpStyle.Render("  Esc/q: back  ↑/↓: scroll"))

	return tea.NewView(b.String())
}

func (m ThreadModel) renderPosts() string {
	var b strings.Builder

	for i, p := range m.posts {
		if i > 0 {
			b.WriteString("\n")
			b.WriteString(strings.Repeat("─", min(m.width, 80)))
			b.WriteString("\n\n")
		}

		voteStr := ""
		if p.Votes.TotalUp > 0 || p.Votes.TotalDown > 0 {
			voteStr = fmt.Sprintf("  ↑%d↓%d", p.Votes.TotalUp, p.Votes.TotalDown)
		}

		header := fmt.Sprintf("#%d · user_%d · %s%s",
			p.Number,
			p.AuthorID,
			relativeAge(p.PostTime),
			voteStr,
		)
		b.WriteString(styles.ViewsStyle.Render(header))
		b.WriteString("\n\n")

		body := stripHTML(p.Body)
		body = wrapText(body, m.width-4)
		b.WriteString(body)
		b.WriteString("\n")
	}

	return b.String()
}

var (
	brRe      = regexp.MustCompile(`<br\s*/?>`)
	htmlTagRe = regexp.MustCompile(`<[^>]+>`)
	quoteRe   = regexp.MustCompile(`&quot;`)
	ampRe     = regexp.MustCompile(`&amp;`)
	nbspRe    = regexp.MustCompile(`&nbsp;`)
)

func stripHTML(s string) string {
	s = brRe.ReplaceAllString(s, "\n")
	s = htmlTagRe.ReplaceAllString(s, "")
	s = quoteRe.ReplaceAllString(s, `"`)
	s = ampRe.ReplaceAllString(s, "&")
	s = nbspRe.ReplaceAllString(s, " ")
	s = html.UnescapeString(s)
	return strings.TrimSpace(s)
}

func wrapText(s string, width int) string {
	if width <= 0 {
		width = 80
	}
	var b strings.Builder
	for _, line := range strings.Split(s, "\n") {
		if len(line) <= width {
			b.WriteString(line)
			b.WriteString("\n")
			continue
		}
		for len(line) > width {
			space := strings.LastIndex(line[:width], " ")
			if space <= 0 {
				space = width
			}
			b.WriteString(line[:space])
			b.WriteString("\n")
			line = line[space:]
			if len(line) > 0 && line[0] == ' ' {
				line = line[1:]
			}
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
