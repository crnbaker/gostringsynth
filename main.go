package main

import (
	"time"

	"github.com/crnbaker/gostringsynth/errors"
	"github.com/crnbaker/gostringsynth/sources"

	"github.com/gordonklaus/portaudio"
)

const sampleRate = 44100

func shutdown(source sources.Source) {
	errors.Chk(source.Stop())
	source.Close()
	portaudio.Terminate()
}

func main() {
	portaudio.Initialize()
	s, err := sources.NewSource("Sine", sampleRate)
	errors.Chk(err)

	defer shutdown(s)
	errors.Chk(s.Start())

	for i := 0; i < 8; i++ {
		s.PlayNote(80, 0.15, 0.05)
		time.Sleep(time.Millisecond * 200)
		s.PlayNote(160, 0.15, 0.05)
		time.Sleep(time.Millisecond * 200)
		s.PlayNote(120, 0.15, 0.05)
		time.Sleep(time.Millisecond * 200)
		s.PlayNote(240, 0.15, 0.05)
		time.Sleep(time.Millisecond * 200)
	}
}
