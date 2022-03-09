/* The sources package provides a Source interface for defining and exporting synthesis functions.

A synthesis function is exported from a source and packaged into a Voice struct, which contains
a reference to the synthesis function and attributes relating to the lifetim of the Voice.
*/
package sources

// Source defines an interface for structs that publish an audio synthesis function.
type Source interface {
	PublishVoice(freqHz float64, amplitude float64)
	Synthesize(out [][]float32)
}

// FDTDSource provides Sources with attributes requried to perform 1D finite difference synthesis
type FDTDSource struct {
	fdtdGrid           [][]float64 // N time steps x M spatial points
	numSpatialSections int
	numTimeSteps       int
}

// NewFTDTSource constructs a new FDTDSource for simualtions of a given temporal and spatial size
func NewFTDTSource(numTimeSteps int, numSpatialSections int) FDTDSource {

	fdtdGrid := make([][]float64, 3)
	for i := range fdtdGrid {
		fdtdGrid[i] = make([]float64, numSpatialSections+1)
	}
	return FDTDSource{fdtdGrid, numSpatialSections, numTimeSteps}
}

// EnvelopedSource provides an amplitude envelope attribute to Sources that need one, i.e. traditional
// oscillator-based synth sources rather that finite difference sources. EnvelopedSources are primarily
// intended for testing and developement.
type EnvelopedSource struct {
	envelope SourceEnvelope
}

// SetEnvelope is used to change the amplitude envelope of an already-constructed EnvelopedSource
func (e *EnvelopedSource) SetEnvelope(env SourceEnvelope) {
	e.envelope = env
}

// NewEnvelopedSource constructs a new EnvelopedSource using a given envelopes.Envelope.
func NewEnvelopedSource(envelope SourceEnvelope) EnvelopedSource {
	return EnvelopedSource{envelope}
}
