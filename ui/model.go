package ui

import (
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeisaraja/youmui/api"
	"github.com/jeisaraja/youmui/ui/components"
)

type Tab struct {
	name string
	item components.ContentModel
}

type ModelState int

const (
	IdleMode ModelState = iota
	InputMode
)

var (
	SongTab     Tab
	PlaylistTab = Tab{name: "Playlist", item: components.NewBaseView("", "playlist", "")}
	QueueTab    = Tab{name: "Queue", item: components.NewBaseView("", "queue", "")}
)

type model struct {
	activeTab   Tab
	tabs        []Tab
	state       ModelState
	searchInput string
	client      *http.Client
}

func NewModel(client *http.Client) *model {
	songs, err := api.GetTrendingMusic(client)
	if err != nil {
		panic(err)
	}
	SongTab = Tab{name: "Song", item: components.NewSongList(songs)}
	tabs := []Tab{SongTab, PlaylistTab, QueueTab}
	return &model{
		activeTab: SongTab,
		tabs:      tabs,
		state:     IdleMode,
		client:    client,
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

	switch msg := msg.(type) {
	case *api.SearchResult:
		_, cmd = m.ActiveTab().item.Update(msg)
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			if m.state == IdleMode {
				m.activeTab = SongTab
			}
		case "p":
			if m.state == IdleMode {
				m.activeTab = PlaylistTab
			}
		case "q":
			if m.state == IdleMode {
				m.activeTab = QueueTab
			}
		case "f":
			if m.state == IdleMode && (m.activeTab == SongTab || m.activeTab == PlaylistTab) {
				m.state = InputMode
			} else {
				m.state = IdleMode
			}
		case "esc":
			if m.state == InputMode {
				m.state = IdleMode
			}
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.state == InputMode {
				// Execute search or other functionality based on input
				// You can perform search logic here with m.searchInput
				m.state = IdleMode // Exit input mode after processing
			}
		default:
			if m.state == InputMode {
				// Register alphabet key inputs in Input Mode
				if len(msg.String()) == 1 && 'a' <= msg.String()[0] && msg.String()[0] <= 'z' {
					m.searchInput += msg.String()
				}
			}
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	tabs := lipgloss.JoinVertical(
		lipgloss.Left,
	)
	tabContent := lipgloss.JoinVertical(lipgloss.Left, m.ActiveTab().item.View())
	if m.state == InputMode {
		tabContent += "\n" + components.NewTextInputView(m.activeTab.name, 20, 30, nil).View()
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, tabs, tabContent)
}

func (m *model) ActiveTab() Tab {
	return m.activeTab
}
