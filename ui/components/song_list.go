package components

import (
	tea "github.com/charmbracelet/bubbletea"
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
			return sl, sl.PlaySong()
		}
	}

	return sl, nil
}

func (sl *SongList) View() string {
	view := ""
	for i, song := range sl.Songs {
		if i == sl.hoverIndex {
			view += ">> " + song.Title + "\n"
		} else {
			view += song.Title + "\n"
		}
	}
	return view
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
	song := sl.Songs[sl.hoverIndex]

	return func() tea.Msg {
		return PlaySongMsg{Song: song}
	}
}

func (sl *SongList) PlaySong() tea.Cmd {
	song := sl.Songs[sl.hoverIndex]

	return func() tea.Msg {
		return SongEnqueuedMsg{Song: song}
	}
}

func (s *SongList) GetSelectedSong() api.Song {
	return s.Songs[s.hoverIndex]
}
