package views

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/simon/rfdtui/internal/client"
	"github.com/simon/rfdtui/internal/styles"
	"github.com/simon/rfdtui/internal/types"
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

func NewThread(topic *types.Topic, c *client.Client, w, h int) ThreadModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.LoadingStyle

	vp := viewport.New()
	vp.SetWidth(w)
	if h > 3 {
		vp.SetHeight(h - 3)
	}

	return ThreadModel{
		topic:    topic,
		client:   c,
		page:     1,
		spinner:  s,
		viewport: vp,
		loading:  true,
		width:    w,
		height:   h,
	}
}

func (m ThreadModel) Topic() *types.Topic {
	return m.topic
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
		cmds := []tea.Cmd{}
		for _, p := range msg.Posts {
			if m.client.CachedUsername(p.AuthorID) == fmt.Sprintf("user_%d", p.AuthorID) {
				cmds = append(cmds, m.client.FetchUsername(p.AuthorID))
			} else {
				for i := range m.posts {
					if m.posts[i].AuthorID == p.AuthorID {
						m.posts[i].AuthorName = m.client.CachedUsername(p.AuthorID)
					}
				}
			}
		}
		if len(cmds) > 0 {
			return m, tea.Batch(cmds...)
		}
		m.viewport.SetContent(m.renderPosts())
		return m, nil

	case client.UsernameMsg:
		for i := range m.posts {
			if m.posts[i].AuthorID == msg.AuthorID {
				m.posts[i].AuthorName = msg.Username
			}
		}
		m.viewport.SetContent(m.renderPosts())
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
		case "esc":
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
		var b strings.Builder
		fmt.Fprintf(&b, "\n  %s Loading thread...\n\n", m.spinner.View())
		for i := 0; i < 5; i++ {
			b.WriteString("  ")
			b.WriteString(renderThreadShimmer(m.width - 4))
			b.WriteString("\n")
		}
		return tea.NewView(b.String())
	}

	if m.err != nil && len(m.posts) == 0 {
		return tea.NewView(fmt.Sprintf("\n  %s Error: %v", styles.ErrorStyle.Render("✗"), m.err))
	}

	var b strings.Builder

	title := styles.TitleStyle.Render(m.topic.Title)
	b.WriteString(title)
	b.WriteString("\n")

	var metaParts []string
	if m.topic.DealerName() != "" {
		metaParts = append(metaParts, styles.DealerStyle.Render(m.topic.DealerName()))
	}
	metaParts = append(metaParts,
		styles.ViewsStyle.Render(fmt.Sprintf("%d views", m.topic.TotalViews)),
		styles.RepliesStyle.Render(fmt.Sprintf("%d replies", m.topic.TotalReplies)),
	)
	sep := " " + styles.SeparatorStyle.Render("·") + " "
	meta := "  " + strings.Join(metaParts, sep)
	b.WriteString(meta)
	b.WriteString("\n")

	divider := strings.Repeat("─", min(m.width, 80))
	b.WriteString(styles.SeparatorStyle.Render(divider))
	b.WriteString("\n")

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
			b.WriteString(styles.DividerStyle.Render("  · · · · · · · · · · · · · · · · · · · · · · · · · · · · ·"))
			b.WriteString("\n")
		}

		header := m.renderPostHeader(p)
		b.WriteString("  ")
		b.WriteString(header)
		b.WriteString("\n")

		body := stripHTML(p.Body)
		body = collapseNewlines(body)
		body = wrapText(body, m.width-4)
		body = strings.TrimRight(body, "\n")

		for _, line := range strings.Split(body, "\n") {
			if strings.TrimSpace(line) == "" {
				b.WriteString("\n")
				continue
			}
			b.WriteString("    ")
			b.WriteString(line)
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (m ThreadModel) renderPostHeader(p types.Post) string {
	numStr := styles.PostHeaderStyle.Render(fmt.Sprintf("%2d.", p.Number))
	name := p.AuthorName
	if name == "" {
		name = fmt.Sprintf("user_%d", p.AuthorID)
	}
	authorStr := styles.PostUserStyle.Render(name)
	ageStr := styles.AgeStyle.Render(relativeAge(p.PostTime))

	parts := []string{numStr, authorStr, ageStr}

	if p.Votes.TotalUp > 0 || p.Votes.TotalDown > 0 {
		upStr := styles.PostUpvoteStyle.Render(fmt.Sprintf("↑%d", p.Votes.TotalUp))
		downStr := styles.PostDownvoteStyle.Render(fmt.Sprintf("↓%d", p.Votes.TotalDown))
		parts = append(parts, upStr, downStr)
	}

	sep := " " + styles.SeparatorStyle.Render("·") + " "
	return strings.Join(parts, sep)
}

var (
	brRe           = regexp.MustCompile(`<br\s*/?>`)
	htmlTagRe      = regexp.MustCompile(`<[^>]+>`)
	quoteRe        = regexp.MustCompile(`&quot;`)
	ampRe          = regexp.MustCompile(`&amp;`)
	nbspRe         = regexp.MustCompile(`&nbsp;`)
	multiNewlineRe = regexp.MustCompile(`\n{3,}`)
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

func collapseNewlines(s string) string {
	return multiNewlineRe.ReplaceAllString(s, "\n\n")
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

func renderThreadShimmer(width int) string {
	if width <= 0 {
		width = 60
	}
	line := strings.Repeat("━", width*2/3) + strings.Repeat(" ", width-width*2/3)
	return styles.ShimmerStyle.Render(line)
}
