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
func fetchPlaylist(url string, source SourceType) (*string, []Song, error) {
	var output []byte
	var err error
	if source == YOUTUBE {
		cmd := exec.Command("yt-dlp", "-j", "--flat-playlist", url)
		output, err = cmd.Output()
		fmt.Println(string(output))
		if err != nil {
			return nil, nil, err
		}
	}

	return parsePlaylist(string(output))
}

func parsePlaylist(output string) (*string, []Song, error) {
	var songs []Song
	var titleData struct {
		PlaylistTitle string `json:"playlist_title"`
	}
	lines := strings.Split(output, "\n")
	if len(lines) > 0 {
		firstLine := lines[0]

		err := json.Unmarshal([]byte(firstLine), &titleData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshall playlist_title: %v", err)
		}
	}
	for _, line := range lines {
		if line == "" {
			continue
		}
		var song Song

		if err := json.Unmarshal([]byte(line), &song); err != nil {
			return nil, nil, fmt.Errorf("error unmarshalling song: %v", err)
		}

		songs = append(songs, song)
	}

	return &titleData.PlaylistTitle, songs, nil
}
