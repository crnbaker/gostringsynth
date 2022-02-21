package sources

// Voice is a struct used to package and export a synthesis function from a Source while controlling its lifetime.
type Voice struct {
	SynthesisFunc     func() float32
	AgeInSamples      int
	LifetimeInSamples int
	markedForDeath    bool
}

// KillOnNextCycle marks a Voice for death, so that the next call to ShouldDie will return true
func (v *Voice) KillOnNextCycle() {
	v.markedForDeath = true
}

// ShouldDie returns true if the Vocie is older than its lifetime or if KillOnNextCycle has been previously called
func (v *Voice) ShouldDie() bool {
	return v.AgeInSamples > v.LifetimeInSamples || v.markedForDeath
}
