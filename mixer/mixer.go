package mixer

import (
	"github.com/crnbaker/gostringsynth/errors"
	"github.com/crnbaker/gostringsynth/sources"
	"github.com/gordonklaus/portaudio"
)

type Mixer struct {
	*portaudio.Stream
	voices []sources.Voice
}

func (m *Mixer) addStream(stream *portaudio.Stream) {
	m.Stream = stream
}

func (m *Mixer) addVoice(synthFunction sources.Voice) {
	m.voices = append(m.voices, synthFunction)
}

func (m *Mixer) killVoice(i int) {
	m.voices = append(m.voices[:i], m.voices[i+1:]...)
}

func (m *Mixer) output(out [][]float32) {
	// Initialise buffer with zeros
	for i := range out[0] {
		out[0][i] = 0
		out[1][i] = 0
	}
	// Add samples values synthesized by currently active voices
	for i := range out[0] {
		for j, f := range m.voices {
			newSample := f.SynthesisFunc()
			out[0][i] += newSample
			out[1][i] += newSample
			m.voices[j].AgeInSamples++ // Use index because f is a copy
		}
	}
	// Kill voices that are past their lifetime
	numKilled := 0
	for i, f := range m.voices {
		if f.ShouldDie() {
			m.killVoice(i - numKilled)
			numKilled++
		}
	}
}

func newMixer(sampleRate float64) *Mixer {

	synthFunctions := make([]sources.Voice, 0)

	mixer := &Mixer{nil, synthFunctions}

	stream, err := portaudio.OpenDefaultStream(0, 2, sampleRate, 0, mixer.output)
	errors.Chk(err)
	mixer.addStream(stream)
	return mixer
}

func MixController(voiceReceiveChan chan sources.Voice, sampleRate float64) {
	mixer := newMixer(sampleRate)
	mixer.Start()
	for f := range voiceReceiveChan {
		mixer.addVoice(f)
	}
	mixer.Stop()
	portaudio.Terminate()
}
