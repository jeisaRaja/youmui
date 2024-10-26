package ui

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeisaraja/youmui/api"
	"github.com/jeisaraja/youmui/storage"
	"github.com/jeisaraja/youmui/stream"
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
	PlaylistTab Tab
	QueueTab    Tab

	textStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
)

type model struct {
	activeTab        Tab
	tabs             []Tab
	state            ModelState
	client           *http.Client
	input            components.TextInput
	inputHeader      string
	loading          bool
	spinner          spinner.Model
	queue            *components.Queue
	hasActiveSong    bool
	height           int
	width            int
	player           *stream.Player
	store            *sql.DB
	isSelectPlaylist bool
	selectPlaylist   *components.SelectPlaylist
	playlistTab      *components.PlaylistComponent
}

func NewModel(client *http.Client, db *sql.DB) *model {

	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = spinnerStyle
	SongTab = Tab{name: "Song", item: components.NewSongList([]api.Song{})}
	queueItem := components.NewQueue()
	QueueTab = Tab{name: "Queue", item: queueItem}

	ps, _ := storage.GetPlaylists(db)
	playlistComp := components.NewPlaylistComponent(db)
	playlistComp.SetPlaylists(ps)
	PlaylistTab = Tab{name: "Playlist", item: playlistComp}

	tabs := []Tab{SongTab, PlaylistTab, QueueTab}

	sp := components.NewSelectPlaylist(ps)
	return &model{
		activeTab:        SongTab,
		tabs:             tabs,
		state:            IdleMode,
		client:           client,
		input:            *components.NewTextInputView(50, 50),
		loading:          false,
		spinner:          s,
		queue:            queueItem,
		player:           stream.NewPlayer(),
		store:            db,
		selectPlaylist:   sp,
		isSelectPlaylist: false,
		playlistTab:      playlistComp,
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
	batchCmds = append(batchCmds, m.input.Init())
	batchCmds = append(batchCmds, m.spinner.Tick)
	return tea.Batch(
		batchCmds...,
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.loading {
		_, cmd := m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case components.SongsFetch:
		_, cmd := m.playlistTab.Update(msg)
		return m, cmd
	case PlaybackError:
		panic(msg)
	case PlayEntirePlaylistMsg:
		cmd := m.AddManyToQueue(msg.Songs)
		return m, cmd
	case SongAddedToPlaylistMsg:
		m.playlistTab.IncrementCount(msg.PlaylistID)
	case CreatePlaylistMsg:
		m.selectPlaylist.AppendPlaylist(msg.Title, msg.ID, msg.Count)
		m.playlistTab.AppendPlaylist(msg.Title, msg.ID, msg.Count)
		m.loading = false
	case components.InputStateMsg:
		m.state = InputMode
	case tea.WindowSizeMsg:
		m.width = msg.Height
		m.height = msg.Height
	case components.PlaySongMsg:
		return m, m.PlaySelectedSong(msg.Song)
	case PlaybackFinished:
		return m, m.handlePlaybackFinished()
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case SongListMsg:
		return m, m.handleSongListMsg(msg)
	case tea.KeyMsg:
		if m.state == InputMode {
			switch msg.String() {
			case "esc":
				m.state = IdleMode
				m.isSelectPlaylist = false
				return m, nil

			case "enter":
				cmd := m.handleEnterKey(msg)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
				m.state = IdleMode
			case "ctrl+c":
				return m, tea.Quit

			default:
				_, cmd := m.input.Update(msg)
				cmds = append(cmds, cmd)
			}

			return m, tea.Batch(cmds...)
		}

		switch msg.String() {
		case "P":
			pid := m.playlistTab.GetCurrent()
			return m, PlayEntirePlaylistCallback(m.store, pid)
		case "i":
			return m, m.handleAddToPlaylist()
		case "enter":
			return m, m.handleEnterKey(msg)
		case "-":
			m.player.VolumeDown()
			return m, nil
		case "+":
			m.player.VolumeUp()
			return m, nil
		case "n":
			cmds = append(cmds, m.PlayNextSong())
		case "x":
			m.player.PlayPause()
			m.queue.PlayPause()
			return m, nil
		case "a":
			m, cmd := m.handleAddToQueue()
			return m, cmd
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
			m.isSelectPlaylist = false
			if m.state == IdleMode {
				m.state = InputMode
			}
			if m.activeTab == PlaylistTab {
				m.input.SetUsage(components.SearchPlaylist)
				m.inputHeader = "Search for playlist"
			} else {
				m.inputHeader = "Search for song"
				m.input.SetUsage(components.SearchSong)
			}
			m.input.SetPlaceholder(fmt.Sprintf("%s title", m.activeTab.name))
		case "c":
			if m.state == IdleMode {
				m.state = InputMode
			}
			m.inputHeader = "Create a new playlist"
			m.input.SetUsage(components.CreatePlaylist)
			m.input.SetPlaceholder("Playlist title")
			m.isSelectPlaylist = false
		case "esc":
			m.state = IdleMode
			m.isSelectPlaylist = false
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		default:
			if m.isSelectPlaylist {
				_, cmd := m.selectPlaylist.Update(msg)
				cmds = append(cmds, cmd)
			} else if m.state == InputMode {
				_, cmd := m.input.Update(msg)
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
			Render(m.input.View())
		tabContent += "\n\n    " + m.inputHeader + "\n\n" + searchBar
	}

	if m.loading {
		tabContent += "\n\n" + m.spinner.View()
	}

	if m.isSelectPlaylist {
		tabContent += "\n\n" + m.selectPlaylist.View()
	}

	queueStyled := queueTabStyle.Height(m.height).Render(m.queue.View())
	return lipgloss.JoinHorizontal(lipgloss.Left, tabs, tabContent, queueStyled)
}

func (m *model) ActiveTab() Tab {
	return m.activeTab
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
type PlaybackError string

func (m *model) PlayNextSong() tea.Cmd {
	if len(m.queue.Songs) == 0 {
		m.hasActiveSong = false
		return nil
	}
	m.hasActiveSong = true

	nextSong := m.queue.RemoveFromQueue()
	m.queue.SetPlayingSong(nextSong)

	return func() tea.Msg {
		err := m.player.FetchAndPlayAudio(nextSong.URL)
		if err != nil {
			return PlaybackError(err.Error())
		}

		return PlaybackFinished{}
	}
}

func (m *model) AddToQueue(song api.Song) tea.Cmd {
	m.queue.AddToQueue(song)
	return nil
}

func (m *model) AddManyToQueue(songs []api.Song) tea.Cmd {
	m.queue.Clear()
	for _, song := range songs {
		m.AddToQueue(song)
	}
	return m.PlayNextSong()
}

func (m *model) PlaySelectedSong(selectedSong api.Song) tea.Cmd {
	m.queue.Clear()
	m.AddToQueue(selectedSong)

	return m.PlayNextSong()
}

func (m *model) handlePlaybackFinished() tea.Cmd {
	if len(m.queue.Songs) > 0 {
		return m.PlayNextSong()
	} else {
		m.hasActiveSong = false
		m.queue.PlayingSong = nil
	}
	return nil
}

func (m *model) handleSongListMsg(msg SongListMsg) tea.Cmd {
	m.loading = false
	if songList, ok := m.activeTab.item.(*components.SongList); ok {
		songList.UpdateSongs(msg.Songs)
	}
	m.state = IdleMode
	return nil
}

func (m *model) handleEnterKey(msg tea.KeyMsg) tea.Cmd {
	if m.state == InputMode {
		inputValue := m.input.Input.Value()
		m.input.Input.SetValue("")
		m.loading = true
		switch m.input.Usage {
		case components.SearchSong:
			return SearchSongCallback(inputValue)
		case components.CreatePlaylist:
			return CreatePlaylistCallback(m.store, inputValue)
		}
	}

	if songList, ok := m.ActiveTab().item.(*components.SongList); ok {
		if m.isSelectPlaylist {
			m.isSelectPlaylist = false
			cmd := AddSongToPlaylistCallback(m.store, m.selectPlaylist.GetSelectedPlaylist(), *songList.GetSelectedSong())
			return cmd
		}
		_, cmd := songList.Update(msg)
		return cmd
	}

	if playlist, ok := m.ActiveTab().item.(*components.PlaylistComponent); ok {
		_, cmd := playlist.Update(msg)
		return cmd
	}

	return nil
}

func (m *model) handleAddToQueue() (*model, tea.Cmd) {
	if songList, ok := m.ActiveTab().item.(*components.SongList); ok {
		selectedSong := songList.GetSelectedSong()
		if selectedSong == nil {
			return m, nil
		}

		m.queue.AddToQueue(*selectedSong)

		if len(m.queue.Songs) <= 1 && !m.hasActiveSong {
			return m, m.PlayNextSong()
		}

		return m, nil
	}

	return m, nil
}

func (m *model) handleAddToPlaylist() tea.Cmd {
	m.isSelectPlaylist = !m.isSelectPlaylist
	if !m.isSelectPlaylist {
		return nil
	}
	switch tab := m.activeTab.item.(type) {
	case *components.SongList:
		if tab.GetSelectedSong() == nil {
			m.isSelectPlaylist = false
			return nil
		}
		m.isSelectPlaylist = true
	}
	return nil
}
