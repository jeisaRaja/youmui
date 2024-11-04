package ui

import (
	"database/sql"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaraja/youmui/api"
	"github.com/jeisaraja/youmui/storage"
)

type ErrorMsg struct {
	Error error
}

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
			return ErrorMsg{Error: err}
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
			return ErrorMsg{Error: err}
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
			return ErrorMsg{Error: err}
		}
		return PlayPlaylistMsg{Songs: songs}
	}
}

type SongDeleted struct{}

func DeleteSongFromPlaylist(db *sql.DB, pid int64, sid int64) tea.Cmd {
	return func() tea.Msg {
		err := storage.DeleteSongFromPlaylist(db, pid, sid)
		if err != nil {
			return nil
		}
		return SongDeleted{}
	}
}

type PlaylistsData []PlaylistData

type PlaylistData struct {
	Title string
	ID    int64
	Count int64
}

func GetPlaylists(db *sql.DB) tea.Cmd {
	return func() tea.Msg {
		rawData, err := storage.GetPlaylists(db)
		if err != nil {
			return nil
		}

		var data PlaylistsData
		for _, item := range rawData {
			data = append(data, PlaylistData{
				Title: item.Title,
				ID:    item.ID,
				Count: item.Count,
			})
		}

		return data
	}
}

func ImportPlaylistIntoDB(db *sql.DB, url string) tea.Cmd {
	return func() tea.Msg {
		title, songs, err := api.FetchPlaylist(url, api.YOUTUBE)
		if err != nil {
			return ErrorMsg{Error: fmt.Errorf("[ERROR] while calling api.FetchPlaylist: %v", err)}
		}
		if title == nil {
			return ErrorMsg{Error: fmt.Errorf("[ERROR] playlist title is nil")}
		}

		res, err := storage.AddPlaylist(db, *title)
		if err != nil {
			return ErrorMsg{Error: fmt.Errorf("[ERROR] while calling storage.AddPlaylist: %v", err)}
		}

		for _, song := range songs {
			err := storage.AddSongToPlaylist(db, *res, song)
			if err != nil {
				return ErrorMsg{Error: fmt.Errorf("[ERROR] while calling storage.AddSongToPlaylist: %v", err)}
			}
		}

		return CreatePlaylistMsg{
			Title: *title,
			ID:    *res,
			Count: int64(len(songs)),
		}
	}
}
