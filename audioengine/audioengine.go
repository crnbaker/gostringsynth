/*
The audioengine package creates a portaudio stream and synthesises audio samples into the stream using Voices
it receives.
*/
package audioengine

import (
	"sync"

	"github.com/crnbaker/gostringsynth/errors"
	"github.com/crnbaker/gostringsynth/sources"
	"github.com/gordonklaus/portaudio"
)

// VoiceController provides a portaudio output stream and attributes for keeping track of currently enabled voices.
type VoiceController struct {
	*portaudio.Stream
	activeVoices []sources.Voice
	stagedVoices []sources.Voice
}

// setStream sets a portaudio stream to a VoiceController
func (m *VoiceController) setStream(stream *portaudio.Stream) {
	m.Stream = stream
}

// stageVoice adds a Voice to the list of voices that will be enabled at the beginning of the next iteration of
// the audio output loop
func (m *VoiceController) stageVoice(voice sources.Voice) {
	m.stagedVoices = append(m.stagedVoices, voice)
}

// addVoice adds a voice to the list of currently active voices
func (m *VoiceController) addVoice(voice sources.Voice) {
	m.activeVoices = append(m.activeVoices, voice)
}

// activateStagedVoices activates all voices in the staged voices list
func (m *VoiceController) activateStagedVoices() {
	for _, v := range m.stagedVoices {
		m.addVoice(v)
		m.stagedVoices = m.stagedVoices[1:]
	}
}

// killVoice deletes a voice from the list of active voices using its index
func (m *VoiceController) killVoice(i int) {
	m.activeVoices = append(m.activeVoices[:i], m.activeVoices[i+1:]...)
}

// output is provided to portaudio as the audio generation callback function. It generates audio samples by
// summing the samples provided by the synthesis functions of the currently active Voices.
func (m *VoiceController) output(out [][]float32) {
	// Kill voices that are past their lifetime or have been "stolen"
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

// newVoiceController constructs a VoiceController with a portaudio stream configured to use the output function
// as the audio generationo callback.
func newVoiceController(sampleRate float64) *VoiceController {

	activeVoices := make([]sources.Voice, 0)
	stagedVoices := make([]sources.Voice, 0)
	mixer := &VoiceController{nil, activeVoices, stagedVoices}

	stream, err := portaudio.OpenDefaultStream(0, 2, sampleRate, 0, mixer.output)
	errors.Chk(err)
	mixer.setStream(stream)
	return mixer
}

// ControlVoices receives voices from the voiceReceiveChan and stages them for activation by the voice controller.
// It also implements voice stealing by marking the oldest voice for death if a maximum number of activated voices
// is exceeded.
func ControlVoices(waitGroup *sync.WaitGroup, voiceReceiveChan chan sources.Voice, sampleRate float64, voiceLimit int) {
	defer waitGroup.Done()
	portaudio.Initialize()
	mixer := newVoiceController(sampleRate)
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
