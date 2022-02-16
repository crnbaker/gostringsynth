package main

import (
	"time"

	"github.com/crnbaker/gostringsynth/errors"
	"github.com/crnbaker/gostringsynth/sources"

	"github.com/gordonklaus/portaudio"
)

const sampleRate = 44100

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()
	s, err := sources.NewSource("Sine", sampleRate)
	errors.Chk(err)
	defer s.Close()
	errors.Chk(s.Start())

	s.PlayNote(200, 0.2)
	time.Sleep(time.Millisecond * 200)
	s.PlayNote(400, 0.2)
	time.Sleep(time.Millisecond * 200)
	s.PlayNote(300, 0.2)
	time.Sleep(time.Millisecond * 200)
	s.PlayNote(600, 0.2)
	time.Sleep(time.Millisecond * 500)

	errors.Chk(s.Stop())
}
