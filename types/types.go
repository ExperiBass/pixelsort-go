package types

type PixelWithMask struct {
	R, G, B, A uint8
	Mask       uint8
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
