package ui

import (
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeisaraja/youmui/api"
	"github.com/jeisaraja/youmui/ui/components"
)

type SearchResultMsg struct {
	Songs []api.Song
}

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
	activeTab Tab
	tabs      []Tab
	state     ModelState
	client    *http.Client
	searchbar components.TextInput
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
		searchbar: *components.NewTextInputView(20, 20, nil),
	}
}

func (m *model) Init() tea.Cmd {
	batchCmds := []tea.Cmd{
		tea.EnterAltScreen,
	}
	for _, t := range m.tabs {
		batchCmds = append(batchCmds, t.item.Init())
	}
	batchCmds = append(batchCmds, m.searchbar.Init())
	return tea.Batch(
		batchCmds...,
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case SearchResultMsg:
		if songList, ok := m.activeTab.item.(*components.SongList); ok {
			songList.UpdateSongs(msg.Songs)
		}
		m.state = IdleMode
		return m, nil
	case tea.KeyMsg:
		if m.state == InputMode {
			if msg.String() == "esc" {
				m.state = IdleMode
			} else if msg.String() == "enter" {
				if m.ActiveTab().name == "Song" {
					searchQuery := m.searchbar.Input.Value()
					cmd = SearchSongCallback(searchQuery)
					cmds = append(cmds, cmd)
				}
				m.state = IdleMode
			} else if msg.String() == "ctrl+c" {
				return m, tea.Quit
			} else {
				_, cmd := m.searchbar.Update(msg)
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

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
			}
		case "esc":
			if m.state == InputMode {
				m.state = IdleMode
			}
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
		default:
			if m.state == InputMode {
				_, cmd := m.searchbar.Update(msg)
				cmds = append(cmds, cmd)
			} else {
				_, cmd := m.ActiveTab().item.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	tabs := tabGroupStyle.Render(lipgloss.JoinVertical(lipgloss.Left, m.activeTab.name))
	contentPaddingStyle := lipgloss.NewStyle().
		PaddingLeft(4)

	tabContent := contentPaddingStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.ActiveTab().item.View(),
		),
	)
	if m.state == InputMode {
		tabContent += m.searchbar.View()
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, tabs, tabContent)
}

func (m *model) ActiveTab() Tab {
	return m.activeTab
}

func SearchSongCallback(input string) tea.Cmd {
	return func() tea.Msg {
		file, err := tea.LogToFile("debug.log", "log:\n")
		defer file.Close()
		if err != nil {
			panic("err while opening debug.log")
		}
		file.WriteString("\nTHIS IS FROM INSIDE CALLBACK\n")
		res, err := api.SearchWithKeyword(api.NewClient(), input, 5)
		return SearchResultMsg{Songs: res}
	}
}
