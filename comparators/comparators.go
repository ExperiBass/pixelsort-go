package comparators

import (
	"pixelsort_go/shared"
	"pixelsort_go/types"
)

var ComparatorFunctionMappings = map[string]types.ComparatorFunc{
	"red":        Red,
	"green":      Green,
	"blue":       Blue,
	"hue":        Hue,
	"saturation": Saturation,
	"lightness":  Lightness,
	"darkness":   Darkness,
}

func Red(a, b types.PixelWithMask) int {
	if checkPixel(a) || checkPixel(b) {
		return 0
	}
	if checkPixelThresholds(int16(a.R)) || checkPixelThresholds(int16(b.R)) {
		return 0
	}
	return int(a.R) - int(b.R)
}

func Green(a, b types.PixelWithMask) int {
	if checkPixel(a) || checkPixel(b) {
		return 0
	}
	if checkPixelThresholds(int16(a.G)) || checkPixelThresholds(int16(b.G)) {
		return 0
	}
	return int(a.G) - int(b.G)
}

func Blue(a, b types.PixelWithMask) int {
	if checkPixel(a) || checkPixel(b) {
		return 0
	}
	if checkPixelThresholds(int16(a.B)) || checkPixelThresholds(int16(b.B)) {
		return 0
	}
	return int(a.B) - int(b.B)
}

func Hue(a, b types.PixelWithMask) int {
	if checkPixel(a) || checkPixel(b) {
		return 0
	}

	aHue := calculateHue(a)
	bHue := calculateHue(b)

	if checkPixelThresholds(int16(aHue)) || checkPixelThresholds(int16(bHue)) {
		return 0
	}
	return int(aHue - bHue)
}

func Saturation(a, b types.PixelWithMask) int {
	if checkPixel(a) || checkPixel(b) {
		return 0
	}

	aSat := calculateSaturation(a)
	bSat := calculateSaturation(b)

	if checkPixelThresholds(int16(aSat)) || checkPixelThresholds(int16(bSat)) {
		return 0
	}
	return int(aSat - bSat)
}

func Lightness(a, b types.PixelWithMask) int {
	if checkPixel(a) || checkPixel(b) {
		return 0
	}

	aLightness := calculateLightness(a)
	bLightness := calculateLightness(b)
	//println(aLightness, bLightness)
	if checkPixelThresholds(int16(aLightness)) || checkPixelThresholds(int16(bLightness)) {
		return 0
	}
	return int(aLightness - bLightness)
}

func Darkness(a, b types.PixelWithMask) int {
	return -Lightness(a, b)
}

// MAYBE: arbitrary thresholds (ex: green comparison with blue threhsolds)
func checkPixelThresholds(val int16) bool {
	// skip if beyond thresholds
	if val < int16(shared.Config.Thresholds.Lower*255) || val > int16(shared.Config.Thresholds.Upper*255) {
		return true
	}
	return false
}
func checkPixel(pixel types.PixelWithMask) bool {
	// skip if masked
	if pixel.Mask == 255 {
		return true
	}
	// and if null
	if pixel.R == 0 && pixel.G == 0 && pixel.B == 0 && pixel.A == 0 {
		return true
	}
	return false
}
func calculateLightness(pixel types.PixelWithMask) float32 {
	// 299, 587, 114
	return float32(pixel.R)*0.29 + float32(pixel.G)*0.59 + float32(pixel.B)*0.11
	//return int16(pixel.R)*29 + int16(pixel.G)*59 + int16(pixel.B)*11
}

func calculateHue(pixel types.PixelWithMask) int16 {
	hue := int16(0)
	maxV := max(pixel.R, pixel.G, pixel.B)
	minV := min(pixel.R, pixel.G, pixel.B)
	switch maxV {
	case pixel.R:
		{
			hue = int16(pixel.G - pixel.B)
			break
		}
	case pixel.G:
		{
			hue = 2 + int16(pixel.B-pixel.R)
			break
		}
	case pixel.B:
		{
			hue = 4 + int16(pixel.R-pixel.G)
		}
	}
	// finish formula and convert to degrees
	hue = (hue / (int16(maxV) - int16(minV))) * 60
	if hue < 0 {
		hue += 360
	}
	return hue
}
func calculateSaturation(pixel types.PixelWithMask) float32 {
	saturation := float32(0)
	// pixels are RGBA so skip the A
	minc := int(min(pixel.R, pixel.G, pixel.B))
	maxc := int(max(pixel.R, pixel.G, pixel.B))

	if minc == maxc {
		return saturation
	}

	sum := maxc + minc
	diff := maxc - minc
	lightness := float32((sum / 2) / 255)
	if lightness < 0.5 {
		saturation = float32(diff / sum)
	} else {
		saturation = float32(diff / (2 - diff))
	}
	return saturation
}
