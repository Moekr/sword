package algo

import "math"

func Var(vs []float64) float64 {
	if len(vs) == 0 {
		return 0
	}
	var sum float64
	for _, v := range vs {
		sum = sum + v
	}
	avg, sum := sum / float64(len(vs)), 0
	for _, v := range vs {
		sum = sum + (avg-v)*(avg-v)
	}
	return sum / float64(len(vs))
}

func Std(vs []float64) float64 {
	return math.Sqrt(Var(vs))
}
