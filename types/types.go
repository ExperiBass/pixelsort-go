package types

type PixelWithMask struct {
	R, G, B, A uint8
	Mask       uint8
}

type PixelStretch struct {
	Start int
	End   int
}

type SortConfig struct {
	Pattern       string
	Interval      string
	Comparator    string
	SectionLength int
	Randomness    float32
	Reverse       bool
	Thresholds    ThresholdConfig
	Angle         float64
}

type ThresholdConfig struct {
	Lower, Upper float32
}

type ComparatorFunc func(a, b PixelWithMask) int

type SorterFunc func(interval []PixelWithMask)
