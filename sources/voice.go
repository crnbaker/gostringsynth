package sources

// Voice is a struct used to package and export a synthesis function from a Source while controlling its lifetime.
type Voice struct {
	synthesisFunc     func() float32
	ageInSamples      int
	lifetimeInSamples int
	markedForDeath    bool
}

// KillOnNextCycle marks a Voice for death, so that the next call to ShouldDie will return true
func (v *Voice) KillOnNextCycle() {
	v.markedForDeath = true
}

// ShouldDie returns true if the Voice is older than its lifetime or if KillOnNextCycle has been previously called
func (v *Voice) ShouldDie() bool {
	return v.ageInSamples > v.lifetimeInSamples || v.markedForDeath
}

// IncrementAgeInSamples adds one to the age of the voice in samples
func (v *Voice) IncrementAgeInSamples() {
	v.ageInSamples++
}

// SynthesizeSample synthesises an audio sample using the voice's synthesis function
func (v *Voice) SynthesizeSample() float32 {
	return v.synthesisFunc()
}
