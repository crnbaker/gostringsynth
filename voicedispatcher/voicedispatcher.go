package voicedispatcher

import (
	"time"

	"github.com/crnbaker/gostringsynth/envelopes"
	"github.com/crnbaker/gostringsynth/keypress"
	"github.com/crnbaker/gostringsynth/sources"
)

func VoiceDispatcher(noteInChan chan keypress.MidiNote, voiceSendChan chan sources.Voice, quitChan chan bool, sampleRate float64, osc string) {
	for note := range noteInChan {
		if osc == "sine" {
			go spawnSineVoiceSource(note, voiceSendChan, sampleRate)
		} else if osc == "string" {
			go spawnStringVoiceSource(note, voiceSendChan, sampleRate)
		}
	}
	close(quitChan)
	close(voiceSendChan)
}

func spawnSineVoiceSource(note keypress.MidiNote, voiceSendChan chan sources.Voice, sampleRate float64) {
	envelope := envelopes.NewTriangleEnvelope(time.Millisecond*100, time.Millisecond*400, sampleRate)
	s := sources.NewSineVoiceSource(sampleRate, envelope, voiceSendChan)
	s.DispatchAndPlayVoice(keypress.MidiPitchToFreq(note.Pitch), keypress.MidiVelocityToAmplitude(note.Velocity))
}

func spawnStringVoiceSource(note keypress.MidiNote, voiceSendChan chan sources.Voice, sampleRate float64) {
	const waveSpeedMpS = 200
	const pickupPos = 0.2
	const decayTimeS = 5.0
	lengthM := sources.FreqToStringLength(keypress.MidiPitchToFreq(note.Pitch), waveSpeedMpS)
	s := sources.NewStringVoiceSource(sampleRate, voiceSendChan, lengthM, waveSpeedMpS, pickupPos, decayTimeS)
	s.DispatchAndPlayVoice(keypress.MidiPitchToFreq(note.Pitch), keypress.MidiVelocityToAmplitude(note.Velocity))
}
