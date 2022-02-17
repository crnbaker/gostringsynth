package envelopes

type triangleEnvelope struct {
	envelope EnvelopeImpl
}

func (e *triangleEnvelope) GetLength() int {
	return e.envelope.GetLength()
}

func (e *triangleEnvelope) GetAmplitude() float64 {
	return e.envelope.GetAmplitude()
}

func (e *triangleEnvelope) Trigger() {
	e.envelope.Trigger()
}

func (e *triangleEnvelope) Step() {
	e.envelope.Step()
}

func NewTriangleEnvelope(attackInSeconds float64, decayInSeconds float64, sampleRate float64) *triangleEnvelope {

	attackInSamples := int(attackInSeconds * sampleRate)
	decayInSamples := int(decayInSeconds * sampleRate)

	lengthInSamples := attackInSamples + decayInSamples
	amps := make([]float64, lengthInSamples)

	for i := 0; i < attackInSamples; i++ {
		amps[i] = float64(i) / float64(attackInSamples-1)
	}

	decayIndex := 0
	for i := attackInSamples; i < lengthInSamples; i++ {
		amplitude := (float64(decayInSamples-decayIndex-1) / float64(decayInSamples-1))
		amps[i] = amplitude
		decayIndex++
	}

	return &triangleEnvelope{*NewEnvelope(amps)}
}
