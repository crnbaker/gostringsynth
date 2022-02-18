package mixer

import (
	"github.com/crnbaker/gostringsynth/errors"
	"github.com/gordonklaus/portaudio"
)

type Mixer struct {
	*portaudio.Stream
	clips         [][]float32
	clipReadHeads []int
}

func (m *Mixer) addStream(stream *portaudio.Stream) {
	m.Stream = stream
}

func (m *Mixer) addClip(clip []float32) {
	m.clips = append(m.clips, clip)
	m.clipReadHeads = append(m.clipReadHeads, 0)
}

func (m *Mixer) deleteClip(clipIndex int) {
	m.clips = append(m.clips[:clipIndex], m.clips[clipIndex+1:]...)
	m.clipReadHeads = append(m.clipReadHeads[:clipIndex], m.clipReadHeads[clipIndex+1:]...)
}

func (m *Mixer) readClipSample(clipIndex int) (sampleValue float32) {
	if m.clipReadHeads[clipIndex] >= len(m.clips[clipIndex]) {
		m.deleteClip(clipIndex)
	} else if len(m.clips[clipIndex]) > 0 {
		sampleValue = m.clips[clipIndex][m.clipReadHeads[clipIndex]]
		m.clipReadHeads[clipIndex]++
	}
	return
}

func (m *Mixer) output(out [][]float32) {
	for i := range out[0] {
		out[0][i] = 0
		out[1][i] = 0
		for j := range m.clips {
			out[0][i] += m.readClipSample(j)
			out[1][i] += m.readClipSample(j)
		}
	}
}

func newMixer(sampleRate float64) *Mixer {

	clips := make([][]float32, 0)
	clipReadHeads := make([]int, 0)

	mixer := &Mixer{nil, clips, clipReadHeads}

	stream, err := portaudio.OpenDefaultStream(0, 2, sampleRate, 0, mixer.output)
	errors.Chk(err)
	mixer.addStream(stream)
	return mixer
}

func MixClips(clipReceiveChannel chan []float32, sampleRate float64) {
	mixer := newMixer(sampleRate)
	mixer.Start()
	for clip := range clipReceiveChannel {
		mixer.addClip(clip)
	}
}
