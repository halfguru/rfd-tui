package views

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/simon/rfd/internal/client"
	"github.com/simon/rfd/internal/styles"
	"github.com/simon/rfd/internal/types"
)

type DealListModel struct {
	topics  []types.Topic
	cursor  int
	page    int
	loading bool
	err     error
	spinner spinner.Model
	width   int
	height  int
	client  *client.Client
}

func NewDealList(c *client.Client) DealListModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.LoadingStyle

	return DealListModel{
		client:  c,
		page:    1,
		spinner: s,
	}
}

func (m DealListModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.client.FetchTopics(1))
}

func (m DealListModel) Update(msg tea.Msg) (DealListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case client.TopicsMsg:
		m.topics = msg.Topics
		m.page = msg.Page
		m.loading = false
		m.err = nil
		if m.cursor >= len(m.topics) {
			m.cursor = len(m.topics) - 1
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
			if m.cursor < len(m.topics)-1 {
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
		}
	}

	return m, nil
}

func (m DealListModel) View() tea.View {
	if m.loading {
		return tea.NewView(fmt.Sprintf("\n  %s Loading deals...", m.spinner.View()))
	}

	if m.err != nil {
		return tea.NewView(fmt.Sprintf("\n  %s Error: %v\n\n  Press q to quit.", styles.ErrorStyle.Render("✗"), m.err))
	}

	if len(m.topics) == 0 {
		return tea.NewView("\n  No deals found.")
	}

	var b strings.Builder

	header := styles.TitleStyle.Render("🔥 RedFlagDeals — Hot Deals")
	b.WriteString(header)
	b.WriteString("\n\n")

	visibleHeight := m.height - 4
	if visibleHeight < 1 {
		visibleHeight = 20
	}

	start := 0
	end := len(m.topics)

	if end > visibleHeight {
		if m.cursor > visibleHeight/2 {
			start = m.cursor - visibleHeight/2
		}
		if start+visibleHeight < end {
			end = start + visibleHeight
		} else {
			end = len(m.topics)
			start = end - visibleHeight
			if start < 0 {
				start = 0
			}
		}
	}

	for i := start; i < end; i++ {
		t := m.topics[i]
		line := formatDealLine(t, m.width)

		if i == m.cursor {
			line = styles.SelectedStyle.Render("▶ " + line)
		} else {
			line = "  " + line
		}

		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	status := fmt.Sprintf(" Page %d | j/k: navigate  Enter: open  n/p: page  o: browser  q: quit", m.page)
	b.WriteString(styles.StatusBarStyle.Render(status))

	return tea.NewView(b.String())
}

func (m DealListModel) SelectedTopic() *types.Topic {
	if m.cursor >= 0 && m.cursor < len(m.topics) {
		return &m.topics[m.cursor]
	}
	return nil
}

func formatDealLine(t types.Topic, width int) string {
	scoreStr := fmt.Sprintf("%+d", t.Score)
	score := styles.ScoreStyle(t.Score).Render(scoreStr)

	dealer := ""
	if t.DealerName() != "" {
		dealer = styles.DealerStyle.Render(t.DealerName())
	}

	views := styles.ViewsStyle.Render(fmt.Sprintf("%d views", t.TotalViews))
	replies := styles.RepliesStyle.Render(fmt.Sprintf("%d replies", t.TotalReplies))
	age := styles.AgeStyle.Render(relativeAge(t.PostTime))

	price := ""
	if t.Price() != "" {
		price = styles.PriceStyle.Render(t.Price())
	}
	if t.Savings() != "" && price != "" {
		price += styles.PriceStyle.Render(fmt.Sprintf(" (%s off)", t.Savings()))
	}

	title := t.Title
	if width > 0 {
		maxTitleWidth := width - 40
		if maxTitleWidth > 0 && len(title) > maxTitleWidth {
			title = title[:maxTitleWidth-1] + "…"
		}
	}

	parts := []string{score, title}
	if dealer != "" {
		parts = append(parts, dealer)
	}
	if price != "" {
		parts = append(parts, price)
	}
	parts = append(parts, views, replies, age)

	return strings.Join(parts, "  ")
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
