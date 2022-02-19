package main

import (
	"fmt"

	"github.com/crnbaker/gostringsynth/keypress"
	"github.com/crnbaker/gostringsynth/mixer"
	"github.com/crnbaker/gostringsynth/notedispatcher"
	"github.com/crnbaker/gostringsynth/sources"

	"github.com/gordonklaus/portaudio"
)

const sampleRate = 44100

func main() {
	portaudio.Initialize()

	noteChan := make(chan keypress.MidiNote)
	quitChan := make(chan bool)
	voiceChan := make(chan sources.Voice)

	go mixer.MixController(voiceChan, sampleRate)
	go notedispatcher.NoteDispatcher(noteChan, voiceChan, quitChan, sampleRate)
	go keypress.KeyDispatcher(noteChan)

	for range quitChan {
		fmt.Println("not quitting")
	}
}
