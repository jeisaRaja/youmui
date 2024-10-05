package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jeisaraja/youmui/api"
	"github.com/jeisaraja/youmui/stream"
	// "github.com/jeisaraja/youmui/ui"
)

func main() {
	keyword := flag.String("s", "rick roll", "search for a song")
	flag.Parse()
	client := api.NewClient()
	result, err := api.SearchWithKeyword(client, *keyword, 5)
	if err != nil {
		log.Fatalf("Error getting the songUrl")
	}

	videoID := result.Items[0].ID.VideoID
	videoUrl := "https://youtube.com/watch?v=" + videoID

	fmt.Println(videoUrl)

	if err := stream.FetchAndPlayAudio(videoUrl); err != nil {
		log.Fatalf("Error fetching and playing audio: %v", err)
	}
	// ui.Start()
}
