package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Playlist struct {
	ID    int64
	Title string
	Count int64
}

type PlaylistComponent struct {
	playlists []Playlist
}

type InputStateMsg struct{}

func NewPlaylistComponent() *PlaylistComponent {
	ti := textinput.New()
	ti.Placeholder = "Playlist title"
	ti.Focus()
	return &PlaylistComponent{}
}

func (p *PlaylistComponent) Init() tea.Cmd {
	return textinput.Blink
}

func (p *PlaylistComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		}
	}
	return p, nil
}

func (p *PlaylistComponent) View() string {
	var builder strings.Builder

	builder.WriteString("\n\n")

	builder.WriteString("Playlists:\n")
	for i, playlist := range p.playlists {
		builder.WriteString(fmt.Sprintf("%d. %s (%d songs)\n", i+1, playlist.Title, playlist.Count))
	}

	return builder.String()
}

func (p *PlaylistComponent) AppendPlaylist(title string, id int64, count int64) {
	p.playlists = append(p.playlists, Playlist{Title: title, ID: id, Count: count})
}

func (p *PlaylistComponent) SetPlaylists(data []struct {
	Title string
	ID    int64
	Count int64
}) {
	for _, ps := range data {
		p.AppendPlaylist(ps.Title, ps.ID, ps.Count)
	}
}

func (p *PlaylistComponent) IncrementCount(pid int64) {
	for i := range p.playlists {
		if p.playlists[i].ID == pid {
			p.playlists[i].Count++
		}
	}
}
