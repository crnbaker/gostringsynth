package excitors

import "math"

type FDTDBow struct {
	PosReStrLen   float64
	WidthReStrLen float64
}

// bowString uses a stick-slip model to bowString an FDTD string. It must be called by synthesize.``
func (b *FDTDBow) Excite(source Excitable, amplitude float64) {
	numSpatialPoints := len(source.GetFDTDGrid()[0])
	bowVelocitycmps := amplitude * 75 // scale amplitude between 0 and 75 cm per second
	bowDisplacement := bowVelocitycmps / source.GetSampleRate()
	bowPosition := int(math.Floor(float64(numSpatialPoints-1) * b.PosReStrLen))
	stringDisplacementLastTimestep := source.GetFDTDGrid()[2][bowPosition] - source.GetFDTDGrid()[1][bowPosition]
	if stringDisplacementLastTimestep >= 0.0 && stringDisplacementLastTimestep < bowDisplacement {
		source.GetFDTDGrid()[2][bowPosition-1] = source.GetFDTDGrid()[2][bowPosition-1] + bowDisplacement
		source.GetFDTDGrid()[2][bowPosition] = source.GetFDTDGrid()[2][bowPosition] + bowDisplacement
		source.GetFDTDGrid()[2][bowPosition+1] = source.GetFDTDGrid()[2][bowPosition+1] + bowDisplacement

	}
}

func NewFTDTBow(PosReStrLen float64, WidthReStrLen float64) FDTDBow {
	return FDTDBow{PosReStrLen, WidthReStrLen}
}
