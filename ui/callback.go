package ui

import (
	"database/sql"
	"fmt"

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
	Count int64
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
			Count: 0,
		}
	}
}

type SongAddedToPlaylistMsg struct {
	PlaylistID int64
	SongID     string
}

func AddSongToPlaylistCallback(db *sql.DB, playlist int64, song api.Song) tea.Cmd {
	return func() tea.Msg {
		err := storage.AddSongToPlaylist(db, playlist, song)
		if err != nil {
			panic(fmt.Errorf("error when adding a song to a playlist: %w", err))
		}
		return SongAddedToPlaylistMsg{
			PlaylistID: playlist,
			SongID:     song.ID,
		}
	}
}

type PlayPlaylistMsg struct {
	Songs []api.Song
}

func PlayPlaylistCallback(db *sql.DB, pid int64) tea.Cmd {
	return func() tea.Msg {
		songs, err := storage.GetSongsFromPlaylist(db, pid)
		if err != nil {
			panic(fmt.Errorf("error when getting songs from playlist: %w", err))
		}
		return PlayPlaylistMsg{Songs: songs}
	}
}
