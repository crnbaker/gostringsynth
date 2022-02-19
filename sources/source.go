package sources

import (
	"github.com/crnbaker/gostringsynth/envelopes"
)

type Source interface {
	DispatchAndPlayVoice(pitch float64, amplitude float64)
	Synthesize(out [][]float32)
}

type VoiceSource struct {
	voiceSendChan chan Voice
}

func NewVoiceSource(voiceSendChan chan Voice) VoiceSource {
	return VoiceSource{voiceSendChan}
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
