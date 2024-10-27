package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeisaraja/youmui/api"
)

type SongList struct {
	Songs      []api.Song
	hoverIndex int
}

func NewSongList(songs []api.Song) *SongList {
	var songComponents []*SongComponent
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if (msg.String() == "down" || msg.String() == "j") && sl.hoverIndex < len(sl.Songs)-1 {
			sl.hoverIndex++
		} else if (msg.String() == "up" || msg.String() == "k") && sl.hoverIndex > 0 {
			sl.hoverIndex--
		}
		switch msg.String() {
		case "enter":
			return sl, sl.ForcePlaySong()
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

func (sl *SongList) ForcePlaySong() tea.Cmd {
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
