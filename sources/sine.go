package sources

import (
	"math"

	"github.com/crnbaker/gostringsynth/envelopes"
)

type sineSource struct {
	SourceImpl
	EnvelopedSource
	SampleRate  float64
	Phase, Freq float64
}

func (g *sineSource) Play(pitch float64, amplitude float64) {
	lifetime := int(float64(g.envelope.GetLength()) * 1.1)
	g.synthFunctionOutputChannel <- SynthFunction{g.Synthesize, 0, lifetime}
	g.Freq = pitch
	g.envelope.Trigger(amplitude)
}

func (g *sineSource) step() float64 {
	return g.Freq / g.SampleRate
}

func (g *sineSource) Synthesize() (sampleValue float32) {
	sampleValue = float32(math.Sin(2*math.Pi*g.Phase)) * float32(g.envelope.GetAmplitude())
	_, g.Phase = math.Modf(g.Phase + g.step())
	g.envelope.Step()
	return
}

func NewSineSource(sampleRate float64, envelope envelopes.Envelope, outputChannel chan SynthFunction) *sineSource {
	s := &sineSource{NewSourceImpl(outputChannel), NewEnvelopedSource(envelope), sampleRate, 0, 0}
	return s
}
