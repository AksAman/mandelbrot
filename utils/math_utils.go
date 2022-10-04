package utils

import "math"

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func Clamp(x, min, max int) int {
	return Max(min, Min(x, max))
}

func ClampFloat(x, min, max float64) float64 {
	return math.Max(min, math.Min(x, max))
}

func RoundPlaces(x float64, places int) float64 {
	if places <= 0 {
		return math.Round(x)
	} else {
		multi := math.Pow(10, ClampFloat(float64(places), 0, 16))
		return (math.Round(x*multi) / multi)
	}
}
