package api

import (
	"bytes"
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

	res, err := getTrendingMusic(&mockClient)
	if err != nil {
		t.Fatalf("err while get trending music: %s", err)
	}
	if len(res.Items) == 0 {
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

	result, err := searchWithKeyword(&mockClient, "test")
	if err != nil {
		t.Fatalf("err while search with keyword: %s", err)
	}
  if len(result.Items) != 1 || result.Items[0].ID.VideoID != "mockVideoId1" {
		t.Errorf("Expected videoId 'mockVideoId1', got %v", result.Items[0].ID.VideoID)
	}
}
