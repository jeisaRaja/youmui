package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaraja/youmui/api"
)

type SearchContent struct {
	TextInput    *TextInput
	Results      []string
	SearchResult *api.SearchResult
	SearchBar    bool
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
		SearchBar:    true,
	}
}

func (sc *SearchContent) Init() tea.Cmd {
	return sc.TextInput.Init()
}

func (sc *SearchContent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case *api.SearchResult:
		sc.SearchResult = msg
		sc.SearchBar = false // Hide the text input after getting the result
	default:
		if sc.SearchBar {
			model, cmd := sc.TextInput.Update(msg)
			cmds = append(cmds, cmd)
			if textInputModel, ok := model.(*TextInput); ok {
				sc.TextInput = textInputModel
			}
		}
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

	var inputView string
	if sc.SearchBar {
		file, err := tea.LogToFile("debug.log", "log from search content view:\n")
		defer file.Close()
		if err != nil {
			panic("err while opening debug.log")
		}
		file.WriteString(fmt.Sprintf("the sc.SearchBar is %v\n", sc.SearchBar))
		inputView = sc.TextInput.View()
	}

	return fmt.Sprintf(
		"%s\n%s\n\n",
		inputView,
		resultsView,
	)
}
