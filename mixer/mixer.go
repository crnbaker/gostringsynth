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
	// Initialise buffer with zeros
	for i := range out[0] {
		out[0][i] = 0
		out[1][i] = 0
	}
	// Add samples values synthesized by currently active voices
	for i := range out[0] {
		for j, f := range m.synthFunctions {
			newSample := f.Synthesize()
			out[0][i] += newSample
			out[1][i] += newSample
			m.synthFunctions[j].AgeInSamples++ // Use index because f is a copy
		}
	}
	// Destroy voices that are past their lifetime
	numRemoved := 0
	for i, f := range m.synthFunctions {
		if f.AgeInSamples > f.LifetimeInSamples {
			m.deleteSynthFunctionByIndex(i - numRemoved)
			numRemoved++
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
