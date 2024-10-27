package ui

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
	Select      int
	isPlaying   bool
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

func (q *Queue) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return q, nil
}

func (q *Queue) View() string {
	var view string
	var sign string
	sign = "⏸ "
	if q.isPlaying {
		sign = "⏵ "
	}
	if q.PlayingSong != nil {
		view += "\n\nPlaying:\n\n" + sign + q.PlayingSong.Title + "\n"
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
	q.isPlaying = true
	q.PlayingSong = &song
}

func (q *Queue) PlayPause() {
	q.isPlaying = !q.isPlaying
}

func (q *Queue) AddManyToQueue(songs []api.Song) {
	q.Clear()
	for _, song := range songs {
		q.AddToQueue(song)
	}
}
