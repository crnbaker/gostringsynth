package sources

import (
	"fmt"

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
	PlayNote(pitch float64, lengthInSeconds float64, attackInSeconds float64)
	setStream(*portaudio.Stream)
	synthesize(out [][]float32)
}

func NewSource(name string, sampleRate float64) (Source, error) {
	if name == "Sine" {
		return NewStereoSine(sampleRate), nil
	} else {
		return nil, &InvalidSourceError{name}
	}
}
