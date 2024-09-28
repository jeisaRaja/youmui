package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewBaseView(header, data, footer string) *BaseView {
	return &BaseView{
		header: header,
		data:   data,
		footer: footer,
	}
}

type BaseView struct {
	header string
	footer string
	data   string
	ContentModel
}

func (v *BaseView) Init() tea.Cmd {
	return nil
}

func (v *BaseView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return v, nil
}

func (v *BaseView) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		v.header,
		v.data,
		v.footer,
	)
}

