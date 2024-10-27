package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type SelectPlaylist struct {
	cur       int
	playlists []Playlist
}

func NewSelectPlaylist(ps []struct {
	Title string
	ID    int64
	Count int64
}) *SelectPlaylist {
	var playlists []Playlist
	for _, p := range ps {
		pl := Playlist{
			Title: p.Title,
			ID:    p.ID,
		}
		playlists = append(playlists, pl)
	}

	return &SelectPlaylist{
		cur:       0,
		playlists: playlists,
	}
}

func (sp *SelectPlaylist) Init() tea.Cmd {
	return nil
}

func (sp *SelectPlaylist) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if sp.cur < len(sp.playlists)-1 {
				sp.cur++
			}
		case "k", "up":
			if sp.cur > 0 {
				sp.cur--
			}
		}
	}
	return sp, nil
}

func (sp *SelectPlaylist) View() string {
	var output strings.Builder

	output.WriteString("Select playlist:\n")

	for i, playlist := range sp.playlists {
		if i == sp.cur {
			output.WriteString(fmt.Sprintf("> [%s]\n", playlist.Title))
		} else {
			output.WriteString(fmt.Sprintf("  %s\n", playlist.Title))
		}
	}

	return output.String()
}

func (sp *SelectPlaylist) GetSelectedPlaylist() int64 {
	return sp.playlists[sp.cur].ID
}

func (sp *SelectPlaylist) AppendPlaylist(title string, id int64, count int64) {
	sp.playlists = append(sp.playlists, Playlist{Title: title, ID: id, Count: count})
}
