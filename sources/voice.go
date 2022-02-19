package sources

type Voice struct {
	Synthesize        func() float32
	AgeInSamples      int
	LifetimeInSamples int
}
