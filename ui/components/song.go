package components

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaraja/youmui/api"
)

type SongComponent struct {
	song      api.Song
	isHovered bool
}

func NewSong(song api.Song) *SongComponent {
	return &SongComponent{
		song: song,
	}
}

func (sc *SongComponent) Init() tea.Cmd {
	return nil
}

func (sc *SongComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if sc.isHovered && msg.String() == "enter" {
			fmt.Println("Selected song:", sc.song.Title)
		}
	}

	return sc, nil
}

func (sc *SongComponent) View() string {
	return fmt.Sprintf("%s\n", sc.song.Title)
}

func (sc *SongComponent) SetHovered(hovered bool) {
	sc.isHovered = hovered
}
