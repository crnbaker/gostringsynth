package sources

import (
	"github.com/crnbaker/gostringsynth/envelopes"
)

type Source interface {
	DispatchAndPlayVoice(freqHz float64, amplitude float64)
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

type FDTDSource struct {
	fdtdGrid           [][]float64 // N time steps x M spatial points
	numSpatialSections int
	numTimeSteps       int
}

func NewFTDTSource(numTimeSteps int, numSpatialSections int) FDTDSource {

	fdtdGrid := make([][]float64, 3)
	for i := range fdtdGrid {
		fdtdGrid[i] = make([]float64, numSpatialSections+1)
	}
	return FDTDSource{fdtdGrid, numSpatialSections, numTimeSteps}
}
