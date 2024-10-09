package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestGetTrendingMusic(t *testing.T) {

	mockResponse := `{
  "kind": "youtube#videoListResponse",
  "etag": "_iGkOxUBzeohaATmNcnBc2ZkiSg",
  "items": [
    {
      "kind": "youtube#video",
      "etag": "WWeRyzQXMDheSS_M_uhFRPJXZbQ",
      "id": "V9PVRfjEBTI",
      "snippet": {
        "publishedAt": "2024-09-27T15:00:38Z",
        "channelId": "UCDGmojLIoWpXok597xYo8cg",
        "title": "Billie Eilish - BIRDS OF A FEATHER (Official Music Video)",
        "description": "Listen to HIT ME HARD AND SOFT: http://BillieEilish.lnk.to/HITMEHARDANDSOFT\nDownload BIRDS OF A FEATHER Live from Billie’s Amazon Music Songline performance: https://billieeilish.lnk.to/BIRDSOFAFEATHER-AMAZONDOWNLOAD \nGet tickets: https://BillieEilish.lnk.to/TourDates \n\nFollow Billie Eilish: \nTikTok: https://BillieEilish.lnk.to/TikTok \nInstagram: https://BillieEilish.lnk.to/Instagram \nFacebook: https://BillieEilish.lnk.to/Facebook \nTwitter: https://BillieEilish.lnk.to/Twitter \nYouTube: / billieeilish \nWhatsApp: https://BillieEilish.lnk.to/WhatsApp \nEmail: https://BillieEilish.lnk.to/SignUp \nStore: https://BillieEilish.lnk.to/Store \nCell: +1 (310) 807-3956\n\nMusic video by Billie Eilish performing BIRDS OF A FEATHER. © 2024 Darkroom/Interscope Records"
      }
    }
  ]
}
`
	mockClient := MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
			}, nil
		},
	}

	res, err := GetTrendingMusic(&mockClient)
	if err != nil {
		t.Fatalf("err while get trending music: %s", err)
	}
	if len(res) == 0 {
		t.Fatal("res.Items have a length of 0")
	}
}

func TestSearchKeyword(t *testing.T) {
	mockResponse := `{
		"kind": "youtube#searchListResponse",
		"items": [
			{
				"kind": "youtube#searchResult",
				"id": {
					"videoId": "mockVideoId1"
				},
				"snippet": {
					"title": "Test Video 1",
					"description": "Test Description 1"
				}
			}
		]
	}`

	mockClient := MockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
		}, nil
	}}

	result, err := SearchWithKeyword(&mockClient, "test", 10)
	if err != nil {
		t.Fatalf("err while search with keyword: %s", err)
	}
	if len(result) != 1 || result[0].ID != "mockVideoId1" {
		t.Errorf("Expected videoId 'mockVideoId1', got %v", result[0].ID)
	}
}

func TestSearchWithKeywordWithoutApi(t *testing.T) {
	res, err := SearchWithKeywordWithoutApi("every summertime")
  fmt.Println("the result:")
	fmt.Println(res)
	if len(res) != 5 {
		t.Fatalf("length of songs not 5 but: %d", len(res))
	}
	if err != nil {
		t.Fatalf("failed because an error occurs: %v", err)
	}
}

func TestParseSong(t *testing.T) {
	line := "Niki ride on children's car and stuck in the ground Vlad tows on the tractor https://www.youtube.com/watch?v=a7Oh4dKDDuU"
	expectedTitle := "Niki ride on children's car and stuck in the ground Vlad tows on the tractor"
	expectedUrl := "https://www.youtube.com/watch?v=a7Oh4dKDDuU"

	title, url := parseSong(line)

	if title != expectedTitle {
		t.Fatalf("title is not the same as expected: %s, the title from parseSong: %s", expectedTitle, title)
	}

	if url != expectedUrl {
		t.Fatalf("url is not the same as expected: %s, the url from parseSong: %s", expectedUrl, url)
	}
}
