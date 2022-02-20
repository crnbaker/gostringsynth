package sources

import (
	"math"

	"github.com/crnbaker/gostringsynth/envelopes"
)

type sineVoiceSource struct {
	VoiceSource
	EnvelopedSource
	SampleRate  float64
	Phase, Freq float64
}

func (s *sineVoiceSource) calculateVoiceLifetime() int {
	return int(float64(s.envelope.GetLength()) * 1.1)
}

func (s *sineVoiceSource) DispatchAndPlayVoice(freqHz float64, amplitude float64) {
	s.voiceSendChan <- Voice{s.Synthesize, 0, s.calculateVoiceLifetime()}
	s.Freq = freqHz
	s.envelope.Trigger(amplitude)
}

func (s *sineVoiceSource) step() float64 {
	return s.Freq / s.SampleRate
}

func (s *sineVoiceSource) Synthesize() (sampleValue float32) {
	sampleValue = float32(math.Sin(2*math.Pi*s.Phase)) * float32(s.envelope.GetAmplitude())
	_, s.Phase = math.Modf(s.Phase + s.step())
	s.envelope.Step()
	return
}

func NewSineVoiceSource(sampleRate float64, envelope envelopes.Envelope, voiceSendChan chan Voice) *sineVoiceSource {
	s := &sineVoiceSource{NewVoiceSource(voiceSendChan), NewEnvelopedSource(envelope), sampleRate, 0, 0}
	return s
}
