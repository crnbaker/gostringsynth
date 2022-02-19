package voicedispatcher

import (
	"time"

	"github.com/crnbaker/gostringsynth/envelopes"
	"github.com/crnbaker/gostringsynth/keypress"
	"github.com/crnbaker/gostringsynth/sources"
)

func VoiceDispatcher(noteInChan chan keypress.MidiNote, voiceOutputChan chan sources.Voice, quitChan chan bool, sampleRate float64) {
	for note := range noteInChan {
		go spawnVoiceSource(note, voiceOutputChan, sampleRate)
	}
	close(quitChan)
	close(voiceOutputChan)
}

func spawnVoiceSource(note keypress.MidiNote, voiceOutputChan chan sources.Voice, sampleRate float64) {
	envelope := envelopes.NewTriangleEnvelope(time.Millisecond*100, time.Millisecond*400, sampleRate)
	s := sources.NewSineVoiceSource(sampleRate, envelope, voiceOutputChan)
	s.DispatchAndPlayVoice(keypress.MidiPitchToFreq(note.Pitch), keypress.MidiVelocityToAmplitude(note.Velocity))
}
