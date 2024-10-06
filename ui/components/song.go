package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaraja/youmui/api"
)

type SongComponent struct {
	api.Song
}

func (sc *SongComponent) Init() tea.Cmd {
	return nil
}

func (sc *SongComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return sc, nil
}

func (sc *SongComponent) View() string {
	title := sc.Title + "\n"
	return title
}
