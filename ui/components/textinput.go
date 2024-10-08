package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaraja/youmui/ui/types"
)

type TextInput struct {
	tea.Model
	input     textinput.Model
	callback  InputCallback
	searchFor string
}

type InputCallback func(input string) tea.Cmd

func NewTextInputView(searchFor string, charLimit, width int, callback InputCallback) *TextInput {
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.CharLimit = charLimit
	ti.Width = width

	return &TextInput{
		input:     ti,
		callback:  callback,
		searchFor: searchFor,
	}
}

func (tm *TextInput) Init() tea.Cmd {
	tm.input.Focus()
	return textinput.Blink
}

func (tm *TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case types.FocusMsg:
		if msg.Level == types.ContentFocus {
			if keyMsg, ok := msg.Msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
				return tm, tm.callback(tm.input.Value())
			}
		}
	}

	tm.input, cmd = tm.input.Update(msg)
	cmds = append(cmds, cmd)
	return tm, tea.Batch(cmds...)
}

func (tm *TextInput) View() string {
	return fmt.Sprintf(
		"Search for %s\n\n%s\n\n%s",
		tm.searchFor,
		tm.input.View(),
		"(esc to quit)",
	) + "\n"
}
