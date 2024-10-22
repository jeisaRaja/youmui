package components

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaraja/youmui/api"
)

type Playlist struct {
	Title string
	Songs []api.Song
}

type PlaylistComponent struct {
	TitleInput textinput.Model
	ShowInput  bool
	Storage    *sql.DB
}

type InputStateMsg struct{}

func NewPlaylistComponent(db *sql.DB) *PlaylistComponent {
	ti := textinput.New()
	ti.Placeholder = "Playlist title"
	ti.Focus()
	return &PlaylistComponent{
		TitleInput: ti,
		Storage:    db,
	}
}

func (p PlaylistComponent) Init() tea.Cmd {
	return textinput.Blink
}

func (p PlaylistComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "c":
			// return p, func() tea.Msg {
			// 	return InputStateMsg{}
			// }
		}
	}
	if p.ShowInput {
		var cmd tea.Cmd
		p.TitleInput, cmd = p.TitleInput.Update(msg)
		return p, cmd
	}
	return p, nil
}

func (p PlaylistComponent) View() string {
	if p.ShowInput {
		return "waduh asu " + p.TitleInput.View()
	}
	return "This is the playlist"
}
