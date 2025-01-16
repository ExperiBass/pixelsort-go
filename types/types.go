package types

import "image/color"


type PixelWithMask struct {
	R, G, B, A uint8
	Mask       uint8
}
func (pixel PixelWithMask) ToColor() color.RGBA {
    return color.RGBA{
        R: pixel.R,
        G: pixel.G,
        B: pixel.B,
        A: pixel.A,
    }
}

type PixelStretch struct {
	Start int
	End   int
}

type ThresholdConfig struct {
	Lower, Upper float32
}

type ComparatorFunc func(a, b PixelWithMask) int

type SorterFunc func(interval []PixelWithMask)
