package utils

import (
	"image"
	"math"
)

const precision = 2

func HsvToRgb(hue float64, saturation float64, value float64) (uint8, uint8, uint8) {
	r, g, b := hsv2rgb(hue, saturation, value)

	return uint8(r * 255.0), uint8(g * 255.0), uint8(b * 255.0)
}

func hsv2rgb(hueDegrees float64, saturation float64, value float64) (float64, float64, float64) {
	return hs2rgb(true, hueDegrees, saturation, value)
}

func hs2rgb(isValue bool, hueDegrees float64, saturation float64, lightOrVal float64) (float64, float64, float64) {
	var r, g, b float64

	hueDegrees = math.Mod(hueDegrees, 360)

	if saturation == 0 {
		r = lightOrVal
		g = lightOrVal
		b = lightOrVal
	} else {
		var chroma float64
		var m float64

		if isValue {
			chroma = lightOrVal * saturation
		} else {
			chroma = (1 - math.Abs((2*lightOrVal)-1)) * saturation
		}

		hueSector := hueDegrees / 60

		intermediate := chroma * (1 - math.Abs(
			math.Mod(hueSector, 2)-1,
		))

		switch {
		case hueSector >= 0 && hueSector <= 1:
			r = chroma
			g = intermediate
			b = 0

		case hueSector > 1 && hueSector <= 2:
			r = intermediate
			g = chroma
			b = 0

		case hueSector > 2 && hueSector <= 3:
			r = 0
			g = chroma
			b = intermediate

		case hueSector > 3 && hueSector <= 4:
			r = 0
			g = intermediate
			b = chroma
		case hueSector > 4 && hueSector <= 5:
			r = intermediate
			g = 0
			b = chroma

		case hueSector > 5 && hueSector <= 6:
			r = chroma
			g = 0
			b = intermediate

		default:
			// panic(fmt.Errorf("hue input %v yielded sector %v", hueDegrees, hueSector))
			r = chroma
			g = 0
			b = intermediate
		}

		if isValue {
			m = lightOrVal - chroma
		} else {
			m = lightOrVal - (chroma / 2)
		}

		r += m
		g += m
		b += m

	}

	r = RoundPlaces(r, precision)
	g = RoundPlaces(g, precision)
	b = RoundPlaces(b, precision)

	return r, g, b
}

func SumOfAllPixels(img image.Image) uint64 {
	var sum uint64
	bounds := img.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, _ := img.At(x, y).RGBA()

			sum += uint64(r) + uint64(g) + uint64(b)
			// fmt.Printf("x, y: %v, %v\n", x, y)
			// fmt.Printf("\tr: %v\n", r)
			// fmt.Printf("\tg: %v\n", g)
			// fmt.Printf("\tb: %v\n", b)
			// fmt.Printf("\tsum: %v\n", sum)
		}
	}

	return sum
}
