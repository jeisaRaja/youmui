package components

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jeisaraja/youmui/api"
)

type ResultComponent struct {
	Songs []*SongComponent
	Cur   int
}

func NewResult() *ResultComponent {
	return &ResultComponent{
		Cur: 0,
	}
}

func (rc *ResultComponent) SetSearchResult(results *api.SearchResult) {
	for _, item := range results.Items {
		var song api.Song
		song.ID = item.ID.VideoID
		song.URL = "https://youtube.com/watch?v=" + song.ID
		song.Title = item.Snippet.Title

		songComp := NewSong(song)
		rc.Songs = append(rc.Songs, songComp)
	}
}

func (rc *ResultComponent) Init() tea.Cmd {
	return nil
}

func (rc *ResultComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			rc.Increment()
		case "up":
			rc.Decrement()
		}
	}
	return rc, nil
}

func (rc *ResultComponent) View() string {
	var stringView = ""
	if rc.Songs != nil && len(rc.Songs) > 0 {
		for _, item := range rc.Songs {
			stringView = stringView + item.song.Title + "\n"
		}
	}
	strCur := strconv.Itoa(rc.Cur)
	return stringView + "\n" + strCur
}

func (rc *ResultComponent) Increment() {
	if rc.Cur > len(rc.Songs)-1 {
		rc.Cur = 0
	} else {
		rc.Cur++
	}
}

func (rc *ResultComponent) Decrement() {
	if rc.Cur < 1 {
		rc.Cur = len(rc.Songs) - 1
	} else {
		rc.Cur--
	}
}
