package styles

import (
	"charm.land/lipgloss/v2"
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 2)

	SelectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#333366")).
			Foreground(lipgloss.Color("#FFFFFF"))

	ScorePositive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).Bold(true)

	ScoreNeutral = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).Bold(true)

	ScoreNegative = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).Bold(true)

	DealerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6BB8FF"))

	ViewsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	AgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	PriceStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).Bold(true)

	RepliesStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA"))

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))

	LoadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).Bold(true)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA")).
			Background(lipgloss.Color("#222222"))

	BorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#444444"))
)

func ScoreStyle(score int) lipgloss.Style {
	if score > 5 {
		return ScorePositive
	}
	if score > 0 {
		return ScoreNeutral
	}
	return ScoreNegative
}
