package stream

import (
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type Player struct {
	ctrl      *beep.Ctrl
	streamer  beep.StreamSeekCloser
	format    beep.Format
	isPaused  bool
	volume    *effects.Volume
	ytdlpCmd  *exec.Cmd
	ffmpegCmd *exec.Cmd
}

func NewPlayer() *Player {
	return &Player{
		isPaused: false,
		ctrl:     &beep.Ctrl{Paused: false},
		volume: &effects.Volume{
			Base:   2,
			Volume: 0,
			Silent: false,
		},
	}
}

func (p *Player) terminatePreviousCommands() {
	if p.ytdlpCmd != nil {
		if err := p.ytdlpCmd.Process.Kill(); err != nil {
			fmt.Printf("Failed to kill yt-dlp: %v\n", err)
		}
		if err := p.ytdlpCmd.Wait(); err != nil && isSignalKilled(err) {
			fmt.Printf("yt-dlp command failed: %v\n", err)
		}
	}

	if p.ffmpegCmd != nil {
		if err := p.ffmpegCmd.Process.Kill(); err != nil {
			fmt.Printf("Failed to kill ffmpeg: %v\n", err)
		}
		if err := p.ffmpegCmd.Wait(); err != nil && isSignalKilled(err) {
			fmt.Printf("ffmpeg command failed: %v\n", err)
		}
	}
}

func (p *Player) FetchAndPlayAudio(url string) error {
	done := make(chan bool)
	p.terminatePreviousCommands()
	p.ytdlpCmd = exec.Command("yt-dlp", "-f", "bestaudio", "-o", "-", url)

	stdout, err := p.ytdlpCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout from yt-dlp: %v", err)
	}

	if err := p.ytdlpCmd.Start(); err != nil {
		return fmt.Errorf("failed to run yt-dlp: %v", err)
	}

	p.ffmpegCmd = exec.Command("ffmpeg", "-i", "pipe:0", "-acodec", "libmp3lame", "-f", "mp3", "pipe:1")
	p.ffmpegCmd.Stdin = stdout
	convertedStdout, err := p.ffmpegCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout from ffmpeg: %v", err)
	}

	if err := p.ffmpegCmd.Start(); err != nil {
		return fmt.Errorf("failed to run ffmpeg: %v", err)
	}

	type teeReadCloser struct {
		io.Reader
		io.Closer
	}
	teeReader := io.TeeReader(convertedStdout, io.Discard)
	readCloser := teeReadCloser{
		Reader: teeReader,
		Closer: convertedStdout,
	}

	streamer, format, err := mp3.Decode(readCloser)

	p.ctrl.Streamer = streamer
	p.streamer = streamer
	p.format = format
	if err != nil {
		return fmt.Errorf("mp3.NewDecoder failed: %v", err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	streamerWithCallback := beep.Callback(func() {
		done <- true
	})

	p.volume.Streamer = p.ctrl

	speaker.Play(beep.Seq(p.volume, streamerWithCallback))
	<-done

	if err := p.ytdlpCmd.Wait(); err != nil {
		return fmt.Errorf("yt-dlp command failed: %v", err)
	}

	if err := p.ffmpegCmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg command failed: %v", err)
	}

	return nil
}

func (p *Player) PlayPause() {
	speaker.Lock()
	p.ctrl.Paused = !p.ctrl.Paused
	speaker.Unlock()
}

func (p *Player) VolumeUp() {
	speaker.Lock()
	p.volume.Volume += 0.5
	speaker.Unlock()
}

func (p *Player) VolumeDown() {
	speaker.Lock()
	p.volume.Volume -= 0.5
	speaker.Unlock()
}

func isSignalKilled(err error) bool {
	if exitError, ok := err.(*exec.ExitError); ok {
		if exitError.ProcessState.ExitCode() == 137 {
			return true
		}
	}
	return false
}
