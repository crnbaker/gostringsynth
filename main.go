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

func makeAndPlay(pitch float64, sourceOutChan chan sources.Source) {
	envelope := envelopes.NewTriangleEnvelope(attackTime, decayTime, sampleRate)
	s := sources.NewStereoSine(sampleRate, envelope)
	errors.Chk(s.Start())
	s.PlayNote(pitch, 0.1)
	sourceOutChan <- s
}

func main() {
	portaudio.Initialize()

	noteLength := attackTime + decayTime
	returnedSourcesChan := make(chan sources.Source)
	var returnedSources [numVoices]sources.Source

	go makeAndPlay(80, returnedSourcesChan)
	time.Sleep(time.Millisecond * 60)
	go makeAndPlay(160, returnedSourcesChan)
	time.Sleep(time.Millisecond * 60)
	go makeAndPlay(120, returnedSourcesChan)
	time.Sleep(time.Millisecond * 60)
	go makeAndPlay(240, returnedSourcesChan)
	time.Sleep(time.Millisecond * 60)

	time.Sleep(noteLength)

	for i := 0; i < 4; i++ {
		returnedSource := <-returnedSourcesChan
		returnedSources[i] = returnedSource
	}

	shutdown(returnedSources)
}
