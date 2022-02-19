package notedispatcher

import (
	"time"

	"github.com/crnbaker/gostringsynth/envelopes"
	"github.com/crnbaker/gostringsynth/sources"
)

func NoteDispatcher(pitchInChan chan float64, voiceOutputChan chan sources.Voice, quitChan chan bool, sampleRate float64) {
	for pitch := range pitchInChan {
		go spawnVoice(pitch, voiceOutputChan, sampleRate)
	}
	close(quitChan)
	close(voiceOutputChan)
}

func spawnVoice(pitch float64, voiceOutputChan chan sources.Voice, sampleRate float64) {
	envelope := envelopes.NewTriangleEnvelope(time.Millisecond*100, time.Millisecond*400, sampleRate)
	s := sources.NewSineSource(sampleRate, envelope, voiceOutputChan)
	s.Play(pitch, 0.1)
}
