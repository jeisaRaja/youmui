package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeisaraja/youmui/ui/components"
	"github.com/jeisaraja/youmui/ui/types"
)

type model struct {
	tabs   *Tabs
	search *components.SearchContent
	level  types.FocusLevel
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
	m.level = types.TabsFocus
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
	var focusMsg types.FocusMsg
	focusMsg.Level = m.level
	focusMsg.Msg = msg

	switch msg := msg.(type) {
	case types.Mockres:
		_, cmd = m.ActiveTab().Update(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			if m.tabs.content_focus == false {
				m.tabs.SetTab("Search")
				// m.tabs.content_focus = true
				m.level = types.ContentFocus
				return m, nil
			}
		}
		switch msg.Type {
		case tea.KeyEsc:
			m.level = m.DecrementFocus()
			// m.tabs.content_focus = false
			return m, nil
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			// m.tabs.content_focus = true
			_, cmd = m.ActiveTab().Update(focusMsg)
			cmds = append(cmds, cmd)
			_, cmd = m.tabs.Update(focusMsg)
			cmds = append(cmds, cmd)
			m.level = m.IncrementFocus()
		default:
			if m.level == types.TabsFocus {
				_, cmd = m.tabs.Update(msg)
			}
			if m.level == types.ContentFocus {
				_, cmd = m.ActiveTab().Update(msg)
			}
			// if m.tabs.content_focus {
			// 	_, cmd = m.ActiveTab().Update(msg)
			// } else {
			// 	_, cmd = m.tabs.Update(msg)
			// }
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

func (m *model) IncrementFocus() types.FocusLevel {
	if m.level == types.TabsFocus {
		return types.ContentFocus
	} else {
		return types.SubContentFocus
	}
}

func (m *model) DecrementFocus() types.FocusLevel {
	if m.level == types.SubContentFocus {
		return types.ContentFocus
	} else {
		return types.TabsFocus
	}
}
