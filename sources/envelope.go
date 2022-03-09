package sources

type SourceEnvelope interface {
	//Causes the envelope to open
	Trigger(amplitude float64)
	//Returns current length of the envelope in samples
	GetLength() int
	//Returns the amplitude of the envelope at the current position of its internal play head
	GetAmplitude() float64
	//Steps the position of the internal play head forwards by one sample
	Step()
}
