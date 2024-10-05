package stream

import (
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func FetchAndPlayAudio(url string) error {
	start := time.Now()

	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "-o", "-", url)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout from yt-dlp: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to run yt-dlp: %v", err)
	}

	downloadTime := time.Since(start)
	fmt.Printf("Time taken to start downloading: %v\n", downloadTime)

	convertCmd := exec.Command("ffmpeg", "-i", "pipe:0", "-acodec", "libmp3lame", "-f", "mp3", "pipe:1")
	convertCmd.Stdin = stdout
	convertedStdout, err := convertCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout from ffmpeg: %v", err)
	}

	conversionTime := time.Since(start) - downloadTime
	fmt.Printf("Time taken to start conversion: %v\n", conversionTime)

	if err := convertCmd.Start(); err != nil {
		return fmt.Errorf("failed to run ffmpeg: %v", err)
	}

	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()
	defer pipeWriter.Close()

	go func() {
		_, err := io.Copy(pipeWriter, convertedStdout)
		if err != nil {
			fmt.Printf("failed to copy ffmpeg output: %v", err)
		}
		pipeWriter.Close()
	}()

	streamer, format, err := mp3.Decode(pipeReader)
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
	fmt.Println("Playback has started.")

	<-done

	fmt.Println("Playback has finished.")

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("yt-dlp command failed: %v", err)
	}
	downloadFinish := time.Now()
	fmt.Printf("Download finished in %v seconds.\n", downloadFinish.Sub(start).Seconds())

	if err := convertCmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg command failed: %v", err)
	}

	conversionFinish := time.Now()
	fmt.Printf("Conversion finished in %v seconds.\n", conversionFinish.Sub(start).Seconds())

	return nil
}

