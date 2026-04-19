package views

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/simon/rfdtui/internal/styles"
)

type HelpModel struct {
	viewName string
	width    int
	height   int
}

func NewHelp(viewName string) HelpModel {
	return HelpModel{viewName: viewName}
}

func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyPressMsg:
		return m, nil
	}
	return m, nil
}

func (m HelpModel) View() tea.View {
	var b strings.Builder

	bindings := m.getBindings()

	maxKey := 0
	for _, bnd := range bindings {
		if len(bnd.key) > maxKey {
			maxKey = len(bnd.key)
		}
	}

	w := 50
	if m.width > 0 && m.width < 60 {
		w = m.width - 10
	}
	if w < 30 {
		w = 30
	}

	title := styles.TitleStyle.Render(fmt.Sprintf(" Help — %s ", m.viewName))
	b.WriteString(title)
	b.WriteString("\n\n")

	for _, bnd := range bindings {
		keyBadge := styles.KeyBadgeStyle.Render(bnd.key)
		pad := strings.Repeat(" ", maxKey-len(bnd.key)+2)
		desc := styles.DescStyle.Render(bnd.desc)
		fmt.Fprintf(&b, "  %s%s%s\n", keyBadge, pad, desc)
	}

	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("  Press any key to close"))

	content := b.String()
	lines := strings.Split(content, "\n")

	borderW := w
	for _, line := range lines {
		if len(line) > borderW {
			borderW = len(line) + 4
		}
	}

	padded := make([]string, len(lines))
	for i, line := range lines {
		if len(line) < borderW {
			line += strings.Repeat(" ", borderW-len(line))
		}
		padded[i] = line
	}

	borderStyle := styles.BorderStyle.Width(borderW + 2)
	bordered := borderStyle.Render(strings.Join(padded, "\n"))

	return tea.NewView(bordered)
}

type binding struct {
	key  string
	desc string
}

func (m HelpModel) getBindings() []binding {
	switch m.viewName {
	case "Deal List":
		return []binding{
			{"j / ↓", "Move down"},
			{"k / ↑", "Move up"},
			{"Enter", "Open deal thread"},
			{"n", "Next page"},
			{"p", "Previous page"},
			{"/", "Search deals"},
			{"s", "Cycle sort (default/score/views)"},
			{"f", "Cycle min score filter"},
			{"o", "Open in browser"},
			{"?", "Toggle help"},
			{"q", "Quit"},
		}
	case "Thread":
		return []binding{
			{"j / ↓", "Scroll down"},
			{"k / ↑", "Scroll up"},
			{"Space", "Load more replies"},
			{"o", "Open deal in browser"},
			{"?", "Toggle help"},
			{"Esc / q", "Back to deal list"},
		}
	default:
		return []binding{
			{"?", "Toggle help"},
			{"q", "Quit"},
		}
	}
}
