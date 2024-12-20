package stream

import "testing"

func TestFetchAndPlayAudio(t *testing.T) {
	url := "https://www.youtube.com/watch?v=K9yaiDG29TM"

  player := NewPlayer()

	err := player.FetchAndPlayAudio(url)
	if err != nil {
		t.Fatalf("Expected no error, %s", err)
	}
}
