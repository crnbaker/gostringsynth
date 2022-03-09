package sources

import (
	"math"

	"github.com/crnbaker/gostringsynth/numeric"
)

// stringSource provides attributes that define a finite-difference simulation of a vibrating string
type stringSource struct {
	FDTDSource
	voiceSendChan chan Voice
	SampleRate    float64
	stringLengthM float64
	physics       StringSettings
	pluck         PluckSettings
}

// calculateLossFactor returns a loss factor used to attenuated the string vibration during synthesis
func (s *stringSource) calculateLossFactor() float64 {
	g := s.physics.DecayTimeS * s.SampleRate / (s.physics.DecayTimeS*s.SampleRate + 6*math.Log(10)) // Stefan Bilbao's loss factor
	return g
}

// calculateVoiceLifetime determines the lifetime to give to the exported Voice in samples
func (s *stringSource) calculateVoiceLifetime() int {
	return int(math.Round(s.physics.DecayTimeS)) * int(s.SampleRate)
}

// PublishVoice packages the synthesis function as Voice struct and publishes it to the voiceChannel
func (s *stringSource) PublishVoice() {
	s.voiceSendChan <- Voice{s.Synthesize, 0, s.calculateVoiceLifetime(), false}
}

// Synthesize simulates the state of the string at the next time stemp and generates an audio output sample
func (s *stringSource) Synthesize() (sampleValue float32) {
	defer s.stepGrid()
	dt2 := math.Pow(1/s.SampleRate, 2)
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

type PluckSettings struct {
	PosReStrLen   float64
	WidthReStrLen float64
	Amplitude     float64
}

func (s *stringSource) SoftPluck() []float64 {
	pluckShape := createTrianglePluck(s.pluck.Amplitude, s.numSpatialSections+1, s.pluck.PosReStrLen)
	if s.pluck.WidthReStrLen < 1.0 {
		stringLengthInPoints := s.numSpatialSections + 1
		fingerWidthInSections := s.pluck.WidthReStrLen * float64(s.numSpatialSections)
		fingerHalfWidthInPoints := int(math.Round(fingerWidthInSections+1) / 2)
		fingerWidthInPoints := fingerHalfWidthInPoints * 2

		if fingerWidthInPoints > 2 {
			var start int
			var stop int

			for i := fingerHalfWidthInPoints; i < stringLengthInPoints-fingerHalfWidthInPoints; i++ {
				start = i - fingerHalfWidthInPoints
				stop = i + fingerHalfWidthInPoints
				pluckShape[i] = mean(pluckShape[start:stop])
			}
		}

	}
	s.fdtdGrid[0] = pluckShape
	s.fdtdGrid[1] = pluckShape
	return s.fdtdGrid[0]
}

type StringSettings struct {
	WaveSpeedMpS         float64
	DecayTimeS           float64
	PickupPosReStringLen float64
}

// NewStringSource constructs a StringSource from the physical properties of a string
func NewStringSource(sampleRate float64, voiceSendChan chan Voice, lengthM float64, physics StringSettings,
	pluck PluckSettings) stringSource {

	physics.PickupPosReStringLen = numeric.Clip(physics.PickupPosReStringLen, 0, 1)
	pluck.PosReStrLen = numeric.Clip(pluck.PosReStrLen, 0, 1)
	pluck.WidthReStrLen = numeric.Clip(pluck.WidthReStrLen, 0, 1)

	numSpatialSections := int(math.Floor(lengthM / (physics.WaveSpeedMpS * (1 / sampleRate)))) // Stability condition
	s := stringSource{
		NewFTDTSource(3, numSpatialSections),
		voiceSendChan,
		sampleRate,
		lengthM,
		physics,
		pluck,
	}
	return s
}

func mean(slice []float64) float64 {
	var sum float64 = slice[0]
	for _, value := range slice {
		sum += value
	}
	return sum / float64(len(slice))
}

// trianglePluck creates a trianglePluck shape in a slice
func createTrianglePluck(amplitude float64, length int, pluckPosFraction float64) []float64 {
	pluckPoint := int(math.Floor(float64(length) * pluckPosFraction))
	if pluckPoint < 1 {
		pluckPoint = 1
	} else if pluckPoint >= length {
		pluckPoint = length - 1
	}
	pluck := make([]float64, length)
	for point := 0; point <= pluckPoint; point++ {
		pluck[point] = amplitude * float64(point) / float64(pluckPoint)
	}
	for point := pluckPoint; point < length; point++ {
		pluck[point] = amplitude * float64(length-point-1) / float64(length-pluckPoint-1)
	}
	return pluck
}
