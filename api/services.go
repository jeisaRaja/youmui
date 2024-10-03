package api

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var apikey = os.Getenv("YOUTUBE_API_KEY")

type IPInfo struct {
	Country string `json:"country"`
}

func GetLocationCode() (string, error) {
	resp, err := http.Get("https://ipinfo.io/json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ipInfo IPInfo
	if err := json.Unmarshal(body, &ipInfo); err != nil {
		return "", err
	}

	return ipInfo.Country, nil
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	*http.Client
	Country string
}

func NewClient() *Client {
	var client Client

	country, err := GetLocationCode()
	if err != nil {
		client.Country = "US"
	} else {
		client.Country = country
	}
	client.Client = &http.Client{}
	return &client
}

type Song struct {
	Title string
	ID    string
	URL   string
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

func GetTrendingMusic(client HTTPClient) (*TrendingMusicResponse, error) {
	endpoint := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=snippet&chart=mostPopular&videoCategoryId=10&regionCode=US&key=%s", apikey)
	if apiClient, ok := client.(*Client); ok {
		endpoint = fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=snippet&chart=mostPopular&videoCategoryId=10&regionCode=%s&key=%s", apiClient.Country, apikey)
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		e := fmt.Errorf("Error while creating new request: %v", err)
		return nil, e
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

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
	Items []*struct {
		Kind string `json:"kind"`
		ID   struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet *struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Url         string
		}
	} `json:"items"`
}

func SearchWithKeyword(client HTTPClient, keyword string, limit int) (*SearchResult, error) {
	keyword = url.QueryEscape(keyword)
	endpoint := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&maxResults=%d&q=%s&regionCode=US&type=video&videoCategoryId=10&key=%s", limit, keyword, apikey)

	if apiClient, ok := client.(*Client); ok {
		endpoint = fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&maxResults=%d&q=%s&regionCode=%s&type=video&videoCategoryId=10&key=%s", limit, keyword, apiClient.Country, apikey)
	}

	file, err := tea.LogToFile("debug.log", "log:\n")
	defer file.Close()
	if err != nil {
		panic("err while opening debug.log")
	}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		file.WriteString(endpoint)
		file.WriteString(" ")
		file.WriteString(keyword)
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

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

	for _, item := range searchResult.Items {
		item.Snippet.Title = html.UnescapeString(item.Snippet.Title)
		item.Snippet.Url = "https://youtube.com/watch?v=" + item.ID.VideoID
	}

	return &searchResult, nil
}
