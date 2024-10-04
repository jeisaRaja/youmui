package stream

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

func FetchAndPlayAudio(url string) error {
	outputFile := "downloaded_audio.mp3"
	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "-o", outputFile, url)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Failed to run yt-dlp: %v", err)
	}

	file, err := os.Open(outputFile)
	if err != nil {
		return fmt.Errorf("Opening audio file failed: %v", err)
	}
	defer file.Close()

	decodeMp3, err := mp3.NewDecoder(file)
	if err != nil {
		return fmt.Errorf("mp3.NewDecoder failed: %v", err)
	}
	file.Close()

	options := oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 2,
	}

	ctx, otoChan, err := oto.NewContext(&options)
	if err != nil {
		return fmt.Errorf("failed to create new context: %v", err)
	}
	<-otoChan

	player := ctx.NewPlayer(decodeMp3)
	player.Play()

	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	err = player.Close()
	if err != nil {
		return fmt.Errorf("failed to close player: %v", err)
	}

	return nil
}
