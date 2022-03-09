package sources

import (
	"math"
)

type sineSource struct {
	envelopedSource
	SampleRate  float64
	Phase, Freq float64
}

func (s *sineSource) calculateVoiceLifetime() int {
	return int(float64(s.envelope.GetLength()) * 1.1)
}

func (s *sineSource) Voice(freqHz float64, amplitude float64) Voice {
	s.Freq = freqHz
	s.envelope.Trigger(amplitude)
	return Voice{s.Synthesize, 0, s.calculateVoiceLifetime(), false}
}

func (s *sineSource) step() float64 {
	return s.Freq / s.SampleRate
}

func (s *sineSource) Synthesize() float32 {
	sampleValue := float32(math.Sin(2*math.Pi*s.Phase)) * float32(s.envelope.GetAmplitude())
	_, s.Phase = math.Modf(s.Phase + s.step())
	s.envelope.Step()
	return sampleValue
}

func NewSineVoiceSource(sampleRate float64, envelope SourceEnvelope, voiceSendChan chan Voice) *sineSource {
	s := &sineSource{newEnvelopedSource(envelope), sampleRate, 0, 0}
	return s
}
