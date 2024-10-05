package stream

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type bytesReadCloser struct {
	*bytes.Reader
}

func (b *bytesReadCloser) Close() error {
	return nil
}

func FetchAndPlayAudio(url string) error {
	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "-o", "-", url)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout from yt-dlp: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to run yt-dlp: %v", err)
	}

	convertCmd := exec.Command("ffmpeg", "-i", "pipe:0", "-acodec", "libmp3lame", "-f", "mp3", "pipe:1")
	convertCmd.Stdin = stdout
	convertedStdout, err := convertCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout from ffmpeg: %v", err)
	}

	if err := convertCmd.Start(); err != nil {
		return fmt.Errorf("failed to run ffmpeg: %v", err)
	}

	var buf bytes.Buffer
	go func() {
		_, err := io.Copy(&buf, convertedStdout)
		if err != nil {
			fmt.Printf("failed to copy ffmpeg output: %v", err)
		}
	}()

	if err := convertCmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg command failed: %v", err)
	}

	reader := &bytesReadCloser{bytes.NewReader(buf.Bytes())}

	streamer, format, err := mp3.Decode(reader)
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

	return nil
}
