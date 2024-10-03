package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaraja/youmui/api"
	"github.com/jeisaraja/youmui/ui/types"
)

type SearchContent struct {
	TextInput    *TextInput
	Results      []string
	SearchResult *api.SearchResult
}

func NewSearchContent(placeholder string, charLimit, width int) *SearchContent {
	callback := func(input string) tea.Cmd {
		return func() tea.Msg {
			file, err := tea.LogToFile("debug.log", "log:\n")
			defer file.Close()
			if err != nil {
				panic("err while opening debug.log")
			}
			res, err := api.SearchWithKeyword(api.NewClient(), input, 3)
			if err != nil {
				file.WriteString(strings.Join([]string{"\n", err.Error()}, ""))
			}
			return res
		}
	}
	return &SearchContent{
		TextInput:    NewTextInputView(placeholder, charLimit, width, callback),
		Results:      []string{},
		SearchResult: &api.SearchResult{},
	}
}

func (sc *SearchContent) Init() tea.Cmd {
	return sc.TextInput.Init()
}

func (sc *SearchContent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	model, cmd := sc.TextInput.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case *api.SearchResult:
		sc.SearchResult = msg
	case types.Mockres:
		sc.Results = msg
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			sc.TextInput.callback(sc.TextInput.input.Value())
		}
	}

	if textInputModel, ok := model.(*TextInput); ok {
		sc.TextInput = textInputModel
	}
	return sc, tea.Batch(cmds...)
}

func (sc *SearchContent) View() string {
	resultsView := ""

	if sc.SearchResult != nil && len(sc.SearchResult.Items) > 0 {
		resultsView = "\nSearch Results:\n"
		for _, item := range sc.SearchResult.Items {
			resultsView += item.Snippet.Title + "\n" + item.Snippet.Url + "\n\n"
		}
	}

	return fmt.Sprintf(
		"%s\n%s\n\n",
		sc.TextInput.View(),
		resultsView,
	)
}
