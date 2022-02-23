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
