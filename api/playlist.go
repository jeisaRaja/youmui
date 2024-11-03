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

func fetchPlaylist(url string, source SourceType) error {
	var songs []Song

	if source == YOUTUBE {
		cmd := exec.Command("yt-dlp", "-j", "--flat-playlist", url)
		output, err := cmd.Output()
		if err != nil {
			return err
		}

		// Split by newline
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			var song Song

			if err := json.Unmarshal([]byte(line), &song); err != nil {
				return fmt.Errorf("error unmarshalling song: %v", err)
			}

			songs = append(songs, song)
		}

		fmt.Printf("Songs are: %#v", songs)
	}

	return nil
}
