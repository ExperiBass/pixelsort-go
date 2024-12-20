package shared

/// global var to hold config instead of passing it everywhere
import (
	"pixelsort_go/types"
)

var Config struct {
	Pattern       string
	Interval      string
	Comparator    string
	SectionLength int
	Randomness    float32
	Reverse       bool
	Thresholds    types.ThresholdConfig
	Angle         float64
}
