package ui

import (
	"net/http"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeisaraja/youmui/api"
	"github.com/jeisaraja/youmui/stream"
	"github.com/jeisaraja/youmui/ui/components"
)

type SongListMsg struct {
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
	QueueTab    Tab

	textStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
)

type model struct {
	activeTab Tab
	tabs      []Tab
	state     ModelState
	client    *http.Client
	searchbar components.TextInput
	loading   bool
	spinner   spinner.Model
	queue     *components.Queue
	isPlaying bool
	height    int
	width     int
	player    *stream.Player
}

func NewModel(client *http.Client) *model {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = spinnerStyle
	SongTab = Tab{name: "Song", item: components.NewSongList([]api.Song{})}
	queueItem := components.NewQueue()
	QueueTab = Tab{name: "Queue", item: queueItem}
	tabs := []Tab{SongTab, PlaylistTab, QueueTab}
	return &model{
		activeTab: SongTab,
		tabs:      tabs,
		state:     IdleMode,
		client:    client,
		searchbar: *components.NewTextInputView(50, 50, nil),
		loading:   false,
		spinner:   s,
		queue:     queueItem,
		player:    stream.NewPlayer(),
	}
}

func (m *model) Init() tea.Cmd {
	batchCmds := []tea.Cmd{
		tea.EnterAltScreen,
	}
	for _, t := range m.tabs {
		batchCmds = append(batchCmds, t.item.Init())
	}
	batchCmds = append(batchCmds, m.loadSongsAsync())
	batchCmds = append(batchCmds, m.searchbar.Init())
	batchCmds = append(batchCmds, m.spinner.Tick)
	return tea.Batch(
		batchCmds...,
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.loading {
		_, cmd := m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Height
		m.height = msg.Height
	case components.PlaySongMsg:
		if songList, ok := m.activeTab.item.(*components.SongList); ok {
			selectedSong := songList.GetSelectedSong()
			m.PlaySelectedSong(selectedSong)
		}
	case PlaybackFinished:
		if len(m.queue.Songs) > 0 {
			cmds = append(cmds, m.PlayNextSong())
		} else {
			m.isPlaying = false
			m.queue.PlayingSong = nil
		}
	case components.SongEnqueuedMsg:
		m.AddToQueue(msg.Song)
		if len(m.queue.Songs) == 1 && !m.isPlaying {
			cmds = append(cmds, m.PlayNextSong())
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case SongListMsg:
		m.loading = false
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
					m.searchbar.Input.SetValue("")
					cmd = SearchSongCallback(searchQuery)
					m.loading = true
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
		case "enter":
			if songList, ok := m.ActiveTab().item.(*components.SongList); ok {
				selectedSong := songList.GetSelectedSong()
				cmds = append(cmds, m.PlaySelectedSong(selectedSong))
				m.queue.View()
			}
		case "-":
			m.player.VolumeDown()
			return m, nil
		case "=":
			m.player.VolumeUp()
			return m, nil
		case "n":
			m.StopPlayback()
			cmds = append(cmds, m.PlayNextSong())
		case "x":
			m.player.PlayPause()
      m.queue.PlayPause()
			return m, nil
		case "a":
			if songList, ok := m.ActiveTab().item.(*components.SongList); ok {
				selectedSong := songList.GetSelectedSong()
				m.queue.AddToQueue(selectedSong)

				if len(m.queue.Songs) == 1 && !m.isPlaying {
					cmds = append(cmds, m.PlayNextSong())
				}
				if m.ActiveTab().name == "Queue" {
					m.ActiveTab().item.View()
				}
			}
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
	tabs := activeTabStyle.Render(lipgloss.JoinVertical(lipgloss.Left, m.activeTab.name))
	contentPaddingStyle := lipgloss.NewStyle().
		PaddingLeft(4)

	tabContent := contentPaddingStyle.Width(80).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.ActiveTab().item.View(),
		),
	)

	if m.state == InputMode {
		searchBar := lipgloss.NewStyle().
			Width(80).
			PaddingLeft(4).
			Render(m.searchbar.View())
		tabContent += "\n\n    Search For " + m.ActiveTab().name + "\n\n" + searchBar
	}

	if m.loading {
		tabContent += "\n\n" + m.spinner.View()
	}

	queueStyled := queueTabStyle.Height(m.height).Render(m.queue.View())
	return lipgloss.JoinHorizontal(lipgloss.Left, tabs, tabContent, queueStyled)
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
		file.WriteString("Search API call with being called!\n")
		res, err := api.SearchWithKeyword(api.NewClient(), input, 5)
		if err != nil {
			file.WriteString("using yt-dlp to search")
			res, err = api.SearchWithKeywordWithoutApi(input)
		}
		return SongListMsg{Songs: res}
	}
}

func (m *model) loadSongsAsync() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		m.loading = true
		songs, err := api.GetTrendingMusic(m.client)
		if err != nil {
			m.loading = false
			return SongListMsg{Songs: nil}
		}

		if songlist, ok := m.ActiveTab().item.(*components.SongList); ok {
			songlist.UpdateSongs(songs)
		}

		m.loading = false
		return SongListMsg{Songs: songs}
	})
}

type PlaybackFinished struct{}

func (m *model) PlayNextSong() tea.Cmd {
	if len(m.queue.Songs) == 0 {
		return nil
	}

	nextSong := m.queue.RemoveFromQueue()
	m.queue.SetPlayingSong(nextSong)

	return func() tea.Msg {
		m.isPlaying = true
		err := m.player.FetchAndPlayAudio(nextSong.URL)
		if err != nil {
			m.isPlaying = false
			return tea.Msg("Error playing audio")
		}

		return PlaybackFinished{}
	}
}

func (m *model) AddToQueue(song api.Song) tea.Cmd {
	m.queue.AddToQueue(song)
	return nil
}

func (m *model) StopPlayback() tea.Cmd {
	m.isPlaying = false
	return nil
}

func (m *model) PlaySelectedSong(selectedSong api.Song) tea.Cmd {
	m.StopPlayback()

	m.queue.Clear()
	m.AddToQueue(selectedSong)

	return m.PlayNextSong()
}
