package notedispatcher

import (
	"time"

	"github.com/crnbaker/gostringsynth/envelopes"
	"github.com/crnbaker/gostringsynth/keypress"
	"github.com/crnbaker/gostringsynth/sources"
)

func NoteDispatcher(noteInChan chan keypress.MidiNote, voiceOutputChan chan sources.Voice, quitChan chan bool, sampleRate float64) {
	for note := range noteInChan {
		go spawnVoice(note, voiceOutputChan, sampleRate)
	}
	close(quitChan)
	close(voiceOutputChan)
}

func spawnVoice(note keypress.MidiNote, voiceOutputChan chan sources.Voice, sampleRate float64) {
	envelope := envelopes.NewTriangleEnvelope(time.Millisecond*100, time.Millisecond*400, sampleRate)
	s := sources.NewSineSource(sampleRate, envelope, voiceOutputChan)
	s.Play(keypress.MidiPitchToFreq(note.Pitch), keypress.MidiVelocityToAmplitude(note.Velocity))
}
