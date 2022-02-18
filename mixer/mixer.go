package mixer

import (
	"sync"

	"github.com/crnbaker/gostringsynth/errors"
	"github.com/crnbaker/gostringsynth/sources"
	"github.com/gordonklaus/portaudio"
)

var mutex = &sync.RWMutex{}

type Mixer struct {
	*portaudio.Stream
	synthFunctions map[int]sources.SynthFunction
}

func (m *Mixer) addStream(stream *portaudio.Stream) {
	m.Stream = stream
}

func (m *Mixer) addSynthFunction(synthFunction sources.SynthFunction) {
	keyUnavailable := true
	var newKey int
	for i := 0; keyUnavailable; i++ {
		if _, ok := m.synthFunctions[i]; !ok {
			keyUnavailable = false
			newKey = i
		}
	}
	mutex.Lock()
	m.synthFunctions[newKey] = synthFunction
	mutex.Unlock()
}

func (m *Mixer) deleteSynthFunctionByKey(key int) {
	mutex.Lock()
	delete(m.synthFunctions, key)
	mutex.Unlock()
}

func (m *Mixer) output(out [][]float32) {
	// Initialise buffer with zeros
	for i := range out[0] {
		out[0][i] = 0
		out[1][i] = 0
	}
	// Add samples values synthesized by currently active voices
	for i := range out[0] {
		for key, f := range m.synthFunctions {
			newSample := f.Synthesize()
			out[0][i] += newSample
			out[1][i] += newSample
			f.AgeInSamples++
			m.synthFunctions[key] = f // Use key to assign because f is a copy
		}
	}
	// Destroy voices that are past their lifetime
	for key, f := range m.synthFunctions {
		if f.AgeInSamples > f.LifetimeInSamples {
			m.deleteSynthFunctionByKey(key)
		}
	}
}

func newMixer(sampleRate float64) *Mixer {

	synthFunctions := make(map[int]sources.SynthFunction)

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
