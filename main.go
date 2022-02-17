package main

import (
	"time"

	"github.com/crnbaker/gostringsynth/envelopes"
	"github.com/crnbaker/gostringsynth/errors"
	"github.com/crnbaker/gostringsynth/sources"

	"github.com/gordonklaus/portaudio"
)

const sampleRate = 44100

func shutdown(source sources.Source) {
	time.Sleep(time.Millisecond * 250)
	errors.Chk(source.Stop())
	source.Close()
	portaudio.Terminate()
}

func main() {
	portaudio.Initialize()

	envelope := envelopes.NewTriangleEnvelope(0.1, 0.1, sampleRate)
	s := sources.NewStereoSine(sampleRate, envelope)

	defer shutdown(s)
	errors.Chk(s.Start())

	for i := 0; i < 4; i++ {
		s.PlayNote(80 * float64(i+1) * 0.9)
		time.Sleep(time.Millisecond * 200)
		s.PlayNote(160 * float64(i+1) * 0.9)
		time.Sleep(time.Millisecond * 200)
		s.PlayNote(120 * float64(i+1) * 0.9)
		time.Sleep(time.Millisecond * 200)
		s.PlayNote(240 * float64(i+1) * 0.9)
		time.Sleep(time.Millisecond * 200)
	}
}
