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

type Player struct {
	ctrl     *beep.Ctrl
	streamer beep.StreamSeekCloser
	format   beep.Format
	isPaused bool
}

func NewPlayer() *Player {
	return &Player{
		isPaused: false,
		ctrl:     &beep.Ctrl{Paused: false},
	}
}

func (p *Player) FetchAndPlayAudio(url string) error {

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
	p.ctrl.Streamer = streamer
	p.streamer = streamer
	p.format = format
	if err != nil {
		return fmt.Errorf("mp3.NewDecoder failed: %v", err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	streamerWithCallback := beep.Callback(func() {
		done <- true
	})

	speaker.Play(beep.Seq(p.ctrl, streamerWithCallback))
	<-done

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("yt-dlp command failed: %v", err)
	}

	if err := convertCmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg command failed: %v", err)
	}

	return nil
}

func (p *Player) PlayPause() {
	speaker.Lock()
	p.ctrl.Paused = !p.ctrl.Paused
	speaker.Unlock()
}
