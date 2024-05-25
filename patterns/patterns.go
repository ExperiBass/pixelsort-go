package patterns

import (
	"image"
	"image/color"
	"pixelsort_go/types"
)

var Loader = map[string]func(img image.RGBA, mask image.RGBA) [][]types.PixelWithMask{
	/// theres a better way, right? please tell me im dumb
	"rowload":    LoadRow,
	"spiralload": LoadSpiral,
}
var Saver = map[string]func(rows [][]types.PixelWithMask, dims image.Rectangle) *image.RGBA{
	"rowsave":    SaveRow,
	"spiralsave": SaveSpiral,
}

func LoadRow(img image.RGBA, mask image.RGBA) [][]types.PixelWithMask {
	dims := img.Bounds().Max
	/// split image into rows
	rows := make([][]types.PixelWithMask, dims.Y)
	for y := 0; y < dims.Y; y++ {
		row := make([]types.PixelWithMask, dims.X)

		for x := 0; x < dims.X; x++ {
			pixel := img.RGBAAt(x, y)
			masked := mask.RGBAAt(x, y).R
			wrapped := types.PixelWithMask{R: pixel.R, G: pixel.G, B: pixel.B, A: pixel.A, Mask: masked}
			row[x] = wrapped
		}
		rows[y] = row
	}
	return rows
}
func SaveRow(rows [][]types.PixelWithMask, dims image.Rectangle) *image.RGBA {
	outputImg := image.NewRGBA(dims)
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		for j := 0; j < len(row); j++ {
			currPixWithMask := row[j]
			pixel := color.RGBA{currPixWithMask.R, currPixWithMask.G, currPixWithMask.B, currPixWithMask.A}
			outputImg.SetRGBA(j, i, pixel)
		}
	}
	return outputImg
}

// https://github.com/jeffThompson/PixelSorting/blob/master/SpiralSortPixels/SpiralSortPixels.pde
// prayge, i'm not a mathy fomx
// this also half-works; when giving a mask, or skipping the sorting step, the image comes out
// missing half its pixels in all but one direction (usually the left side) and the
// remaining three parts are flipped
// only on images with even dimensions tho!!!! wheeeee!!!!!!!!!
func LoadSpiral(img image.RGBA, mask image.RGBA) [][]types.PixelWithMask {
	dims := img.Bounds().Max
	width := dims.X
	height := dims.Y

	seams := make([][]types.PixelWithMask, 0)

	for offset := 0; offset < height/2; offset++ {

		seam := make([]types.PixelWithMask, 0)

		for x := offset; x < width-offset; x++ {
			topOffset := (offset*width + x) * 4
			top := img.Pix[topOffset : topOffset+4]
			topMask := mask.Pix[topOffset]

			bottomOffset := ((height-offset-1)*width + x) * 4
			bottom := img.Pix[bottomOffset : bottomOffset+4]
			bottomMask := mask.Pix[bottomOffset]

			seam = append(seam, types.PixelWithMask{R: top[0], G: top[1], B: top[2], A: top[3], Mask: topMask})
			seam = append(seam, types.PixelWithMask{R: bottom[0], G: bottom[1], B: bottom[2], A: bottom[3], Mask: bottomMask})
		}

		// right & left
		for y := offset + 1; y < height-offset-1; y++ {
			rightOffset := (y*width + offset) * 4
			right := img.Pix[rightOffset : rightOffset+4]
			rightMask := mask.Pix[rightOffset]

			leftOffset := (y*width + width - offset) * 4
			left := img.Pix[leftOffset : leftOffset+4]
			leftMask := mask.Pix[leftOffset]

			seam = append(seam, types.PixelWithMask{R: right[0], G: right[1], B: right[2], A: right[3], Mask: rightMask})
			seam = append(seam, types.PixelWithMask{R: left[0], G: left[1], B: left[2], A: left[3], Mask: leftMask})
		}

		seams = append(seams, seam)
	}

	return seams
}
func SaveSpiral(seams [][]types.PixelWithMask, dims image.Rectangle) *image.RGBA {
	outputImg := image.NewRGBA(dims)

	width := dims.Max.X
	height := dims.Max.Y

	for offset := 0; offset < height/2; offset++ {
		seam := seams[offset]
		index := 0

		// top
		for x := offset; x < width-offset; x++ {
			idx := (offset*width + x) * 4
			pixel := seam[index]
			outputImg.Pix[idx] = pixel.R
			outputImg.Pix[idx+1] = pixel.G
			outputImg.Pix[idx+2] = pixel.B
			outputImg.Pix[idx+3] = pixel.A
			index++
		}

		// right
		for y := offset + 1; y < height-offset-1; y++ {
			idx := (y*width + width - offset) * 4
			pixel := seam[index]
			outputImg.Pix[idx] = pixel.R
			outputImg.Pix[idx+1] = pixel.G
			outputImg.Pix[idx+2] = pixel.B
			outputImg.Pix[idx+3] = pixel.A
			index++
		}

		// bottom
		for x := width - offset - 1; x >= offset; x-- {
			idx := ((height-offset-1)*width + x) * 4
			pixel := seam[index]
			outputImg.Pix[idx] = pixel.R
			outputImg.Pix[idx+1] = pixel.G
			outputImg.Pix[idx+2] = pixel.B
			outputImg.Pix[idx+3] = pixel.A
			index++
		}

		// left
		for y := height - offset - 2; y > offset+1; y-- {
			idx := (y*width + offset) * 4
			pixel := seam[index]
			outputImg.Pix[idx] = pixel.R
			outputImg.Pix[idx+1] = pixel.G
			outputImg.Pix[idx+2] = pixel.B
			outputImg.Pix[idx+3] = pixel.A
			index++
		}
	}

	return outputImg
}

/*
func LoadSeamCarving(img image.RGBA, mask image.RGBA) [][]types.PixelWithMask {
	dims := img.Bounds().Max
	width := dims.X
	height := dims.Y

	seams := make([][]types.PixelWithMask, 0)
	// TODO
	return seams
}
func SaveSeamCarving(seams [][]types.PixelWithMask, dims image.Rectangle) *image.RGBA {
	outputImg := image.NewRGBA(dims)

	width := dims.Max.X
	height := dims.Max.Y
	// TODO
	return outputImg
}
*/
