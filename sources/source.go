package sources

import (
	"fmt"

	"github.com/crnbaker/gostringsynth/envelopes"

	"github.com/gordonklaus/portaudio"
)

type InvalidSourceError struct {
	invalidSource string
}

func (e *InvalidSourceError) Error() string {
	return fmt.Sprintf("%s is not a valid source.", e.invalidSource)
}

type Source interface {
	Start() error
	Stop() error
	Close() error
	PlayNote(pitch float64)
	SetEnvelope(envelope envelopes.Envelope)
	setStream(*portaudio.Stream)
	synthesize(out [][]float32)
}

func NewSource(name string, sampleRate float64) (Source, error) {
	env := envelopes.NewEnvelope(make([]float64, 1))
	if name == "Sine" {
		return NewStereoSine(sampleRate, env), nil
	} else {
		return nil, &InvalidSourceError{name}
	}
}
