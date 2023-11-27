package sorters

import (
	"math"
	mathRand "math/rand"
	"pixelsort_go/comparators"
	"pixelsort_go/shared"
	"pixelsort_go/types"
	"slices"
)

var SortingFunctionMappings = map[string]func([]types.PixelWithMask){
	"shuffle":     Shuffle,
	"row":         Row,
	"random":      Random,
	"randomnoisy": RandomNoisy,
	"wave":        Wave,
}

func Sort(row []types.PixelWithMask) {
	stretches := getUnmaskedStretches(row)
	for i := 0; i < len(stretches); i++ {
		stretch := stretches[i]
		SortingFunctionMappings[shared.Config.Sorter](row[stretch.Start:stretch.End])
	}
}

// sorters
func Shuffle(interval []types.PixelWithMask) {
	mathRand.Shuffle(len(interval), func(i, j int) {
		interval[i], interval[j] = interval[j], interval[i]
	})
}

func Row(interval []types.PixelWithMask) {
	if mathRand.Float32() < shared.Config.Randomness {
		return
	}
	commonSort([]types.Stretch{{Start: 0, End: len(interval)}}, interval)
}

func Random(interval []types.PixelWithMask) {
	stretches := make([]types.Stretch, 0)
	intervalLength := len(interval)
	section_length := shared.Config.SectionLength

	j := 0
	for {
		if j >= intervalLength {
			break
		}
		// so many broken limbs, call me the cast fox
		stretchLength := max(mathRand.Int(), int(math.Floor(float64(float32(section_length)*shared.Config.Randomness))))
		endIdx := min(j+stretchLength, intervalLength)
		stretches = append(stretches, types.Stretch{Start: j, End: endIdx})
		j += stretchLength
	}
	commonSort(stretches, interval)
}
func RandomNoisy(interval []types.PixelWithMask) {
	stretches := make([]types.Stretch, 0)
	intervalLength := len(interval)
	min_section_length := shared.Config.SectionLength

	j := 0
	for {
		if j >= intervalLength {
			break
		}
		randLength := randBetween((intervalLength - j), min_section_length)
		endIdx := min(j+randLength, intervalLength)
		stretches = append(stretches, types.Stretch{Start: j, End: endIdx})
		j += randLength
	}
	commonSort(stretches, interval)
}

func Wave(interval []types.PixelWithMask) {
	stretches := make([]types.Stretch, 0)
	intervalLength := len(interval)
	sectionLength := shared.Config.SectionLength

	j := 0
	for {
		if j >= intervalLength {
			break
		}
		randFloat := mathRand.Float64() * 100
		randMin := min(randFloat, math.Floor(float64(float32(sectionLength)*shared.Config.Randomness)))

		stretchLength := sectionLength + randBetween(int(randFloat), int(-randMin))

		endIdx := min(j+stretchLength, intervalLength)
		stretches = append(stretches, types.Stretch{Start: j, End: endIdx})
		j += stretchLength
	}
	commonSort(stretches, interval)
}

///

///

// util
func randBetween(max int, min_opt ...int) int {
	min := 0
	if len(min_opt) > 0 {
		min = min_opt[0]
	}
	randNum := mathRand.Float64()
	if min != 0 {
		return int(math.Floor(randNum*float64(((+max)+1)-(+min)))) + (+min)
	}
	return int(math.Floor(randNum * float64((+max)+1)))
}

func commonSort(stretches []types.Stretch, interval []types.PixelWithMask) {
	for stretchIdx := 0; stretchIdx < len(stretches); stretchIdx++ {
		stretch := stretches[stretchIdx]
		pixels := interval[stretch.Start:stretch.End]
		if shared.Config.Reverse {
			for i, j := 0, len(pixels)-1; i < j; i, j = i+1, j-1 {
				pixels[i], pixels[j] = pixels[j], pixels[i]
			}
		}

		slices.SortStableFunc(pixels, comparators.ComparatorFunctionMappings[shared.Config.Comparator])

		if shared.Config.Reverse {
			for i, j := 0, len(pixels)-1; i < j; i, j = i+1, j-1 {
				pixels[i], pixels[j] = pixels[j], pixels[i]
			}
		}

	}
}

func getUnmaskedStretches(interval []types.PixelWithMask) []types.Stretch {
	stretches := make([]types.Stretch, 0)
	baseIdx := 0

	for j := 0; j < len(interval); j++ {
		pixel := interval[j]
		if pixel.Mask == 255 || (pixel.R == 0 && pixel.G == 0 && pixel.B == 0 && pixel.A == 0) {
			// look ahead for the end of the mask
			endMaskIdx := j
			for {
				if endMaskIdx >= len(interval) {
					break
				}
				if interval[endMaskIdx].Mask != 255 || !(pixel.R == 0 && pixel.G == 0 && pixel.B == 0 && pixel.A == 0) {
					//println(interval[endMaskIdx-1].Mask, interval[endMaskIdx].Mask, interval[endMaskIdx+1].Mask)
					break
				}
				endMaskIdx++
			}
			stretch := types.Stretch{Start: baseIdx, End: j}
			//stretches[len(stretches)] = stretch
			stretches = append(stretches, stretch)

			baseIdx = endMaskIdx
			j = baseIdx
		}
	}
	// and then add any remaning unmasked pixels
	//stretches[0] = types.Stretch{Start: baseIdx, End: len(interval)}
	stretches = append(stretches, types.Stretch{Start: baseIdx, End: len(interval)})
	//fmt.Printf("stretches: %v\n", stretches)
	return stretches
}
