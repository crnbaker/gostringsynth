/*
The voicepub package spawns a new Source and publishes it synthesis function as a Voice, each time a new MIDI
note is received.
*/
package voicepub

import (
	"sync"
	"time"

	"github.com/crnbaker/gostringsynth/envelopes"
	"github.com/crnbaker/gostringsynth/notes"
	"github.com/crnbaker/gostringsynth/sources"
)

// PublishVoices listens for new MIDI note on the noteInChan, spawning a source in a new goroutine for every note received.
func PublishVoices(waitGroup *sync.WaitGroup, noteInChan chan notes.MidiNote, voiceSendChan chan sources.Voice, sampleRate float64, osc string) {
	defer waitGroup.Done()
	for note := range noteInChan {
		if osc == "sine" {
			go spawnSineSource(note, voiceSendChan, sampleRate)
		} else if osc == "string" {
			go spawnStringSource(note, voiceSendChan, sampleRate)
		}
	}
	close(voiceSendChan)
}

// spawnSineSource constructs and configrues a finite difference string source, and publishes its Voice
// to the voiceSendChan.
func spawnStringSource(note notes.MidiNote, voiceSendChan chan sources.Voice, sampleRate float64) {
	const waveSpeedMpS = 200
	const pickupPos = 0.5
	const decayTimeS = 3.0
	lengthM := sources.FreqToStringLength(notes.MidiPitchToFreq(note.Pitch), waveSpeedMpS)
	s := sources.NewStringSource(sampleRate, voiceSendChan, lengthM, waveSpeedMpS, pickupPos, decayTimeS)
	s.PublishVoice(notes.MidiPitchToFreq(note.Pitch), notes.MidiVelocityToAmplitude(note.Velocity))
}

func spawnSineSource(note notes.MidiNote, voiceSendChan chan sources.Voice, sampleRate float64) {
	envelope := envelopes.NewTriangleEnvelope(time.Millisecond*100, time.Millisecond*400, sampleRate)
	s := sources.NewSineVoiceSource(sampleRate, envelope, voiceSendChan)
	s.PublishVoice(notes.MidiPitchToFreq(note.Pitch), notes.MidiVelocityToAmplitude(note.Velocity))
}
