package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeisaraja/youmui/ui/components"
)

type model struct {
	tabs   *Tabs
	search *components.SearchContent
}

var (
	tabItems = map[string]string{
		"Home":     "Main Area",
		"Trending": "Browse Trending Musics",
		"Playlist": "My Playlists",
		"Search":   "Search Results",
		"Library":  "My Library",
	}
)

func NewModel() *model {
	var m model
	var tabList []TabItem
	for name, item := range tabItems {
		switch name {
		case "Search":
			searchComp := components.NewSearchContent("placeholder", 120, 50)
			tabList = append(tabList, TabItem{name: name, item: searchComp})
			m.search = searchComp
		default:
			tabList = append(tabList, TabItem{name: name, item: components.NewBaseView("header", item, "footer")})
		}
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
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			if m.tabs.content_focus == false {
				m.tabs.SetTab("Search")
				m.tabs.content_focus = true
				return m, nil
			}
		}
		switch msg.Type {
		case tea.KeyEsc:
			m.tabs.content_focus = false
			return m, nil
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			m.tabs.content_focus = true
			_, cmd = m.ActiveTab().Update(msg)
			cmds = append(cmds, cmd)
		default:
			if m.tabs.content_focus {
				_, cmd = m.ActiveTab().Update(msg)
			} else {
				_, cmd = m.tabs.Update(msg)
			}
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	tabs := lipgloss.JoinVertical(
		lipgloss.Left,
		m.tabs.View(),
	)
	tabContent := lipgloss.JoinVertical(lipgloss.Left, m.ActiveTab().View())
	return lipgloss.JoinHorizontal(lipgloss.Left, tabs, tabContent)
}

func (m *model) ActiveTab() components.ContentModel {
	item := m.tabs.CurrentTab().item
	return item
}
