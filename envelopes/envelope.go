package envelopes

type Envelope struct {
	scale      float64
	amplitudes []float64
	index      int
}

func (e *Envelope) GetLength() int {
	return len(e.amplitudes)
}

func (e *Envelope) GetAmplitude() float64 {
	var amplitude float64
	if e.index < len(e.amplitudes) {
		amplitude = e.amplitudes[e.index] * e.scale
	} else {
		amplitude = 0.0
	}
	return amplitude
}

func (e *Envelope) Trigger(amplitude float64) {
	e.scale = amplitude
	e.index = 0
}

func (e *Envelope) Step() {
	e.index += 1
}

func NewEnvelope(amplitudes []float64) *Envelope {
	return &Envelope{1.0, amplitudes, 0}
}
