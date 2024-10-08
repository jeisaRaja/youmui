package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInput struct {
	tea.Model
	Input     textinput.Model
	callback  InputCallback
}

type InputCallback func(input string) tea.Cmd

func NewTextInputView(charLimit, width int, callback InputCallback) *TextInput {
	file, err := tea.LogToFile("debug.log", "log:\n")
	defer file.Close()
	if err != nil {
		panic("err while opening debug.log")
	}
  file.WriteString("writing this when init text input")
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.CharLimit = charLimit
	ti.Width = width

	return &TextInput{
		Input:     ti,
		callback:  callback,
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
	return fmt.Sprintf(
		"Search for \n\n%s\n\n%s",
		tm.Input.View(),
		"(esc to quit)",
	) + "\n"
}
