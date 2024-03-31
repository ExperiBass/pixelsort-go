package loaders

import (
	"image"
	"image/color"
	"pixelsort_go/types"
)

func LoadRow(img image.RGBA, mask image.RGBA) [][]types.PixelWithMask {
	dims := img.Bounds().Max
	// split image into rows
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
func LoadSpiral(img image.RGBA, mask image.RGBA) [][]types.PixelWithMask {
	dims := img.Bounds().Max
	seams := make([][]types.PixelWithMask, 0)

	for offset := 0; offset < dims.Y/2; offset++ {

		seam := make([]types.PixelWithMask, 0)

		for x := offset; x < dims.X-offset; x++ {
			///spirals = append(spirals, pixels[offset*dims.X+x])            // top
			///spirals = append(spirals, pixels[(dims.Y-offset-1)*dims.X+x]) // bottom
			topOffset := (offset*dims.X + x) * 4
			top := img.Pix[topOffset : topOffset+4]
			topMask := mask.Pix[topOffset+3]

			bottomOffset := ((dims.Y-offset-1)*dims.X + x) * 4
			bottom := img.Pix[bottomOffset : bottomOffset+4]
			bottomMask := mask.Pix[bottomOffset+3]

			seam = append(seam, types.PixelWithMask{R: top[0], G: top[1], B: top[2], A: top[3], Mask: topMask})
			seam = append(seam, types.PixelWithMask{R: bottom[0], G: bottom[1], B: bottom[2], A: bottom[3], Mask: bottomMask})
		}

		// right & left
		for y := offset + 1; y < dims.Y-offset-1; y++ {
			///spiral = append(spiral, pixels[y*dims.X+offset])        // right
			///spiral = append(spiral, pixels[y*dims.X+dims.X-offset]) // left
			rightOffset := (y*dims.X + offset) * 4
			right := img.Pix[rightOffset : rightOffset+4]
			rightMask := img.Pix[rightOffset+3]

			leftOffset := (y*dims.X + dims.X - offset) * 4
			left := img.Pix[leftOffset : leftOffset+4]
			leftMask := img.Pix[leftOffset+3]

			seam = append(seam, types.PixelWithMask{R: right[0], G: right[1], B: right[2], A: right[3], Mask: rightMask})
			seam = append(seam, types.PixelWithMask{R: left[0], G: left[1], B: left[2], A: left[3], Mask: leftMask})
		}

		seams = append(seams, seam)
	}

	return seams
}
func SaveSpiral(spirals [][]types.PixelWithMask, dims image.Rectangle) *image.RGBA {
	outputImg := image.NewRGBA(dims)

	width := dims.Max.X
	height := dims.Max.Y

	for offset := 0; offset < height/2; offset++ {
		seam := spirals[offset]
		index := 0

		// top
		for x := offset; x < width-offset; x++ {
			idx := (offset*width + x) * 4
			pixel := seam[index]
			outputImg.Pix[idx] = pixel.R
			outputImg.Pix[idx+1] = pixel.G
			outputImg.Pix[idx+2] = pixel.B
			outputImg.Pix[idx+3] = pixel.A
			println(idx)
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
			println(idx)
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
			println(idx)
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
			println(idx)
			index++
		}
	}
	println(len(outputImg.Pix))
	return outputImg
}
