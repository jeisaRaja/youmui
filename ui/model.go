package ui

import (

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	tabs *Tabs
}

var (
	tabItems = map[string]string{
		"Home":     "Main Area",
		"Browse":   "Browse Musics",
		"Playlist": "My Playlists",
		"Search":   "Search Results",
		"Library":  "My Library",
	}
)

func NewModel() *model {
	var m model
	var tabList []TabItem
	for name, item := range tabItems {
		tabList = append(tabList, TabItem{name: name, item: NewBaseView("header", item, "footer")})
	}
	m.tabs = NewTabs(tabList)
	return &m
}

func (m *model) Init() tea.Cmd {
	batchCmds := []tea.Cmd{
		tea.EnterAltScreen,
	}
	for _, t := range m.tabs.tabList {
		batchCmds = append(batchCmds, t.item.Init())
	}
	return tea.Batch(
		batchCmds...,
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.tabs.Update(msg)
  switch msg.String(){
    case "ctrl+c":
    return m, tea.Quit
  }
	}
	return m, nil
}

func (m *model) View() string {
	tabs := lipgloss.JoinVertical(
		lipgloss.Left,
		m.tabs.View(),
	)
	tabContent := lipgloss.JoinVertical(lipgloss.Left, m.ActiveTab().View())
	return lipgloss.JoinHorizontal(lipgloss.Left, tabs, tabContent)
}

func (m *model) ActiveTab() ContentModel {
	item := m.tabs.CurrentTab().item
	return item
}