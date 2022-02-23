package main

import (
	"sync"

	"github.com/crnbaker/gostringsynth/audioengine"
	"github.com/crnbaker/gostringsynth/notes"
	"github.com/crnbaker/gostringsynth/plotting"
	"github.com/crnbaker/gostringsynth/sources"
	"github.com/crnbaker/gostringsynth/voicepub"
)

const sampleRate = 44100
const voiceLimit = 8

func main() {
	var wg sync.WaitGroup
	wg.Add(4)

	voiceChan := make(chan sources.Voice)
	noteChan := make(chan notes.MidiNote)
	pluckPlotChan := make(chan []float64)

	go plotting.StartUI(&wg, pluckPlotChan)
	go audioengine.ControlVoices(&wg, voiceChan, sampleRate, voiceLimit)
	go voicepub.PublishVoices(&wg, noteChan, voiceChan, pluckPlotChan, sampleRate, "string")
	go notes.PublishNotes(&wg, noteChan)

	wg.Wait()

}
