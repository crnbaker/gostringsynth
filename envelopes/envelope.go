package envelopes

type Envelope struct {
	amplitudes []float32
	index      int
}

func (e *Envelope) GetLength() int {
	return len(e.amplitudes)
}

func (e *Envelope) GetAmplitude() float32 {
	var amplitude float32
	if e.index < len(e.amplitudes) {
		amplitude = e.amplitudes[e.index]
	} else {
		amplitude = 0.0
	}
	return amplitude
}

func (e *Envelope) Trigger() {
	e.index = 0
}

func (e *Envelope) Step() {
	e.index += 1
}

func NewEnvelope(lengthInSamples int) Envelope {
	amps := make([]float32, lengthInSamples)
	for i := range amps {
		amplitude := (float32(lengthInSamples-i-1) / float32(lengthInSamples-1))
		amps[i] = amplitude
	}
	return Envelope{amps, 0}
}
