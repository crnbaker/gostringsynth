package sources

import (
	"math"

	"github.com/crnbaker/gostringsynth/envelopes"
	"github.com/crnbaker/gostringsynth/errors"

	"github.com/gordonklaus/portaudio"
)

type stereoSine struct {
	*portaudio.Stream
	SampleRate    float64
	PhaseL, FreqL float64
	PhaseR, FreqR float64
	envelope      envelopes.Envelope
}

func (g *stereoSine) SetEnvelope(env envelopes.Envelope) {
	g.envelope = env
}

func (g *stereoSine) PlayNote(pitch float64) {
	g.FreqL = pitch
	g.FreqR = pitch
	g.envelope.Trigger()
}

func (g *stereoSine) setStream(stream *portaudio.Stream) {
	g.Stream = stream
}

func (g *stereoSine) stepL() float64 {
	return g.FreqL / g.SampleRate
}

func (g *stereoSine) stepR() float64 {
	return g.FreqR / g.SampleRate
}

func (g *stereoSine) synthesize(out [][]float32) {
	for i := range out[0] {
		out[0][i] = float32(math.Sin(2*math.Pi*g.PhaseL)) * float32(g.envelope.GetAmplitude())
		_, g.PhaseL = math.Modf(g.PhaseL + g.stepL())
		out[1][i] = float32(math.Sin(2*math.Pi*g.PhaseR)) * float32(g.envelope.GetAmplitude())
		_, g.PhaseR = math.Modf(g.PhaseR + g.stepR())
		g.envelope.Step()
	}
}

func NewStereoSine(sampleRate float64, envelope envelopes.Envelope) *stereoSine {
	s := &stereoSine{nil, sampleRate, 0, 0, 0, 0, envelope}
	var err error
	var stream *portaudio.Stream
	stream, err = portaudio.OpenDefaultStream(0, 2, sampleRate, 0, s.synthesize)
	s.setStream(stream)
	errors.Chk(err)
	return s
}
