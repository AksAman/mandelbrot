package mandelbrot

func mapRange(val, fromMin, fromMax, toMin, toMax float64) float64 {
	fromRange := fromMax - fromMin
	toRange := toMax - toMin
	scaleFactor := toRange / fromRange

	return toMin + ((val - fromMin) * scaleFactor)
}
