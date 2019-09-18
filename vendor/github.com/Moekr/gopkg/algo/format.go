package algo

import "math"

func Round(v float64, n int) float64 {
	p := math.Pow10(-n)
	return float64(int64((v+p/2)/p)) * p
}
