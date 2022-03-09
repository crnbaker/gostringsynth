/*
The voicepub package spawns a new Source and publishes it synthesis function as a Voice, each time a new MIDI
note is received.
*/
package voicepub

import (
	"sync"

	"github.com/crnbaker/gostringsynth/sources"
)

type StringNote interface {
	DecayTimeS() float64
	PickupPos() float64
	PluckPos() float64
	PluckWidth() float64
	Amplitude() float64
	LengthM() float64
	Wavespeed() float64
}

// PublishVoices listens for new MIDI note on the noteInChan, spawning a source in a new goroutine for every note received.
func PublishVoices(waitGroup *sync.WaitGroup, noteInChan chan StringNote, voiceSendChan chan sources.Voice,
	pluckPlotChan chan []float64, sampleRate float64) {
	defer waitGroup.Done()
	defer close(voiceSendChan)
	defer close(pluckPlotChan)
	for note := range noteInChan {
		go spawnStringSource(note, voiceSendChan, pluckPlotChan, sampleRate)
	}
}

// spawnSineSource constructs and configrues a finite difference string source, and publishes its Voice
// to the voiceSendChan.
func spawnStringSource(note StringNote, voiceSendChan chan sources.Voice, pluckPlotChan chan []float64, sampleRate float64) {

	physics := sources.StringSettings{
		WaveSpeedMpS: note.Wavespeed(), DecayTimeS: note.DecayTimeS(), PickupPosReStringLen: note.PickupPos(),
	}
	pluck := sources.PluckSettings{
		PosReStrLen: note.PluckPos(), WidthReStrLen: note.PluckWidth(), Amplitude: note.Amplitude(),
	}
	s := sources.NewStringSource(sampleRate, voiceSendChan, note.LengthM(), physics, pluck)
	pluckPlotChan <- s.SoftPluck()
	s.PublishVoice()
}
