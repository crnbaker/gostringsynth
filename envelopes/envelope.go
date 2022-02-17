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

func NewEnvelope(lengthInSamples int, attackInSamples int) Envelope {
	amps := make([]float32, lengthInSamples)

	for i := 0; i < attackInSamples; i++ {
		amps[i] = float32(i) / float32(attackInSamples-1)
	}

	decayInSamples := lengthInSamples - attackInSamples
	decayIndex := 0
	for i := attackInSamples; i < lengthInSamples; i++ {
		amplitude := (float32(decayInSamples-decayIndex-1) / float32(decayInSamples-1))
		amps[i] = amplitude
		decayIndex++
	}
	return Envelope{amps, 0}
}
