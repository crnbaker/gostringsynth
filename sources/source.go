package sources

import (
	"github.com/crnbaker/gostringsynth/envelopes"
)

type Source interface {
	Play(pitch float64, amplitude float64)
	synthesize(out [][]float32)
}

type SourceImpl struct {
	voiceSendChan chan Voice
}

func NewSourceImpl(voiceSendChan chan Voice) SourceImpl {
	return SourceImpl{voiceSendChan}
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
