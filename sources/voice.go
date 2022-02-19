package sources

type Voice struct {
	Synthesize        func() float32
	AgeInSamples      int
	LifetimeInSamples int
}

func (v *Voice) ShouldDie() bool {
	return v.AgeInSamples > v.LifetimeInSamples
}
