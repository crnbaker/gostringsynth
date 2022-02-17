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
const decayTime = time.Millisecond * 2000
const numVoices = 4

func shutdown(oscillators [numVoices]sources.Source) {
	time.Sleep(time.Millisecond * 450)
	for _, s := range oscillators {
		errors.Chk(s.Stop())
		s.Close()
	}
	portaudio.Terminate()
}

func main() {
	portaudio.Initialize()

	var oscillators [numVoices]sources.Source
	for i := 0; i < numVoices; i++ {
		envelope := envelopes.NewTriangleEnvelope(attackTime, decayTime, sampleRate)
		s := sources.NewStereoSine(sampleRate, envelope)
		errors.Chk(s.Start())
		oscillators[i] = s
	}
	defer shutdown(oscillators)
	noteLength := attackTime + decayTime

	go oscillators[0].PlayNote(80, 0.1)
	time.Sleep(time.Millisecond * 60)
	go oscillators[1].PlayNote(160, 0.1)
	time.Sleep(time.Millisecond * 60)
	go oscillators[2].PlayNote(120, 0.1)
	time.Sleep(time.Millisecond * 60)
	go oscillators[3].PlayNote(240, 0.1)
	time.Sleep(time.Millisecond * 60)

	time.Sleep(noteLength)
}
