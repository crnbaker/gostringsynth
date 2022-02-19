package main

import (
	"fmt"

	"github.com/crnbaker/gostringsynth/keypress"
	"github.com/crnbaker/gostringsynth/mixer"
	"github.com/crnbaker/gostringsynth/sources"
	"github.com/crnbaker/gostringsynth/voicedispatcher"

	"github.com/gordonklaus/portaudio"
)

const sampleRate = 44100

func main() {
	portaudio.Initialize()

	voiceChan := make(chan sources.Voice)
	noteChan := make(chan keypress.MidiNote)
	quitChan := make(chan bool)

	go mixer.MixController(voiceChan, sampleRate)
	go voicedispatcher.VoiceDispatcher(noteChan, voiceChan, quitChan, sampleRate)
	go keypress.NoteDispatcher(noteChan)

	for range quitChan {
		fmt.Println("not quitting")
	}
}
