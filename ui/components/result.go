package components

import tea "github.com/charmbracelet/bubbletea"

type ResultComponent struct {
	Songs []*SongComponent
	Cur   int
}

func (rc *ResultComponent) Init() tea.Cmd {
	return nil
}

func (rc *ResultComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return rc, nil
}

func (rc *ResultComponent) View() string {
	return ""
}
