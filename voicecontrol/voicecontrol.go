package voicecontrol

import (
	"time"

	"github.com/crnbaker/gostringsynth/envelopes"
	"github.com/crnbaker/gostringsynth/sources"
)

func NoteDispatcher(pitchInChan chan float64, synthFunctionOutputChannel chan sources.SynthFunction, quitChan chan bool, sampleRate float64) {
	for pitch := range pitchInChan {
		go play(pitch, synthFunctionOutputChannel, sampleRate)
	}
	close(quitChan)
	close(synthFunctionOutputChannel)
}

func play(pitch float64, synthFunctionOutputChannel chan sources.SynthFunction, sampleRate float64) {
	envelope := envelopes.NewTriangleEnvelope(time.Millisecond*100, time.Millisecond*400, sampleRate)
	s := sources.NewSineSource(sampleRate, envelope, synthFunctionOutputChannel)
	s.Play(pitch, 0.1)
}
