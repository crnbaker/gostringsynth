package sources

import (
	"math"
)

type stringVoiceSource struct {
	VoiceSource
	FDTDSource
	SampleRate           float64
	stringLengthM        float64
	stringWaveSpeedMPerS float64
	pickupPositionFrac   float64
	decayTimeS           float64
}

// Bilbao's loss factor
func (s *stringVoiceSource) calculateLossFactor() float64 {
	g := s.decayTimeS * s.SampleRate / (s.decayTimeS*s.SampleRate + 6*math.Log(10))
	return g
}

func (s *stringVoiceSource) calculateVoiceLifetime() int {
	return int(math.Round(s.decayTimeS*1.1)) * int(s.SampleRate)
}

func (s *stringVoiceSource) DispatchAndPlayVoice(freqHz float64, amplitude float64) {
	s.voiceSendChan <- Voice{s.Synthesize, 0, s.calculateVoiceLifetime()}
	s.pluck(amplitude)
}

func (s *stringVoiceSource) Synthesize() (sampleValue float32) {
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

func (s *stringVoiceSource) readPickup() float32 {
	var pickupPoint int
	if s.pickupPositionFrac < 0.5 {
		pickupPoint = int(math.Ceil(float64(s.numSpatialSections) * s.pickupPositionFrac))
	} else {
		pickupPoint = int(math.Floor(float64(s.numSpatialSections) * s.pickupPositionFrac))
	}
	return float32(s.fdtdGrid[2][pickupPoint])
}

func (s *stringVoiceSource) stepGrid() {
	s.fdtdGrid = append(s.fdtdGrid, make([]float64, s.numSpatialSections+1))
	s.fdtdGrid = s.fdtdGrid[1:]
}

func (s *stringVoiceSource) pluck(amplitude float64) {
	pluckPoint := int(math.Floor(float64(len(s.fdtdGrid[0])) / 2))
	for i := 0; i <= pluckPoint; i++ {
		s.fdtdGrid[0][i] = amplitude * float64(i) / float64(pluckPoint)
		s.fdtdGrid[1][i] = amplitude * float64(i) / float64(pluckPoint)
	}
	for point, i := len(s.fdtdGrid[0])-1, 0; point > pluckPoint; point, i = point-1, i+1 {
		s.fdtdGrid[0][point] = amplitude * float64(i) / float64(pluckPoint)
		s.fdtdGrid[1][point] = amplitude * float64(i) / float64(pluckPoint)
		i++
	}
}

func NewStringVoiceSource(sampleRate float64, voiceSendChan chan Voice, lengthM float64, waveSpeedMpS float64,
	pickupPosFrac float64, decayTimeS float64) stringVoiceSource {

	if pickupPosFrac < 0 {
		pickupPosFrac = 0
	} else if pickupPosFrac > 1 {
		pickupPosFrac = 1
	}

	numSpatialSections := int(math.Floor(lengthM / (waveSpeedMpS * (1 / sampleRate)))) // Stability condition
	s := stringVoiceSource{
		NewVoiceSource(voiceSendChan),
		NewFTDTSource(3, numSpatialSections),
		sampleRate,
		lengthM,
		waveSpeedMpS,
		pickupPosFrac,
		decayTimeS,
	}
	return s
}

func FreqToStringLength(freqHz float64, waveSpeedMpS float64) float64 {
	return waveSpeedMpS / freqHz
}
