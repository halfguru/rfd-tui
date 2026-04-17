package main

import (
	"github.com/simon/rfd/internal/client"
	"github.com/simon/rfd/internal/views"
	tea "charm.land/bubbletea/v2"
)

type activeView int

const (
	viewDealList activeView = iota
)

type Model struct {
	activeView activeView
	dealList   views.DealListModel
	client     *client.Client
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
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.activeView == viewDealList {
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.dealList, cmd = m.dealList.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	switch m.activeView {
	case viewDealList:
		return m.dealList.View()
	default:
		return tea.NewView("")
	}
}
