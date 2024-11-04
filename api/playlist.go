package api

import (
	"encoding/json"
	"fmt"
	"net/url"
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
func FetchPlaylist(url string, source SourceType) (*string, []Song, error) {
	if source == YOUTUBE {
		ytUrl, err := convertToYoutubeURL(url)
		if err != nil {
			return nil, nil, fmt.Errorf("[ERROR] failed to convert url to std youtube url")
		}
		cmd := exec.Command("yt-dlp", "-j", "--flat-playlist", "--compat-options", "no-youtube-unavailable-videos", ytUrl)
		output, err := cmd.Output()
		if err != nil {
			return nil, nil, fmt.Errorf("[error] from FetchPlaylist: %v", err)
		}
		return parsePlaylist(string(output))

	}

	return nil, nil, nil
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

func convertToYoutubeURL(musicURL string) (string, error) {
	parsedURL, err := url.Parse(musicURL)
	if err != nil {
		return "", fmt.Errorf("[ERROR] parsing URL: %w", err)
	}

	playlistID := parsedURL.Query().Get("list")
	if playlistID == "" {
		return "", fmt.Errorf("[ERROR] no playlist ID found in URL")
	}

	youtubeURL := fmt.Sprintf("https://www.youtube.com/playlist?list=%s", playlistID)
	return youtubeURL, nil
}
