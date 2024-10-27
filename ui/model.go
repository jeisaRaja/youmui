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
)

type ContentModel interface {
	tea.Model
}

type Tab struct {
	name string
	item ContentModel
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
	input            TextInput
	inputHeader      string
	loading          bool
	spinner          spinner.Model
	queue            *Queue
	hasActiveSong    bool
	height           int
	width            int
	player           *stream.Player
	store            *sql.DB
	isSelectPlaylist bool
	selectPlaylist   *SelectPlaylist
	playlistTab      *PlaylistComponent
}

func NewModel(client *http.Client, db *sql.DB) *model {

	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = spinnerStyle
	SongTab = Tab{name: "Song", item: NewSongList(db)}
	queueItem := NewQueue()
	QueueTab = Tab{name: "Queue", item: queueItem}

	ps, _ := storage.GetPlaylists(db)
	playlistComp := NewPlaylistComponent(db)
	playlistComp.SetPlaylists(ps)
	PlaylistTab = Tab{name: "Playlist", item: playlistComp}
	tabs := []Tab{SongTab, PlaylistTab}
	sp := NewSelectPlaylist(ps)

	return &model{
		activeTab:        SongTab,
		tabs:             tabs,
		state:            IdleMode,
		client:           client,
		input:            *NewTextInputView(50, 50),
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
	case ShowPlaylistOptions:
		m.handleAddToPlaylist()
	case SongEnqueuedMsg:
		m.queue.AddToQueue(msg.Song)
		if len(m.queue.Songs) <= 1 && !m.hasActiveSong {
			return m, m.PlayNextSong()
		}
	case SongsFromPlaylist:
		_, cmd := m.playlistTab.Update(msg)
		return m, cmd
	case PlaybackError:
		panic(msg)
	case PlayPlaylistMsg:
		m.queue.AddManyToQueue(msg.Songs)
		return m, m.PlayNextSong()
	case SongAddedToPlaylistMsg:
		m.playlistTab.IncrementCount(msg.PlaylistID)
	case CreatePlaylistMsg:
		m.selectPlaylist.AppendPlaylist(msg.Title, msg.ID, msg.Count)
		m.playlistTab.AppendPlaylist(msg.Title, msg.ID, msg.Count)
		m.loading = false
	case InputStateMsg:
		m.state = InputMode
	case tea.WindowSizeMsg:
		m.width = msg.Height
		m.height = msg.Height
	case PlaySongMsg:
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
		case "esc":
			if m.isInputMode() {
				m.state = IdleMode
				return m, nil
			}
			if m.isSelectPlaylist {
				m.isSelectPlaylist = false
			}
			m.ActiveTab().item.Update(msg)
		case "tab":
			m.toggleTab()
			return m, nil
		case "enter":
			if m.state == InputMode {
				m.state = IdleMode
			}
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
		case "f":
			m.isSelectPlaylist = false
			m.state = InputMode
			if m.activeTab == PlaylistTab {
				m.input.SetUsage(SearchPlaylist)
				m.inputHeader = "Search for playlist"
			} else {
				m.inputHeader = "Search for song"
				m.input.SetUsage(SearchSong)
			}
			m.input.SetPlaceholder(fmt.Sprintf("%s title", m.activeTab.name))
		case "c":
			m.state = InputMode
			m.inputHeader = "Create a new playlist"
			m.input.SetUsage(CreatePlaylist)
			m.input.SetPlaceholder("Playlist title")
			m.isSelectPlaylist = false
		case "ctrl+c":
			return m, tea.Quit
		default:
			if m.isSelectPlaylist {
				_, cmd := m.selectPlaylist.Update(msg)
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

	if m.isInputMode() {
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

func (m *model) PlaySelectedSong(selectedSong api.Song) tea.Cmd {
	m.queue.Clear()
	m.queue.AddToQueue(selectedSong)

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
	if songList, ok := m.activeTab.item.(*SongList); ok {
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
		case SearchSong:
			return SearchSongCallback(inputValue)
		case CreatePlaylist:
			return CreatePlaylistCallback(m.store, inputValue)
		}
	}

	if songList, ok := m.ActiveTab().item.(*SongList); ok {
		if m.isSelectPlaylist {
			m.isSelectPlaylist = false
			cmd := AddSongToPlaylistCallback(m.store, m.selectPlaylist.GetSelectedPlaylist(), *songList.GetSelectedSong())
			return cmd
		}
		_, cmd := songList.Update(msg)
		return cmd
	}

	if playlist, ok := m.ActiveTab().item.(*PlaylistComponent); ok {
		_, cmd := playlist.Update(msg)
		return cmd
	}

	return nil
}

func (m *model) handleAddToPlaylist() {
	m.isSelectPlaylist = !m.isSelectPlaylist
	if !m.isSelectPlaylist {
		return
	}
	m.isSelectPlaylist = true
	return
}

func (m *model) toggleTab() {
	currentIndex := -1
	for i, tab := range m.tabs {
		if tab.name == m.activeTab.name {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return
	}
	nextIndex := (currentIndex + 1) % len(m.tabs)
	m.isSelectPlaylist = false
	m.activeTab = m.tabs[nextIndex]
}

func (m *model) isInputMode() bool {
	return m.state == InputMode
}
