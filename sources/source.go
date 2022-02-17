package sources

import (
	"github.com/crnbaker/gostringsynth/envelopes"

	"github.com/gordonklaus/portaudio"
)

type Source interface {
	Start() error
	Stop() error
	Close() error
	PlayNote(pitch float64)
	setStream(*portaudio.Stream)
	synthesize(out [][]float32)
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
