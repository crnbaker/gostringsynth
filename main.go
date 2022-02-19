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

	newNoteChan := make(chan float64)
	quitChan := make(chan bool)
	synthFunctionChan := make(chan sources.Voice)

	go mixer.MixController(synthFunctionChan, sampleRate)
	go notedispatcher.NoteDispatcher(newNoteChan, synthFunctionChan, quitChan, sampleRate)
	go keypress.KeyDispatcher(newNoteChan)

	for range quitChan {
		fmt.Println("not quitting")
	}
}
