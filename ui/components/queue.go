package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaraja/youmui/api"
)

func NewQueue() *Queue {
	return &Queue{
		Songs: []api.Song{},
	}
}

type Queue struct {
	Songs       []api.Song
	PlayingSong *api.Song
}

func (q *Queue) AddToQueue(song api.Song) {
	q.Songs = append(q.Songs, song)
}

func (q *Queue) RemoveFromQueue() api.Song {
	song := q.Songs[0]
	q.Songs = q.Songs[1:]
	return song
}

func (q *Queue) Init() tea.Cmd {
	return nil
}

func (q *Queue) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return q, nil
}

func (q *Queue) View() string {
	var view string
	if q.PlayingSong != nil {
    view += "Playing:\n\n" + q.PlayingSong.Title + "\n"
	}
	view += "\n\nNext:\n\n"
	for i, song := range q.Songs {
		view += fmt.Sprintf("%d. %s\n", i+1, song.Title)
	}
	return view
}

func (q *Queue) Clear() {
	q.Songs = []api.Song{}
}

func (q *Queue) SetPlayingSong(song api.Song) {
	q.PlayingSong = &song
}
