package ui

import (
	"database/sql"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaraja/youmui/api"
	"github.com/jeisaraja/youmui/storage"
)

type SongListMsg struct {
	Songs []api.Song
}

func SearchSongCallback(input string) tea.Cmd {
	return func() tea.Msg {
		res, err := api.SearchWithKeyword(api.NewClient(), input, 5)
		if err != nil {
			res, err = api.SearchWithKeywordWithoutApi(input)
		}
		return SongListMsg{Songs: res}
	}
}

type CreatePlaylistMsg struct {
	Title string
	ID    int64
}

func CreatePlaylistCallback(db *sql.DB, input string) tea.Cmd {
	return func() tea.Msg {
		res, err := storage.AddPlaylist(db, input)
		if err != nil {
			panic(err)
		}
		return CreatePlaylistMsg{
			Title: input,
			ID:    *res,
		}
	}
}
