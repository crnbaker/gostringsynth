/*
The voicepub package spawns a new Source and publishes it synthesis function as a Voice, each time a new MIDI
note is received.
*/
package sources

import (
	"sync"

	"github.com/crnbaker/gostringsynth/audioengine"
	"github.com/crnbaker/gostringsynth/excitors"
)

type StringNote interface {
	DecayTimeS() float64
	PickupPos() float64
	PluckPos() float64
	PluckWidth() float64
	Amplitude() float64
	LengthM() float64
	WavespeedMpS() float64
}

// PublishVoices listens for new MIDI note on the noteInChan, spawning a source in a new goroutine for every note received.
func PublishVoices(waitGroup *sync.WaitGroup, noteInChan chan StringNote, voiceSendChan chan audioengine.SynthVoice,
	pluckPlotChan chan []float64, sampleRate float64) {
	defer waitGroup.Done()
	defer close(voiceSendChan)
	defer close(pluckPlotChan)
	for note := range noteInChan {
		physics := stringSettings{
			WaveSpeedMpS: note.WavespeedMpS(), DecayTimeS: note.DecayTimeS(), PickupPosReStringLen: note.PickupPos(),
		}
		pluck := excitors.NewFTDTPluck(
			note.PluckPos(), note.PluckWidth(), note.Amplitude(),
		)
		s := newStringSource(sampleRate, note.LengthM(), physics, &pluck)
		// pluckPlotChan <- s.pluckString()
		s.pluckString()
		pluckPlotChan <- make([]float64, s.numSpatialSections+1)
		voiceSendChan <- s.createVoice()
	}
}
