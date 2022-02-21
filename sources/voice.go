package sources

type Voice struct {
	SynthesisFunc     func() float32
	AgeInSamples      int
	LifetimeInSamples int
	markedForDeath    bool
}

func (v *Voice) KillOnNextCycle() {
	v.markedForDeath = true
}

func (v *Voice) ShouldDie() bool {
	return v.AgeInSamples > v.LifetimeInSamples || v.markedForDeath
}
