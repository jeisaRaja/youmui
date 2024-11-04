package api

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFetchPlaylist(t *testing.T) {
	var examplePlaylist = "https://music.youtube.com/playlist?list=PLfsJhwuMhRrB71-Ad-kuiY3HkiYqoTlo_&si=lrVlA7zaWnzLjchm"
	title, _, err := fetchPlaylist(examplePlaylist, YOUTUBE)
	if err != nil {
		t.Fatalf("error when fetching playlist: %v", err)
	}
	if title == nil {
		t.Fatalf("error when fetching playlist: %v", err)
	}
}

func TestParsePlaylist(t *testing.T) {
	mockOutput := ` {"title": "Song One", "url": "link1", "playlist_title": "cool songs"}
                  {"title": "Song Two", "url": "link2", "playlist_title": "cool songs"}
                  {"title": "Song Three", "url": "link3", "playlist_title": "cool songs"}`

	expectedSongs := []Song{
		{Title: "Song One", URL: "link1"},
		{Title: "Song Two", URL: "link2"},
		{Title: "Song Three", URL: "link3"},
	}

	title, outSongs, err := parsePlaylist(mockOutput)
	if err != nil {
		t.Fatalf("failed to parse the playlist string: %v", err)
	}

	if title == nil {
		t.Fatalf("failed to get the playlist title: %v", err)
	}

	fmt.Printf("Playlist title: %s\n", *title)

	for i := range expectedSongs {
		if !reflect.DeepEqual(expectedSongs[i], outSongs[i]) {
			t.Fatalf("song and expected song are different: expected %v, got %v", expectedSongs[i], outSongs[i])
		}
	}
}
