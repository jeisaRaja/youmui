package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TextInput struct {
	tea.Model
	Input    textinput.Model
	callback InputCallback
}

type InputCallback func(input string) tea.Cmd

func NewTextInputView(charLimit, width int, callback InputCallback) *TextInput {

	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.CharLimit = charLimit
	ti.Width = width

	return &TextInput{
		Input:    ti,
		callback: callback,
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
