package main

import (
	"sync"

	"github.com/crnbaker/gostringsynth/audioengine"
	"github.com/crnbaker/gostringsynth/gui"
	"github.com/crnbaker/gostringsynth/notes"
	"github.com/crnbaker/gostringsynth/sources"
)

const sampleRate = 44100
const voiceLimit = 6

func main() {
	var wg sync.WaitGroup
	wg.Add(4)

	voiceChan := make(chan audioengine.SynthVoice)
	noteChan := make(chan sources.StringNote)
	pluckPlotChan := make(chan []float64)
	userSettingsChannel := make(chan gui.SynthParameters)

	go gui.StartUILoop(&wg, pluckPlotChan, userSettingsChannel)
	go audioengine.ControlVoices(&wg, voiceChan, sampleRate, voiceLimit)          // receives voices, plays sound
	go sources.PublishVoices(&wg, noteChan, voiceChan, pluckPlotChan, sampleRate) // receives notes, sends voices
	go notes.PublishNotes(&wg, noteChan, userSettingsChannel)                     // listens for keypresses, sends notes

	wg.Wait()

}
