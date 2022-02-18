package voicecontrol

import (
	"time"

	"github.com/crnbaker/gostringsynth/envelopes"
	"github.com/crnbaker/gostringsynth/errors"
	"github.com/crnbaker/gostringsynth/sources"
	"github.com/gordonklaus/portaudio"
)

func NoteDispatcher(pitchInChan chan float64, closeVoiceChan chan sources.Source, quitChan chan bool, sampleRate float64) {
	go voiceCloser(closeVoiceChan, quitChan)
	for pitch := range pitchInChan {
		go startVoice(pitch, closeVoiceChan, sampleRate)
	}
	close(closeVoiceChan)
}

func startVoice(pitch float64, closeVoiceChan chan sources.Source, sampleRate float64) {
	envelope := envelopes.NewTriangleEnvelope(time.Millisecond*100, time.Millisecond*400, sampleRate)
	s := sources.NewStereoSine(sampleRate, envelope)
	errors.Chk(s.Start())
	s.PlayNote(pitch, 0.1)
	time.Sleep(time.Millisecond * 500)
	closeVoiceChan <- s
}

func voiceCloser(closeVoiceChan chan sources.Source, quitChan chan bool) {
	for source := range closeVoiceChan {
		time.Sleep(time.Millisecond * 250)
		source.Stop()
		// source.Close()
	}
	time.Sleep(time.Millisecond * 250)
	portaudio.Terminate()
	close(quitChan)
}
