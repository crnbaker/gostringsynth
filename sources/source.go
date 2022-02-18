package sources

import (
	"github.com/crnbaker/gostringsynth/envelopes"
)

type SynthFunction struct {
	Synthesize        func() float32
	AgeInSamples      int
	LifetimeInSamples int
}

type Source interface {
	Play(pitch float64, amplitude float64)
	synthesize(out [][]float32)
}

type SourceImpl struct {
	synthFunctionOutputChannel chan SynthFunction
}

func NewSourceImpl(synthFunctionOutputChannel chan SynthFunction) SourceImpl {
	return SourceImpl{synthFunctionOutputChannel}
}

type EnvelopedSource struct {
	envelope envelopes.Envelope
}

func (g *EnvelopedSource) SetEnvelope(env envelopes.Envelope) {
	g.envelope = env
}

func NewEnvelopedSource(envelope envelopes.Envelope) EnvelopedSource {
	return EnvelopedSource{envelope}
}
