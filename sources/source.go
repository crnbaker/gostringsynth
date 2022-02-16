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

type source interface {
	Start() error
	Stop() error
	Close() error
	PlayNote(pitch float64, lengthInSeconds float64)
	setStream(*portaudio.Stream)
	synthesize(out [][]float32)
}
