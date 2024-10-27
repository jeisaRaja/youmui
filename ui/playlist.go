package ui

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jeisaraja/youmui/api"
	"github.com/jeisaraja/youmui/storage"
)

type PlaylistComponent struct {
	playlists    []PlaylistData
	cur          int
	openPlaylist *SongList
	isOpen       bool
	db           *sql.DB
}

type InputStateMsg struct{}

func NewPlaylistComponent(db *sql.DB) *PlaylistComponent {
	return &PlaylistComponent{db: db, isOpen: false, openPlaylist: NewSongList(db)}
}

func (p *PlaylistComponent) Init() tea.Cmd {
	return textinput.Blink
}

func (p *PlaylistComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case PlaylistsData:
		p.RefreshPlaylist(msg)
		return p, nil
	case SongDeleted:
		if p.isOpen {
			return p, GetAllSongsFromPlaylist(p.db, p.GetCurrent())
		}
	case SongsFromPlaylist:
		p.openPlaylist.UpdateSongs(msg.songs)
		return p, GetPlaylists(p.db)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			p.isOpen = false
			return p, nil
		case "a":
			if !p.isOpen {
				return p, PlayPlaylistCallback(p.db, p.GetCurrent())
			} else {
				_, cmd := p.openPlaylist.Update(msg)
				return p, cmd
			}
		case "d":
			if p.isOpen {
				song := p.openPlaylist.GetSelectedSong()
				return p, DeleteSongFromPlaylist(p.db, p.GetCurrent(), song.DbID)
			}
		case "enter":
			if !p.isOpen {
				p.isOpen = true
				return p, GetAllSongsFromPlaylist(p.db, p.playlists[p.cur].ID)
			} else {
				_, cmd := p.openPlaylist.Update(msg)
				return p, cmd
			}

		case "j", "down":
			if p.isOpen {
				_, cmd := p.openPlaylist.Update(msg)
				return p, cmd
			}
			if p.cur < len(p.playlists)-1 {
				p.cur++
			}
		case "k", "up":
			if p.isOpen {
				_, cmd := p.openPlaylist.Update(msg)
				return p, cmd
			}
			if p.cur > 0 {
				p.cur--
			}
		default:
			if p.isOpen {
				p.openPlaylist.Update(msg)
			}
		}
	}
	return p, nil
}

func (p *PlaylistComponent) View() string {
	var builder strings.Builder
	builder.WriteString("\n\nPlaylists:\n")

	if p.isOpen {
		return p.openPlaylist.View()
	}

	normalStyle := lipgloss.NewStyle()
	hoverStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("75"))

	for i, playlist := range p.playlists {
		style := normalStyle
		if i == p.cur {
			style = hoverStyle
		}
		builder.WriteString(fmt.Sprintf("%d. %s (%d songs)\n", i+1, style.Render(playlist.Title), playlist.Count))
	}

	return lipgloss.NewStyle().PaddingTop(2).Render(builder.String())
}

func (p *PlaylistComponent) AppendPlaylist(title string, id int64, count int64) {
	p.playlists = append(p.playlists, PlaylistData{Title: title, ID: id, Count: count})
}

func (p *PlaylistComponent) RefreshPlaylist(data PlaylistsData) {
	p.playlists = data
}

func (p *PlaylistComponent) SetPlaylists(data PlaylistsData) {
	for _, ps := range data {
		p.AppendPlaylist(ps.Title, ps.ID, ps.Count)
	}
}

func (p *PlaylistComponent) IncrementCount(pid int64) {
	for i := range p.playlists {
		if p.playlists[i].ID == pid {
			p.playlists[i].Count++
		}
	}
}

func (p *PlaylistComponent) GetCurrent() int64 {
	if len(p.playlists) > 0 {
		return p.playlists[p.cur].ID
	}
	return 0
}

type SongsFromPlaylist struct {
	songs []api.Song
}

func GetAllSongsFromPlaylist(db *sql.DB, pid int64) tea.Cmd {
	return func() tea.Msg {
		songs, err := storage.GetSongsFromPlaylist(db, pid)
		if err != nil {
			panic("error when getting songs from playlist")
		}
		return SongsFromPlaylist{songs}
	}
}
