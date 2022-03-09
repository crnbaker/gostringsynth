/* The sources package provides a Source interface for defining and exporting synthesis functions.

A synthesis function is exported from a source and packaged into a Voice struct, which contains
a reference to the synthesis function and attributes relating to the lifetim of the Voice.
*/
package sources

// fdtdSource provides Sources with attributes requried to perform 1D finite difference synthesis
type fdtdSource struct {
	fdtdGrid           [][]float64 // N time steps x M spatial points
	numSpatialSections int
	numTimeSteps       int
}

// newFtdtSource constructs a new FDTDSource for simualtions of a given temporal and spatial size
func newFtdtSource(numTimeSteps int, numSpatialSections int) fdtdSource {

	fdtdGrid := make([][]float64, 3)
	for i := range fdtdGrid {
		fdtdGrid[i] = make([]float64, numSpatialSections+1)
	}
	return fdtdSource{fdtdGrid, numSpatialSections, numTimeSteps}
}

// envelopedSource provides an amplitude envelope attribute to Sources that need one, i.e. traditional
// oscillator-based synth sources rather that finite difference sources. EnvelopedSources are primarily
// intended for testing and developement.
type envelopedSource struct {
	envelope SourceEnvelope
}

// setEnvelope is used to change the amplitude envelope of an already-constructed EnvelopedSource
func (e *envelopedSource) setEnvelope(env SourceEnvelope) {
	e.envelope = env
}

// newEnvelopedSource constructs a new EnvelopedSource using a given envelopes.Envelope.
func newEnvelopedSource(envelope SourceEnvelope) envelopedSource {
	return envelopedSource{envelope}
}
