package stream

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func FetchAndPlayAudio(url string) error {
	outputFile := "downloaded_audio.mp3"
	convertedFile := "converted_audio.mp3"
	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "-o", outputFile, url)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Failed to run yt-dlp: %v", err)
	}

	convertCmd := exec.Command("ffmpeg", "-i", outputFile, "-acodec", "libmp3lame", "-b:a", "192k", convertedFile)
	if err := convertCmd.Run(); err != nil {
		return fmt.Errorf("failed to convert mp3: %v", err)
	}

	file, err := os.Open(convertedFile)
	if err != nil {
		return fmt.Errorf("Opening audio file failed: %v", err)
	}
	defer file.Close()

	streamer, format, err := mp3.Decode(file)
	if err != nil {
		return fmt.Errorf("mp3.NewDecoder failed: %v", err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)

	streamerWithCallback := beep.Callback(func() {
		done <- true
	})
	speaker.Play(beep.Seq(streamer, streamerWithCallback))

	<-done

	rmMp3FilesCmd := exec.Command("rm", outputFile, convertedFile)
	if err := rmMp3FilesCmd.Run(); err != nil {
		return fmt.Errorf("failed to remove mp3 files: %v", err)
	}
	return nil
}
