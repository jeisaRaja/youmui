package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var apikey = os.Getenv("YOUTUBE_API_KEY")

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type TrendingMusicResponse struct {
	Kind  string `json:"kind"`
	Items []struct {
		ID      string `json:"id"`
		Snippet struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			ChannelId   string `json:"channelId"`
		} `json:"snippet"`
	} `json:"items"`
}

func getTrendingMusic(client HTTPClient) (*TrendingMusicResponse, error) {
	endpoint := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=snippet&chart=mostPopular&videoCategoryId=10&regionCode=US&key=%s", apikey)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		e := fmt.Errorf("Error while creating new request: %v", err)
		return nil, e
	}

	res, err := client.Do(req)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err)
	}

	var trendingMusic TrendingMusicResponse
	if err := json.Unmarshal(body, &trendingMusic); err != nil {
		return nil, fmt.Errorf("error unmarshalling the request body: %s", err)
	}

	return &trendingMusic, nil
}

type SearchResult struct {
	Kind  string `json:"kind"`
	Items []struct {
		Kind string `json:"kind"`
		ID   struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		}
	} `json:"items"`
}

func searchWithKeyword(client HTTPClient, keyword string) (*SearchResult, error) {
	endpoint := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&maxResults=10&q=%s&regionCode=US&type=video&videoCategoryId=10&key=%s", keyword, apikey)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err)
	}

	var searchResult SearchResult
	if err := json.Unmarshal(body, &searchResult); err != nil {
		return nil, fmt.Errorf("error unmarshalling the request body: %s", err)
	}

	return &searchResult, nil
}
