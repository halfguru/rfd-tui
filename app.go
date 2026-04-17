package main

import (
	"fmt"
	"os/exec"
	"runtime"

	tea "charm.land/bubbletea/v2"
	"github.com/simon/rfd/internal/client"
	"github.com/simon/rfd/internal/views"
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
	client     *client.Client
	width      int
	height     int
}

func NewModel() Model {
	c := client.New()
	return Model{
		activeView: viewDealList,
		dealList:   views.NewDealList(c),
		client:     c,
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

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.activeView == viewDealList {
				return m, tea.Quit
			}
			m.activeView = viewDealList
			m.thread = views.ThreadModel{}
			return m, nil
		case "escape":
			if m.activeView == viewThread {
				m.activeView = viewDealList
				m.thread = views.ThreadModel{}
				return m, nil
			}
			return m, tea.Quit
		case "enter":
			if m.activeView == viewDealList {
				if t := m.dealList.SelectedTopic(); t != nil {
					m.activeView = viewThread
					m.thread = views.NewThread(t, m.client)
					return m, m.thread.Init()
				}
				return m, nil
			}
		case "o":
			if m.activeView == viewDealList {
				if t := m.dealList.SelectedTopic(); t != nil {
					return m, openBrowser(t.DealURL())
				}
			}
		}
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

func (m Model) View() tea.View {
	switch m.activeView {
	case viewDealList:
		return m.dealList.View()
	case viewThread:
		return m.thread.View()
	default:
		return tea.NewView("")
	}
}

func openBrowser(url string) tea.Cmd {
	return func() tea.Msg {
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
			return client.ErrMsg{Err: fmt.Errorf("failed to open browser: %w", err)}
		}
		return nil
	}
}
