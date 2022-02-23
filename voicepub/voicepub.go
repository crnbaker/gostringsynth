/*
The voicepub package spawns a new Source and publishes it synthesis function as a Voice, each time a new MIDI
note is received.
*/
package voicepub

import (
	"sync"

	"github.com/crnbaker/gostringsynth/notes"
	"github.com/crnbaker/gostringsynth/sources"
)

// PublishVoices listens for new MIDI note on the noteInChan, spawning a source in a new goroutine for every note received.
func PublishVoices(waitGroup *sync.WaitGroup, noteInChan chan notes.StringMidiNote, voiceSendChan chan sources.Voice,
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
func spawnStringSource(note notes.StringMidiNote, voiceSendChan chan sources.Voice, pluckPlotChan chan []float64, sampleRate float64) {

	physics := sources.StringSettings{note.WaveSpeedMpS, note.DecayTimeS, note.PickupPos}
	pluck := sources.PluckSettings{note.PluckPos, note.PluckWidth, notes.MidiVelocityToAmplitude(note.Velocity)}

	lengthM := sources.FreqToStringLength(notes.MidiPitchToFreq(note.Pitch()), note.WaveSpeedMpS)
	s := sources.NewStringSource(sampleRate, voiceSendChan, lengthM, physics, pluck)
	pluckPlotChan <- s.SoftPluck()
	s.PublishVoice()
}

/* func spawnSineSource(note notes.MidiNote, voiceSendChan chan sources.Voice, sampleRate float64) {
	envelope := envelopes.NewTriangleEnvelope(time.Millisecond*100, time.Millisecond*400, sampleRate)
	s := sources.NewSineVoiceSource(sampleRate, envelope, voiceSendChan)
	s.PublishVoice(notes.MidiPitchToFreq(note.Pitch()), notes.MidiVelocityToAmplitude(note.Velocity))
} */
