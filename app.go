package main

import (
	"fmt"
	"log/slog"
	"os/exec"
	"runtime"

	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"
	"github.com/halfguru/rfd-tui/internal/client"
	"github.com/halfguru/rfd-tui/internal/config"
	"github.com/halfguru/rfd-tui/internal/views"
)

type activeView int

const (
	viewDealList activeView = iota
	viewThread
)

type Model struct {
	activeView activeView
	dealList   views.DealListModel
	thread     views.ThreadModel
	help       views.HelpModel
	showHelp   bool
	client     *client.Client
	config     config.Config
	width      int
	height     int
	altScreen  bool
	statusMsg  string
}

func NewModel(cfg config.Config) Model {
	c := client.New()
	return Model{
		activeView: viewDealList,
		dealList:   views.NewDealList(c),
		client:     c,
		config:     cfg,
		altScreen:  cfg.AltScreen,
	}
}

func (m Model) Init() tea.Cmd {
	return m.dealList.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	if m.showHelp {
		return m.updateHelp(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if m.activeView == viewDealList && m.dealList.SearchActive() {
			var cmd tea.Cmd
			m.dealList, cmd = m.dealList.Update(msg)
			return m, cmd
		}
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "?":
			m.showHelp = true
			viewName := "Deal List"
			if m.activeView == viewThread {
				viewName = "Thread"
			}
			m.help = views.NewHelp(viewName)
			return m, nil
		case "q":
			if m.activeView == viewDealList {
				return m, tea.Quit
			}
			m.activeView = viewDealList
			m.thread = views.ThreadModel{}
			m.statusMsg = ""
			return m, nil
		case "esc":
			if m.activeView == viewThread {
				m.activeView = viewDealList
				m.thread = views.ThreadModel{}
				m.statusMsg = ""
				return m, nil
			}
			return m, tea.Quit
		case "enter":
			if m.activeView == viewDealList {
				if t := m.dealList.SelectedTopic(); t != nil {
					m.activeView = viewThread
					m.thread = views.NewThread(t, m.client, m.width, m.height)
					m.statusMsg = ""
					return m, m.thread.Init()
				}
				return m, nil
			}
		case "o":
			url := m.currentDealURL()
			if url != "" {
				return m, openBrowser(url)
			}
		case "c":
			url := m.currentDealURL()
			if url != "" {
				return m, copyToClipboard(url)
			}
		case "ctrl+a":
			m.altScreen = !m.altScreen
			return m, nil
		}
	case clipboardMsg:
		m.statusMsg = "Copied to clipboard"
		return m, nil
	case client.ErrMsg:
		m.statusMsg = fmt.Sprintf("Error: %v", msg.Err)
	}

	var cmd tea.Cmd
	switch m.activeView {
	case viewDealList:
		m.dealList, cmd = m.dealList.Update(msg)
	case viewThread:
		m.thread, cmd = m.thread.Update(msg)
	}
	return m, cmd
}

func (m Model) currentDealURL() string {
	switch m.activeView {
	case viewDealList:
		if t := m.dealList.SelectedTopic(); t != nil {
			return t.DealURL()
		}
	case viewThread:
		if m.thread.Topic() != nil {
			return m.thread.Topic().DealURL()
		}
	}
	return ""
}

func (m Model) updateHelp(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyPressMsg:
		m.showHelp = false
		return m, nil
	}
	return m, nil
}

func (m Model) View() tea.View {
	var v tea.View
	switch m.activeView {
	case viewDealList:
		v = m.dealList.View()
	case viewThread:
		v = m.thread.View()
	default:
		v = tea.NewView("")
	}

	if m.showHelp {
		return m.help.View()
	}

	v.AltScreen = m.altScreen
	if m.config.Mouse {
		v.MouseMode = tea.MouseModeCellMotion
	}
	return v
}

type clipboardMsg struct{}

func copyToClipboard(text string) tea.Cmd {
	return func() tea.Msg {
		_ = clipboard.WriteAll(text)
		return clipboardMsg{}
	}
}

func openBrowser(url string) tea.Cmd {
	return func() tea.Msg {
		slog.Info("opening browser", "url", url)
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", url)
		case "linux":
			cmd = exec.Command("xdg-open", url)
		default:
			return client.ErrMsg{Err: fmt.Errorf("unsupported platform")}
		}
		if err := cmd.Start(); err != nil {
			slog.Error("failed to open browser", "url", url, "error", err)
			return client.ErrMsg{Err: fmt.Errorf("failed to open browser: %w", err)}
		}
		return nil
	}
}
