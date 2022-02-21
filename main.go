package main

import (
	"sync"

	"github.com/crnbaker/gostringsynth/keypress"
	"github.com/crnbaker/gostringsynth/mixer"
	"github.com/crnbaker/gostringsynth/sources"
	"github.com/crnbaker/gostringsynth/voicedispatcher"
)

const sampleRate = 44100
const voiceLimit = 8

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	voiceChan := make(chan sources.Voice)
	noteChan := make(chan keypress.MidiNote)

	go mixer.MixController(&wg, voiceChan, sampleRate, voiceLimit)
	go voicedispatcher.VoiceDispatcher(&wg, noteChan, voiceChan, sampleRate, "string")
	go keypress.NoteDispatcher(&wg, noteChan)

	wg.Wait()
}
