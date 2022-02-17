package envelopes

type Envelope interface {
	//Causes the envelope to open
	Trigger(amplitude float64)
	//Returns current length of the envelope in samples
	GetLength() int
	//Returns the amplitude of the envelope at the current position of its internal play head
	GetAmplitude() float64
	//Steps the position of the internal play head forwards by one sample
	Step()
}

type EnvelopeImpl struct {
	scale      float64
	amplitudes []float64
	index      int
}

func (e *EnvelopeImpl) GetLength() int {
	return len(e.amplitudes)
}

func (e *EnvelopeImpl) GetAmplitude() float64 {
	var amplitude float64
	if e.index < len(e.amplitudes) {
		amplitude = e.amplitudes[e.index] * e.scale
	} else {
		amplitude = 0.0
	}
	return amplitude
}

func (e *EnvelopeImpl) Trigger(amplitude float64) {
	e.scale = amplitude
	e.index = 0
}

func (e *EnvelopeImpl) Step() {
	e.index += 1
}

func NewEnvelope(amplitudes []float64) *EnvelopeImpl {
	return &EnvelopeImpl{1.0, amplitudes, 0}
}
