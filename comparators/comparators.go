package comparators

import (
	"pixelsort_go/shared"
	"pixelsort_go/types"
)

var ComparatorFunctionMappings = map[string]types.ComparatorFunc{
	"red":        Red,
	"green":      Green,
	"blue":       Blue,
	"saturation": Saturation,
	"lightness":  Lightness,
}

func Red(a, b types.PixelWithMask) int {
	if checkPixel(a) || checkPixel(b) {
		return 0
	}
	if checkPixelThresholds(float32(a.R)) || checkPixelThresholds(float32(b.R)) {
		return 0
	}
	return int(a.R) - int(b.R)
}
func Green(a, b types.PixelWithMask) int {
	if checkPixel(a) || checkPixel(b) {
		return 0
	}
	if checkPixelThresholds(float32(a.G)) || checkPixelThresholds(float32(b.G)) {
		return 0
	}
	return int(a.G) - int(b.G)
}

func Blue(a, b types.PixelWithMask) int {
	if checkPixel(a) || checkPixel(b) {
		return 0
	}
	if checkPixelThresholds(float32(a.B)) || checkPixelThresholds(float32(b.B)) {
		return 0
	}
	return int(a.B) - int(b.B)
}
func Saturation(a, b types.PixelWithMask) int {
	if checkPixel(a) || checkPixel(b) {
		return 0
	}

	aSat := calculateSaturation(a)
	bSat := calculateSaturation(b)

	if checkPixelThresholds(aSat) || checkPixelThresholds(bSat) {
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

	if checkPixelThresholds(aLightness) || checkPixelThresholds(bLightness) {
		return 0
	}
	return int(aLightness - bLightness)
}
func Darkness(a, b types.PixelWithMask) int {
	return -Lightness(a, b)
}

// MAYBE: arbitrary thresholds (ex: green comparison with blue threhsolds)
func checkPixelThresholds(val float32) bool {
	// skip if beyond thresholds
	if val < (shared.Config.Thresholds.Lower*255) || val > (shared.Config.Thresholds.Upper*255) {
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
	return float32(pixel.R)*0.299 + float32(pixel.G)*0.587 + float32(pixel.B)*0.114
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
