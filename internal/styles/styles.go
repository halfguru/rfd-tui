package styles

import (
	"charm.land/lipgloss/v2"
)

var (
	AppStyle = lipgloss.NewStyle().Padding(0, 2)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#ff6600")).
			Padding(0, 2).
			MarginBottom(1)

	SelectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#2a2520")).
			Foreground(lipgloss.Color("#ffffff"))

	CursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff6600")).Bold(true)

	ScorePositive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2e7d32")).Bold(true)

	ScoreNeutral = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f9a825")).Bold(true)

	ScoreNegative = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#c62828")).Bold(true)

	ScoreBadgePositive = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				Background(lipgloss.Color("#2e7d32")).
				Padding(0, 1).Bold(true)

	ScoreBadgeNeutral = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#f9a825")).
				Padding(0, 1).Bold(true)

	ScoreBadgeNegative = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				Background(lipgloss.Color("#c62828")).
				Padding(0, 1).Bold(true)

	DealerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7fb3d8")).Bold(true)

	ViewsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9a918a"))

	AgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9a918a"))

	PriceStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#66bb6a")).Bold(true)

	PriceBadgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#66bb6a")).Bold(true)

	RepliesStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9a918a"))

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7a716a"))

	LoadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff6600"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#c62828")).Bold(true)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d4cfc9")).
			Background(lipgloss.Color("#2a2520")).
			Padding(0, 2)

	BorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#555555"))

	DividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#444444"))

	PostHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9a918a"))

	PostUserStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e0e0e0")).Bold(true)

	PostBodyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#cccccc"))

	PostBorderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5a534d"))

	PostUpvoteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4caf50")).Bold(true)

	PostDownvoteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ef5350")).Bold(true)

	SearchPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff6600")).Bold(true)

	KeyBadgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#555555")).
			Padding(0, 1)

	DescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9a918a"))

	TitleLineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e8e0d8")).Bold(true)

	MetaLineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9a918a"))

	SeparatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5a534d"))

	HotBadgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#ff4444")).
			Padding(0, 1).Bold(true)

	ScrollBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff6600"))

	ScrollTrackStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#333333"))

	ShimmerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#444444"))

	ShimmerHighlightStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#666666"))
)

var CategoryTagStyles = map[int]lipgloss.Style{
	9:  lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#1565c0")).Padding(0, 1),
	10: lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#2e7d32")).Padding(0, 1),
	11: lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#6a1b9a")).Padding(0, 1),
	12: lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#ff8f00")).Padding(0, 1),
	13: lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#c62828")).Padding(0, 1),
	14: lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#ad1457")).Padding(0, 1),
	15: lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#00838f")).Padding(0, 1),
	16: lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#4e342e")).Padding(0, 1),
	17: lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(lipgloss.Color("#283593")).Padding(0, 1),
}

func ScoreStyle(score int) lipgloss.Style {
	if score > 5 {
		return ScorePositive
	}
	if score > 0 {
		return ScoreNeutral
	}
	return ScoreNegative
}

func ScoreBadge(score int) lipgloss.Style {
	if score > 5 {
		return ScoreBadgePositive
	}
	if score > 0 {
		return ScoreBadgeNeutral
	}
	return ScoreBadgeNegative
}
