package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type SearchContent struct {
	TextInput *TextInput
	Results   []string // This will hold the search results
}

func NewSearchContent(placeholder string, charLimit, width int) *SearchContent {
	callback := func(input string) tea.Cmd {
		return func() tea.Msg {
			results := []string{"Song1", "Song2", "Song3"}
			return NewBaseView("results", strings.Join(results, ""), "end of result")
		}
	}
	return &SearchContent{
		TextInput: NewTextInputView(placeholder, charLimit, width, callback),
		Results:   []string{}, 
	}
}

func (sc *SearchContent) Init() tea.Cmd {
	return sc.TextInput.Init()
}

func (sc *SearchContent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	model, cmd := sc.TextInput.Update(msg) // Get the updated model and command

	// Type assertion to convert the tea.Model back to *TextInput
	if textInputModel, ok := model.(*TextInput); ok {
		sc.TextInput = textInputModel // Assign it back to sc.TextInput
	} else {
		// Handle the error or unexpected type case if needed
	}

	return sc, cmd
}

func (sc *SearchContent) View() string {
	// Join the text input view and the results vertically
	resultsView := ""
	if len(sc.Results) > 0 {
		resultsView = "\nSearch Results:\n" + fmt.Sprintf("%v", sc.Results)
	}

	return fmt.Sprintf(
		"%s\n%s\n\n",
		sc.TextInput.View(),
		resultsView,
	)
}
