package algo

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MaxI64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func MaxF64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MinI64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func MinF64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
