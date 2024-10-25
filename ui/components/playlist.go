package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaraja/youmui/api"
)

type Playlist struct {
	ID    int64
	Title string
	Songs []api.Song
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
		builder.WriteString(fmt.Sprintf("%d. %s (%d songs)\n", i+1, playlist.Title, len(playlist.Songs)))
	}

	return builder.String()
}

func (p *PlaylistComponent) AppendPlaylist(title string, id int64) {
	p.playlists = append(p.playlists, Playlist{Title: title, ID: id})
}

func (p *PlaylistComponent) SetPlaylists(data []struct {
	Title string
	ID    int64
}) {
	for _, ps := range data {
		p.AppendPlaylist(ps.Title, ps.ID)
	}
}
