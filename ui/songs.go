package ui

import (
	"database/sql"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeisaraja/youmui/api"
	"github.com/jeisaraja/youmui/storage"
)

type SongList struct {
	Songs      []api.Song
	hoverIndex int
	db         *sql.DB
}

func NewSongList(db *sql.DB) *SongList {
	var songComponents []*SongComponent
	songs, err := storage.GetSongs(db)
	if err != nil {
		panic(fmt.Errorf("[ERROR] while calling storage.GetSongs from NewSongList: %w", err))
	}
	for _, song := range songs {
		songComponents = append(songComponents, NewSong(song))
	}

	return &SongList{
		Songs: songs,
	}
}

func (sl *SongList) Init() tea.Cmd {
	return nil
}

func (sl *SongList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch keyMsg := msg.(type) {
	case tea.KeyMsg:
		key := keyMsg.String()

		switch key {
		case "down", "j":
			if sl.hoverIndex < len(sl.Songs)-1 {
				sl.hoverIndex++
			}
		case "up", "k":
			if sl.hoverIndex > 0 {
				sl.hoverIndex--
			}
		case "enter":
			return sl, sl.PlaySelectedSong()
		case "a":
			return sl, sl.AddToQueue()
		case "i":
			song := sl.GetSelectedSong()
			if song != nil {
				return sl, sl.ShowPlaylistOptions()
			}
		}
	}

	return sl, nil
}

func (sl *SongList) View() string {
	normalStyle := lipgloss.NewStyle().Width(70)
	hoverStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("75")).Width(70)
	view := ""
	for i, song := range sl.Songs {
		if i == sl.hoverIndex {
			view += hoverStyle.Render(song.Title) + "\n"
		} else {
			view += normalStyle.Render(song.Title) + "\n"
		}
	}
	return lipgloss.NewStyle().PaddingTop(2).Render(view)
}

func (sl *SongList) UpdateSongs(songs []api.Song) {
	var songComponents []*SongComponent
	for _, song := range songs {
		songComponents = append(songComponents, NewSong(song))
	}
	sl.Songs = songs
}

type SongEnqueuedMsg struct {
	Song api.Song
}

type PlaySongMsg struct {
	Song api.Song
}

func (sl *SongList) PlaySelectedSong() tea.Cmd {
	if len(sl.Songs) == 0 {
		return nil
	}
	if sl.hoverIndex > len(sl.Songs) {
		return nil
	}
	song := sl.Songs[sl.hoverIndex]

	return func() tea.Msg {
		return PlaySongMsg{Song: song}
	}
}

func (sl *SongList) AddToQueue() tea.Cmd {
	song := sl.Songs[sl.hoverIndex]

	return func() tea.Msg {
		return SongEnqueuedMsg{Song: song}
	}
}

func (s *SongList) GetSelectedSong() *api.Song {
	if len(s.Songs) == 0 {
		return nil
	}
	return &s.Songs[s.hoverIndex]
}

type ShowPlaylistOptions struct{}

func (sl *SongList) ShowPlaylistOptions() tea.Cmd {
	return func() tea.Msg {
		return ShowPlaylistOptions{}
	}
}
