package api

import "testing"

func TestFetchPlaylist(t *testing.T) {
	var examplePlaylist = "https://music.youtube.com/playlist?list=PLfsJhwuMhRrB71-Ad-kuiY3HkiYqoTlo_&si=lrVlA7zaWnzLjchm"
	err := fetchPlaylist(examplePlaylist, YOUTUBE)
	if err != nil {
		t.Fatalf("error when fetching playlist: %v", err)
	}
}
