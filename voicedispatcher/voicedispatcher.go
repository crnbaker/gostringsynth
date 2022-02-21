package voicedispatcher

import (
	"sync"
	"time"

	"github.com/crnbaker/gostringsynth/envelopes"
	"github.com/crnbaker/gostringsynth/notes"
	"github.com/crnbaker/gostringsynth/sources"
)

func VoiceDispatcher(waitGroup *sync.WaitGroup, noteInChan chan notes.MidiNote, voiceSendChan chan sources.Voice, sampleRate float64, osc string) {
	defer waitGroup.Done()
	for note := range noteInChan {
		if osc == "sine" {
			go spawnSineVoiceSource(note, voiceSendChan, sampleRate)
		} else if osc == "string" {
			go spawnStringVoiceSource(note, voiceSendChan, sampleRate)
		}
	}
	close(voiceSendChan)
}

func spawnSineVoiceSource(note notes.MidiNote, voiceSendChan chan sources.Voice, sampleRate float64) {
	envelope := envelopes.NewTriangleEnvelope(time.Millisecond*100, time.Millisecond*400, sampleRate)
	s := sources.NewSineVoiceSource(sampleRate, envelope, voiceSendChan)
	s.PublishVoice(notes.MidiPitchToFreq(note.Pitch), notes.MidiVelocityToAmplitude(note.Velocity))
}

func spawnStringVoiceSource(note notes.MidiNote, voiceSendChan chan sources.Voice, sampleRate float64) {
	const waveSpeedMpS = 200
	const pickupPos = 0.5
	const decayTimeS = 3.0
	lengthM := sources.FreqToStringLength(notes.MidiPitchToFreq(note.Pitch), waveSpeedMpS)
	s := sources.NewStringSource(sampleRate, voiceSendChan, lengthM, waveSpeedMpS, pickupPos, decayTimeS)
	s.PublishVoice(notes.MidiPitchToFreq(note.Pitch), notes.MidiVelocityToAmplitude(note.Velocity))
}
