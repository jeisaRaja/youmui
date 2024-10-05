package main

import (
	"log"

	"github.com/jeisaraja/youmui/stream"
	// "github.com/jeisaraja/youmui/ui"
)

func main() {
	url := "https://www.youtube.com/watch?v=PtJ6yAGjsIs"
	if err := stream.FetchAndPlayAudio(url); err != nil {
		log.Fatalf("Error fetching and playing audio: %v", err)
	}
	// ui.Start()
}
