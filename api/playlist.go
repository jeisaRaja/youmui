package api

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type SourceType string

const (
	YOUTUBE SourceType = "YOUTUBE"
	SPOTIFY SourceType = "SPOTIFY"
)

// Uses yt-dlp’s flat-playlist mode to efficiently gather minimal metadata from YouTube playlists,
// trading detailed info for speed and low resource use, as only basic song data is needed.
func fetchPlaylist(url string, source SourceType) ([]Song, error) {
	var output []byte
	var err error
	if source == YOUTUBE {
		cmd := exec.Command("yt-dlp", "-j", "--flat-playlist", url)
		output, err = cmd.Output()
		fmt.Println(string(output))
		if err != nil {
			return nil, err
		}
	}

	return parsePlaylist(string(output))
}

func parsePlaylist(output string) ([]Song, error) {
	var songs []Song
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		var song Song

		if err := json.Unmarshal([]byte(line), &song); err != nil {
			return nil, fmt.Errorf("error unmarshalling song: %v", err)
		}

		songs = append(songs, song)
	}

	return songs, nil
}
