package numeric

func Max(s []float64) (max float64) {
	max = s[0]
	for _, value := range s {
		if max < value {
			max = value
		}
	}
	return
}

func Clip(value float64, min float64, max float64) float64 {
	if value < 0 {
		return 0
	} else if value > 1 {
		return 1
	} else {
		return value
	}
}
