package types

type PixelWithMask struct {
	R, G, B, A uint8
	Mask       uint8
}

type Stretch struct {
	Start int
	End   int
}

type SortConfig struct {
	Sorter        string
	Comparator    string
	Quality       int
	SectionLength int
	Randomness    float32
	Reverse       bool
	Thresholds    ThresholdConfig
}

type ThresholdConfig struct {
	Lower, Upper float32
}

type ComparatorFunc func(a, b PixelWithMask) int

type SorterFunc func(interval []PixelWithMask)
