package envelopes

type Envelope interface {
	//Causes the envelope to open
	Trigger()
	//Returns current length of the envelope in samples
	GetLength() int
	//Returns the amplitude of the envelope at the current position of its internal play head
	GetAmplitude() float64
	//Steps the position of the internal play head forwards by one sample
	Step()
}

type EnvelopeImpl struct {
	amplitudes []float64
	index      int
}

func (e *EnvelopeImpl) GetLength() int {
	return len(e.amplitudes)
}

func (e *EnvelopeImpl) GetAmplitude() float64 {
	var amplitude float64
	if e.index < len(e.amplitudes) {
		amplitude = e.amplitudes[e.index]
	} else {
		amplitude = 0.0
	}
	return amplitude
}

func (e *EnvelopeImpl) Trigger() {
	e.index = 0
}

func (e *EnvelopeImpl) Step() {
	e.index += 1
}

func NewEnvelope(amplitudes []float64) *EnvelopeImpl {
	return &EnvelopeImpl{amplitudes, 0}
}
