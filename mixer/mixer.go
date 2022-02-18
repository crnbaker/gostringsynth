package mixer

import (
	"github.com/crnbaker/gostringsynth/errors"
	"github.com/crnbaker/gostringsynth/sources"
	"github.com/gordonklaus/portaudio"
)

type Mixer struct {
	*portaudio.Stream
	synthFunctions []sources.SynthFunction
}

func (m *Mixer) addStream(stream *portaudio.Stream) {
	m.Stream = stream
}

func (m *Mixer) addSynthFunction(synthFunction sources.SynthFunction) {
	m.synthFunctions = append(m.synthFunctions, synthFunction)
}

func (m *Mixer) deleteSynthFunctionByIndex(i int) {
	m.synthFunctions = append(m.synthFunctions[:i], m.synthFunctions[i+1:]...)
}

func (m *Mixer) output(out [][]float32) {
	for i := range out[0] {
		out[0][i] = 0
		out[1][i] = 0
		for _, f := range m.synthFunctions {
			f.AgeInSamples++
		}
	}
	for i, f := range m.synthFunctions {
		f.Synthesize(out)
		if f.AgeInSamples > f.LifetimeInSamples {
			m.deleteSynthFunctionByIndex(i)
		}
	}
}

func newMixer(sampleRate float64) *Mixer {

	synthFunctions := make([]sources.SynthFunction, 0)

	mixer := &Mixer{nil, synthFunctions}

	stream, err := portaudio.OpenDefaultStream(0, 2, sampleRate, 0, mixer.output)
	errors.Chk(err)
	mixer.addStream(stream)
	return mixer
}

func MixController(synthFunctionReceiveChannel chan sources.SynthFunction, sampleRate float64) {
	mixer := newMixer(sampleRate)
	mixer.Start()
	for f := range synthFunctionReceiveChannel {
		mixer.addSynthFunction(f)
	}
}
