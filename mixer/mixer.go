package mixer

import (
	"sync"

	"github.com/crnbaker/gostringsynth/errors"
	"github.com/crnbaker/gostringsynth/sources"
	"github.com/gordonklaus/portaudio"
)

type Mixer struct {
	*portaudio.Stream
	activeVoices []sources.Voice
	stagedVoices []sources.Voice
}

func (m *Mixer) addStream(stream *portaudio.Stream) {
	m.Stream = stream
}

func (m *Mixer) stageVoice(voice sources.Voice) {
	m.stagedVoices = append(m.stagedVoices, voice)
}

func (m *Mixer) addVoice(voice sources.Voice) {
	m.activeVoices = append(m.activeVoices, voice)
}

func (m *Mixer) activateStagedVoices() {
	for _, v := range m.stagedVoices {
		m.addVoice(v)
		m.stagedVoices = m.stagedVoices[1:]
	}
}

func (m *Mixer) killVoice(i int) {
	m.activeVoices = append(m.activeVoices[:i], m.activeVoices[i+1:]...)
}

func (m *Mixer) output(out [][]float32) {
	// Kill voices that are past their lifetime or have been stolen
	numKilled := 0
	for i, f := range m.activeVoices {
		if f.ShouldDie() {
			m.killVoice(i - numKilled)
			numKilled++
		}
	}
	// Activate new voices that have been staged for activation
	m.activateStagedVoices()

	// Initialise buffer with zeros
	for i := range out[0] {
		out[0][i] = 0
		out[1][i] = 0
	}
	// Add samples values synthesized by currently active voices
	for i := range out[0] {
		for j, f := range m.activeVoices {
			newSample := f.SynthesisFunc()
			out[0][i] += newSample
			out[1][i] += newSample
			m.activeVoices[j].AgeInSamples++ // Use index because f is a copy
		}
	}
}

func newMixer(sampleRate float64) *Mixer {

	activeVoices := make([]sources.Voice, 0)
	stagedVoices := make([]sources.Voice, 0)
	mixer := &Mixer{nil, activeVoices, stagedVoices}

	stream, err := portaudio.OpenDefaultStream(0, 2, sampleRate, 0, mixer.output)
	errors.Chk(err)
	mixer.addStream(stream)
	return mixer
}

func MixController(waitGroup *sync.WaitGroup, voiceReceiveChan chan sources.Voice, sampleRate float64, voiceLimit int) {
	defer waitGroup.Done()
	portaudio.Initialize()
	mixer := newMixer(sampleRate)
	mixer.Start()
	for f := range voiceReceiveChan {
		if len(mixer.activeVoices) == voiceLimit {
			mixer.activeVoices[0].KillOnNextCycle()
		}
		mixer.stageVoice(f)
	}
	mixer.Stop()
	portaudio.Terminate()
}
