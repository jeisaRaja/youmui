package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeisaraja/youmui/api"
	"github.com/jeisaraja/youmui/ui/components"
	"github.com/jeisaraja/youmui/ui/types"
)

type Tab struct {
	name string
	item components.ContentModel
}

var (
	SongTab     = Tab{name: "Song", item: components.NewTrendingContent()}
	PlaylistTab = Tab{name: "Playlist", item: components.NewBaseView("", "playlist", "")}
	QueueTab    = Tab{name: "Queue", item: components.NewBaseView("", "queue", "")}
)

type model struct {
	activeTab    Tab
	level        types.FocusLevel
	focusOnQueue bool
	tabs         []Tab
}

func NewModel() *model {
	tabs := []Tab{SongTab, PlaylistTab, QueueTab}
	return &model{
		activeTab: SongTab,
		level:     types.TabsFocus,
		tabs:      tabs,
	}
}

func (m *model) Init() tea.Cmd {
	batchCmds := []tea.Cmd{
		tea.EnterAltScreen,
	}
	for _, t := range m.tabs {
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
	case *api.SearchResult:
		_, cmd = m.ActiveTab().Update(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			if m.level == types.TabsFocus {
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


func (m *model) IncrementFocus() types.FocusLevel {
	if m.level == types.TabsFocus {
		return types.ContentFocus
	} else {
		return types.ContentFocus
	}
}

func (m *model) DecrementFocus() types.FocusLevel {
	if m.level == types.SubContentFocus {
		return types.ContentFocus
	} else {
		return types.TabsFocus
	}
}
