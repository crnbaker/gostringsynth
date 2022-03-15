package excitors

import (
	"math"

	"github.com/crnbaker/gostringsynth/numeric"
)

type FDTDPluck struct {
	PosReStrLen   float64
	WidthReStrLen float64
	Amplitude     float64
}

func NewFTDTPluck(PosReStrLen float64, WidthReStrLen float64, Amplitude float64) FDTDPluck {
	return FDTDPluck{PosReStrLen, WidthReStrLen, Amplitude}
}

func (p *FDTDPluck) Excite(fdtdGrid [][]float64) {
	numSpatialPoints := len(fdtdGrid[0])
	pluckShape := createTrianglePluck(p.Amplitude, numSpatialPoints, p.PosReStrLen)
	if p.WidthReStrLen < 1.0 {
		fingerWidthInSections := p.WidthReStrLen * float64(numSpatialPoints-1)
		fingerHalfWidthInPoints := int(math.Round(fingerWidthInSections+1) / 2)
		fingerWidthInPoints := fingerHalfWidthInPoints * 2

		if fingerWidthInPoints > 2 {
			var start int
			var stop int

			for i := fingerHalfWidthInPoints; i < numSpatialPoints-fingerHalfWidthInPoints; i++ {
				start = i - fingerHalfWidthInPoints
				stop = i + fingerHalfWidthInPoints
				pluckShape[i] = numeric.Mean(pluckShape[start:stop])
			}
		}

	}
	fdtdGrid[0] = pluckShape
	fdtdGrid[1] = pluckShape
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
