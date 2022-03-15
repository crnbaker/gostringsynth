package sources

import (
	"math"
	"time"

	"github.com/crnbaker/gostringsynth/envelopes"
	"github.com/crnbaker/gostringsynth/excitors"
	"github.com/crnbaker/gostringsynth/numeric"
)

type Excitor interface {
	Excite(source excitors.Excitable, amplitude float64)
}

// stringSource provides attributes that define a finite-difference simulation of a vibrating string
type stringSource struct {
	FDTDSource
	envelopedSource
	stringLengthM     float64
	physics           stringSettings
	transientExcitors []Excitor
	continousExcitors []Excitor
}

func (s *stringSource) transientExcitation(amplitudes []float64) {
	for i, e := range s.transientExcitors {
		e.Excite(&s.FDTDSource, amplitudes[i])
	}
}

func (s *stringSource) continuousExcitation(amplitudes []float64) {
	for i, e := range s.continousExcitors {
		e.Excite(&s.FDTDSource, amplitudes[i])
	}
}

// calculateLossFactor returns a loss factor used to attenuated the string vibration during synthesis
func (s *stringSource) calculateLossFactor() float64 {
	g := s.physics.DecayTimeS * s.sampleRate / (s.physics.DecayTimeS*s.sampleRate + 6*math.Log(10)) // Stefan Bilbao's loss factor
	return g
}

// calculateVoiceLifetime determines the lifetime to give to the exported Voice in samples
func (s *stringSource) calculateVoiceLifetime() int {
	// return int(math.Round(s.physics.DecayTimeS)) * int(s.sampleRate)
	return int(s.sampleRate * 2)
}

// PublishVoice packages the synthesis function as createVoice struct and publishes it to the voiceChannel
func (s *stringSource) createVoice() *Voice {
	return &Voice{s.synthesize, 0, s.calculateVoiceLifetime(), false}
}

// synthesize simulates the state of the string at the next time stemp and generates an audio output sample
func (s *stringSource) synthesize() float32 {

	defer s.stepGrid()
	defer s.continuousExcitation([]float64{s.envelope.GetAmplitude()})
	defer s.envelope.Step()

	dt := 1 / s.sampleRate
	dt2 := math.Pow(dt, 2)
	a2 := math.Pow(s.physics.WaveSpeedMpS, 2)
	dx2 := math.Pow(s.stringLengthM/float64(s.numSpatialSections), 2)
	coeff := (dt2 * a2) / dx2
	g := s.calculateLossFactor()
	for m := 1; m < s.numSpatialSections; m++ {
		s.fdtdGrid[2][m] = g *
			(coeff*(s.fdtdGrid[1][m+1]-2*s.fdtdGrid[1][m]+s.fdtdGrid[1][m-1]) +
				2*s.fdtdGrid[1][m] - (2-1/g)*s.fdtdGrid[0][m])
	}

	return s.readPickup()

}

// readPickup is used by Synthesize to generate an output sample from a chosen point on the string
func (s *stringSource) readPickup() float32 {
	var pickupPoint int
	if s.physics.PickupPosReStringLen < 0.5 {
		pickupPoint = int(math.Ceil(float64(s.numSpatialSections) * s.physics.PickupPosReStringLen))
	} else {
		pickupPoint = int(math.Floor(float64(s.numSpatialSections) * s.physics.PickupPosReStringLen))
	}
	return float32(s.fdtdGrid[2][pickupPoint])
}

// stepGrid updates the finite difference simulation grid by one timestamp, providing a new, empty
// string state for simulation with Synthesize()
func (s *stringSource) stepGrid() {
	s.fdtdGrid = append(s.fdtdGrid, make([]float64, s.numSpatialSections+1))
	s.fdtdGrid = s.fdtdGrid[1:]
}

type pluckSettings struct {
	PosReStrLen   float64
	WidthReStrLen float64
	Amplitude     float64
}

type stringSettings struct {
	WaveSpeedMpS         float64
	DecayTimeS           float64
	PickupPosReStringLen float64
}

// newStringSource constructs a StringSource from the physical properties of a string
func newStringSource(sampleRate float64, lengthM float64, physics stringSettings,
	pluck Excitor, bow Excitor) stringSource {

	physics.PickupPosReStringLen = numeric.Clip(physics.PickupPosReStringLen, 0, 1)

	numSpatialSections := int(math.Floor(lengthM / (physics.WaveSpeedMpS * (1 / sampleRate)))) // Stability condition
	s := stringSource{
		newFtdtSource(sampleRate, 3, numSpatialSections),
		newEnvelopedSource(envelopes.NewTriangleEnvelope(time.Millisecond*1000, time.Second, sampleRate)),
		lengthM,
		physics,
		[]Excitor{pluck},
		[]Excitor{bow},
	}
	return s
}
