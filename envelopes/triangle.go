package envelopes

import "time"

type triangleEnvelope struct {
	Envelope
}

func NewTriangleEnvelope(attackTime time.Duration, decayTime time.Duration, sampleRate float64) *triangleEnvelope {

	attackInSamples := int(attackTime.Seconds() * sampleRate)
	decayInSamples := int(decayTime.Seconds() * sampleRate)

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
