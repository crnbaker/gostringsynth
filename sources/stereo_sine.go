package sources

import (
	"math"

	"github.com/crnbaker/gostringsynth/envelopes"
)

type sineSource struct {
	EnvelopedSource
	SampleRate  float64
	Phase, Freq float64
}

func (g *sineSource) PlayNote(pitch float64, amplitude float64) {
	g.Freq = pitch
	g.envelope.Trigger(amplitude)
}

func (g *sineSource) step() float64 {
	return g.Freq / g.SampleRate
}

func (g *sineSource) newClipSynthesizer(clipOutChannel chan []float32) {
	clip := make([]float32, g.envelope.GetLength())
	for i := 0; i < g.envelope.GetLength(); i++ {
		clip[i] = float32(math.Sin(2*math.Pi*g.Phase)) * float32(g.envelope.GetAmplitude())
		_, g.Phase = math.Modf(g.Phase + g.step())
		g.envelope.Step()
	}
	clipOutChannel <- clip
}

func NewSineSource(sampleRate float64, envelope envelopes.Envelope) *sineSource {
	s := &sineSource{NewEnvelopedSource(envelope), sampleRate, 0, 0}
	return s
}
