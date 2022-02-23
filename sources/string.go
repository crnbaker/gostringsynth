package sources

import (
	"math"
)

// stringSource provides attributes that define a finite-difference simulation of a vibrating string
type stringSource struct {
	FDTDSource
	voiceSendChan        chan Voice
	SampleRate           float64
	stringLengthM        float64
	stringWaveSpeedMPerS float64
	pickupPositionFrac   float64
	decayTimeS           float64
}

// calculateLossFactor returns a loss factor used to attenuated the string vibration during synthesis
func (s *stringSource) calculateLossFactor() float64 {
	g := s.decayTimeS * s.SampleRate / (s.decayTimeS*s.SampleRate + 6*math.Log(10)) // Stefan Bilbao's loss factor
	return g
}

// calculateVoiceLifetime determines the lifetime to give to the exported Voice in samples
func (s *stringSource) calculateVoiceLifetime() int {
	return int(math.Round(s.decayTimeS)) * int(s.SampleRate)
}

// PublishVoice packages the synthesis function as Voice struct and publishes it to the voiceChannel
func (s *stringSource) PublishVoice() {
	s.voiceSendChan <- Voice{s.Synthesize, 0, s.calculateVoiceLifetime(), false}
}

// Synthesize simulates the state of the string at the next time stemp and generates an audio output sample
func (s *stringSource) Synthesize() (sampleValue float32) {
	defer s.stepGrid()
	dt2 := math.Pow(1/s.SampleRate, 2)
	a2 := math.Pow(s.stringWaveSpeedMPerS, 2)
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
	if s.pickupPositionFrac < 0.5 {
		pickupPoint = int(math.Ceil(float64(s.numSpatialSections) * s.pickupPositionFrac))
	} else {
		pickupPoint = int(math.Floor(float64(s.numSpatialSections) * s.pickupPositionFrac))
	}
	return float32(s.fdtdGrid[2][pickupPoint])
}

// stepGrid updates the finite difference simulation grid by one timestamp, providing a new, empty
// string state for simulation with Synthesize()
func (s *stringSource) stepGrid() {
	s.fdtdGrid = append(s.fdtdGrid, make([]float64, s.numSpatialSections+1))
	s.fdtdGrid = s.fdtdGrid[1:]
}

func (s *stringSource) SoftPluck(amplitude float64) []float64 {
	const fingerWidthM = 0.05
	const pluckPosition = 0.5
	pluck := createTrianglePluck(amplitude, s.numSpatialSections+1, pluckPosition)
	if fingerWidthM < s.stringLengthM {
		dx := s.stringLengthM / float64(s.numSpatialSections)
		stringLengthInPoints := s.numSpatialSections + 1
		fingerWidthInSections := fingerWidthM / dx
		fingerHalfWidthInPoints := int(math.Round(fingerWidthInSections+1) / 2)
		fingerWidthInPoints := fingerHalfWidthInPoints * 2

		if fingerWidthInPoints > 2 {
			var start int
			var stop int

			for i := fingerHalfWidthInPoints; i < stringLengthInPoints-fingerHalfWidthInPoints; i++ {
				start = i - fingerHalfWidthInPoints
				stop = i + fingerHalfWidthInPoints
				pluck[i] = mean(pluck[start:stop])
			}
		}

	}
	s.fdtdGrid[0] = pluck
	s.fdtdGrid[1] = pluck
	return s.fdtdGrid[0]
}

// NewStringSource constructs a StringSource from the physical properties of a string
func NewStringSource(sampleRate float64, voiceSendChan chan Voice, lengthM float64, waveSpeedMpS float64,
	pickupPosFrac float64, decayTimeS float64) stringSource {

	if pickupPosFrac < 0 {
		pickupPosFrac = 0
	} else if pickupPosFrac > 1 {
		pickupPosFrac = 1
	}

	numSpatialSections := int(math.Floor(lengthM / (waveSpeedMpS * (1 / sampleRate)))) // Stability condition
	s := stringSource{
		NewFTDTSource(3, numSpatialSections),
		voiceSendChan,
		sampleRate,
		lengthM,
		waveSpeedMpS,
		pickupPosFrac,
		decayTimeS,
	}
	return s
}

// FreqToStringLength converts a fundemental frequency to a string length, given a string wave speed in m/s
func FreqToStringLength(freqHz float64, waveSpeedMpS float64) float64 {
	return waveSpeedMpS / freqHz
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
	pluck := make([]float64, length)
	for point := 0; point <= pluckPoint; point++ {
		pluck[point] = amplitude * float64(point) / float64(pluckPoint)
	}
	for point := pluckPoint; point < length; point++ {
		pluck[point] = amplitude * float64(length-point-1) / float64(length-pluckPoint-1)
	}
	return pluck
}
