package sources

import (
	"github.com/crnbaker/gostringsynth/envelopes"
)

type Source interface {
	PlayNote(pitch float64, amplitude float64)
	synthesizeClip(clipOutChannel chan []float32)
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
