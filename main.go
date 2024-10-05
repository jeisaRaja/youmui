package main

import (
	"flag"
	"log"

	"github.com/jeisaraja/youmui/stream"
	// "github.com/jeisaraja/youmui/ui"
)

func main() {
  url := flag.String("s", "https://www.youtube.com/watch?v=PtJ6yAGjsIs", "youtube video url")
  flag.Parse()
	if err := stream.FetchAndPlayAudio(*url); err != nil {
		log.Fatalf("Error fetching and playing audio: %v", err)
	}
	// ui.Start()
}
