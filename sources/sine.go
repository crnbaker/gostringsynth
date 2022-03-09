package sources

import (
	"math"
)

type sineSource struct {
	EnvelopedSource
	voiceSendChan chan Voice
	SampleRate    float64
	Phase, Freq   float64
}

func (s *sineSource) calculateVoiceLifetime() int {
	return int(float64(s.envelope.GetLength()) * 1.1)
}

func (s *sineSource) PublishVoice(freqHz float64, amplitude float64) {
	s.voiceSendChan <- Voice{s.Synthesize, 0, s.calculateVoiceLifetime(), false}
	s.Freq = freqHz
	s.envelope.Trigger(amplitude)
}

func (s *sineSource) step() float64 {
	return s.Freq / s.SampleRate
}

func (s *sineSource) Synthesize() (sampleValue float32) {
	sampleValue = float32(math.Sin(2*math.Pi*s.Phase)) * float32(s.envelope.GetAmplitude())
	_, s.Phase = math.Modf(s.Phase + s.step())
	s.envelope.Step()
	return
}

func NewSineVoiceSource(sampleRate float64, envelope SourceEnvelope, voiceSendChan chan Voice) *sineSource {
	s := &sineSource{NewEnvelopedSource(envelope), voiceSendChan, sampleRate, 0, 0}
	return s
}
