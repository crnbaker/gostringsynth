package main

import (
	"time"

	"github.com/crnbaker/gostringsynth/envelopes"
	"github.com/crnbaker/gostringsynth/errors"
	"github.com/crnbaker/gostringsynth/sources"

	"github.com/gordonklaus/portaudio"
)

const sampleRate = 44100
const attackTime = time.Millisecond * 10
const decayTime = time.Millisecond * 190

func shutdown(source sources.Source) {
	time.Sleep(time.Millisecond * 250)
	errors.Chk(source.Stop())
	source.Close()
	portaudio.Terminate()
}

func main() {
	portaudio.Initialize()

	envelope := envelopes.NewTriangleEnvelope(attackTime, decayTime, sampleRate)
	s := sources.NewStereoSine(sampleRate, envelope)
	noteLength := attackTime + decayTime

	defer shutdown(s)
	errors.Chk(s.Start())

	for i := 0; i < 4; i++ {
		s.PlayNote(80*float64(i+1)*0.9, 0.5)
		time.Sleep(noteLength)
		s.PlayNote(160*float64(i+1)*0.9, 0.6)
		time.Sleep(noteLength)
		s.PlayNote(120*float64(i+1)*0.9, 0.7)
		time.Sleep(noteLength)
		s.PlayNote(240*float64(i+1)*0.9, 0.8)
		time.Sleep(noteLength)
	}
}
