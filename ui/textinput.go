package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputUsage int

const (
	SearchSong = iota
	SearchPlaylist
	CreatePlaylist
	ImportPlaylist
)

type TextInput struct {
	tea.Model
	Input textinput.Model
	Usage InputUsage
}

func NewTextInputView(charLimit, width int) *TextInput {

	ti := textinput.New()
	ti.CharLimit = charLimit
	ti.Width = width

	return &TextInput{
		Input: ti,
	}
}

func (tm *TextInput) Init() tea.Cmd {
	tm.Input.Focus()
	return textinput.Blink
}

func (tm *TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	tm.Input, cmd = tm.Input.Update(msg)
	cmds = append(cmds, cmd)
	return tm, tea.Batch(cmds...)
}

func (tm *TextInput) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, tm.Input.View())
}

func (tm *TextInput) SetPlaceholder(p string) *TextInput {
	tm.Input.Placeholder = p
	return tm
}

func (tm *TextInput) SetUsage(usage InputUsage) {
	tm.Usage = usage
}
