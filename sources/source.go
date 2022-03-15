/* The sources package contains structs and functions for synthesising audio samples.

Synthesis functions are exported from the sources package inside Voice structs , which contain
a reference to the synthesis function and attributes relating to its lifetime. The Voice
struct is designed to conform to the SynthVoice interface required by the audioengine.
*/
package sources

// FDTDSource provides Sources with attributes requried to perform 1D finite difference synthesis
type FDTDSource struct {
	sampleRate         float64
	fdtdGrid           [][]float64 // N time steps x M spatial points
	numSpatialSections int
	numTimeSteps       int
}

// newFtdtSource constructs a new FDTDSource for simualtions of a given temporal and spatial size
func newFtdtSource(sampleRate float64, numTimeSteps int, numSpatialSections int) FDTDSource {

	fdtdGrid := make([][]float64, 3)
	for i := range fdtdGrid {
		fdtdGrid[i] = make([]float64, numSpatialSections+1)
	}
	return FDTDSource{sampleRate, fdtdGrid, numSpatialSections, numTimeSteps}
}

func (s *FDTDSource) GetSampleRate() float64 {
	return s.sampleRate
}

func (s *FDTDSource) GetFDTDGrid() [][]float64 {
	return s.fdtdGrid
}

// envelopedSource provides an amplitude envelope attribute to Sources that need one, i.e. traditional
// oscillator-based synth sources rather that finite difference sources. EnvelopedSources are primarily
// intended for testing and development.
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
