package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInput struct {
	tea.Model
	input    textinput.Model
	callback InputCallback
}

type InputCallback func(input string) tea.Cmd

func NewTextInputView(placeholder string, charLimit, width int, callback InputCallback) *TextInput {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = charLimit
	ti.Width = width

	return &TextInput{
		input:    ti,
		callback: callback,
	}
}

func (tm *TextInput) Init() tea.Cmd {
	return textinput.Blink
}

func (tm *TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return tm, tm.callback(tm.input.Value())
		case tea.KeyCtrlC, tea.KeyEsc:
			return tm, tea.Quit
		}
	}
	tm.input, cmd = tm.input.Update(msg)
	return tm, cmd
}

func (tm *TextInput) View() string {
	return fmt.Sprintf(
		"Search for songs\n\n%s\n\n%s",
		tm.input.View(),
		"(esc to quit)",
	) + "\n"
}
