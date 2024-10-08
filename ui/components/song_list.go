package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaraja/youmui/api"
	"github.com/jeisaraja/youmui/stream"
)

type SongList struct {
	Songs      []*SongComponent
	hoverIndex int
}

func NewSongList(songs []api.Song) *SongList {
	var songComponents []*SongComponent
	for _, song := range songs {
		songComponents = append(songComponents, NewSong(song))
	}

	return &SongList{
		Songs: songComponents,
	}
}

func (sl *SongList) Init() tea.Cmd {
	return nil
}

func (sl *SongList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if (msg.String() == "down" || msg.String() == "j") && sl.hoverIndex < len(sl.Songs)-1 {
			sl.Songs[sl.hoverIndex].SetHovered(false)
			sl.hoverIndex++
			sl.Songs[sl.hoverIndex].SetHovered(true)
		} else if (msg.String() == "up" || msg.String() == "k") && sl.hoverIndex > 0 {
			sl.Songs[sl.hoverIndex].SetHovered(false)
			sl.hoverIndex--
			sl.Songs[sl.hoverIndex].SetHovered(true)
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
			view += ">> " + song.View() // Highlight hovered song
		} else {
			view += song.View()
		}
	}
	return view
}

func (sl *SongList) UpdateSongs(songs []api.Song) {
	var songComponents []*SongComponent
	for _, song := range songs {
		songComponents = append(songComponents, NewSong(song))
	}
	sl.Songs = songComponents
}

func (sl *SongList) PlaySong() tea.Cmd {
	songURL := sl.Songs[sl.hoverIndex].song.URL

	return func() tea.Msg {
		err := stream.FetchAndPlayAudio(songURL)
		if err != nil {
			return tea.Msg("Error playing audio")
		}
		return tea.Msg("Playback finished")
	}
}
