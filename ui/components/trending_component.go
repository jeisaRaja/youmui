package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaraja/youmui/api"
)

type TrendingContent struct {
	Songs []*api.Song
}

func NewTrendingContent() *TrendingContent {
	return &TrendingContent{}
}

func (tc *TrendingContent) Init() tea.Cmd {
	return func() tea.Msg {
		file, err := tea.LogToFile("debug.log", "log:\n")
		defer file.Close()
		if err != nil {
			panic("err while opening debug.log")
		}
		res, err := api.GetTrendingMusic(api.NewClient())
		if err != nil {
			file.WriteString(strings.Join([]string{"\n", err.Error()}, ""))
		}
		for _, item := range res.Items {
			url := "https://youtube.com/watch?v=" + item.ID
			tc.Songs = append(tc.Songs, &api.Song{Title: item.Snippet.Title, URL: url})
		}
		return res
	}
}

func (tc *TrendingContent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
	// var cmds []tea.Cmd

	return tc, nil
}

func (tc *TrendingContent) View() string {
	trendingView := ""

	if len(tc.Songs) > 0 {
		trendingView = "\nTrending songs\n\n"
		for _, item := range tc.Songs {
			trendingView += item.Title + "\n" + item.URL + "\n\n"
		}
	}
	return trendingView
}
