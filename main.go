package main

import (
	"fmt"

	"github.com/crnbaker/gostringsynth/keypress"
	"github.com/crnbaker/gostringsynth/sources"
	"github.com/crnbaker/gostringsynth/voicecontrol"

	"github.com/gordonklaus/portaudio"
)

const sampleRate = 44100

func main() {
	portaudio.Initialize()

	closeVoiceChan := make(chan sources.Source)
	newNoteChan := make(chan float64)
	quitChan := make(chan bool)

	go voicecontrol.NoteDispatcher(newNoteChan, closeVoiceChan, quitChan, sampleRate)
	go keypress.KeyDispatcher(newNoteChan)

	for range quitChan {
		fmt.Println("not quitting")
	}
}
