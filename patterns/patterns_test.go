package patterns_test

import (
	//"fmt"
	"math/rand"
	"testing"

	"image"
	"image/color"
	"pixelsort_go/patterns"
	// "pixelsort_go/types"
)

// func TestLoadSpiral(t *testing.T) {
// 	input := image.NewRGBA(image.Rectangle{
// 		Min: image.Point{X: 0, Y: 0},
// 		Max: image.Point{X: 3, Y: 3},
// 	})
// 	mask := image.NewRGBA(image.Rectangle{
// 		Min: image.Point{X: 0, Y: 0},
// 		Max: image.Point{X: 3, Y: 3},
// 	})
// 	colors := make([]color.RGBA, 6)
// 	colors[0] = color.RGBA{
// 		R: 255,
// 		G: 255,
// 		B: 255,
// 		A: 255,
// 	}
// 	colors[1] = color.RGBA{
// 		R: 255,
// 		G: 0,
// 		B: 255,
// 		A: 255,
// 	}
// 	colors[2] = color.RGBA{
// 		R: 255,
// 		G: 255,
// 		B: 0,
// 		A: 255,
// 	}
// 	colors[3] = color.RGBA{
// 		R: 0,
// 		G: 255,
// 		B: 255,
// 		A: 255,
// 	}
// 	colors[4] = color.RGBA{
// 		R: 0,
// 		G: 0,
// 		B: 255,
// 		A: 255,
// 	}

// 	// expected := []uint8{
// 	// 	0,0,0,0, 0,0,0,0, 255,0,255,255,
// 	//  0,255,255,255, 255,255,0,255, 255,255,0,255,
// 	//  0,0,0,0, 0,255,255,255, 255,255,0,255
// 	// }
// 	expected := [][]types.PixelWithMask{
// 		[]types.PixelWithMask{
// 			types.PixelWithMask{R:0,G:0,B:0,A:0,Mask:0},
// 			types.PixelWithMask{R:0,G:0,B:0,A:0,Mask:0},
// 			types.PixelWithMask{R:255,G:0,B:255,A:255,Mask:0},

// 			types.PixelWithMask{R:0,G:255,B:255,A:255,Mask:0},
// 			types.PixelWithMask{R:255,G:255,B:0,A:255,Mask:0},
// 			types.PixelWithMask{R:255,G:255,B:0,A:255,Mask:0},

// 			types.PixelWithMask{R:0,G:0,B:0,A:0,Mask:0},
// 			types.PixelWithMask{R:0,G:255,B:255,A:255,Mask:0},
// 		},
// 		[]types.PixelWithMask{
// 			types.PixelWithMask{R:255,G:255,B:0,A:255,Mask:0},
// 		},
// 	}
// 	/// fuck this shit
// 	for x := 1; x < 3; x++ {
// 		for y := 0; y < 3; y++ {
// 			input.Set(x, y, colors[(x+y)%4])
// 		}
// 	}
// 	slices, _ := patterns.LoadSpiral(input, mask)
// 	fmt.Println(expected)
// 	fmt.Println(*slices)
// 	for i := range *slices {
// 		slice := (*slices)[i]
// 		for ii := range slice {
// 			if (*slices)[i][ii] != expected[i][ii] {
// 				fmt.Println(i, ii, (*slices)[i][ii], expected[i][ii])
// 				t.FailNow()
// 			}
// 		}
// 	}
// 	// res := patterns.SaveSpiral(slices, input.Rect)
// 	// fmt.Println(input.Pix)
// 	// fmt.Println(res.Pix)
// 	// for idx := range input.Pix {
// 	// 	if input.Pix[idx] != res.Pix[idx] {
// 	// 		fmt.Println(idx)
// 	// 		t.FailNow()
// 	// 		break
// 	// 	}
// 	// }
// }

func TestArdenLoadSpiral(t *testing.T) {
    // Constants are not to be modified; test expects 3 and 3
    WIDTH := 3
    HEIGHT := 3
    input := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
    mask := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
    // Fill input with random pixels
    n, err := rand.Read(input.Pix)
    if n != WIDTH*HEIGHT*4 {
        t.Errorf("Read %d bytes, expected %d", n, WIDTH*HEIGHT*4)
    } else if err != nil {
        t.Errorf("Error reading: %s", err)
    }
    // Algorithm reads top then bottom from left to right, then right then left from top to bottom
    // Each spiral is a separate slice
    expected := [][]color.RGBA{
        {
            input.RGBAAt(0, 0), input.RGBAAt(0, 2), input.RGBAAt(1, 0),
            input.RGBAAt(1, 2), input.RGBAAt(2, 0), input.RGBAAt(2, 2),
            input.RGBAAt(2, 1), input.RGBAAt(0, 1),
        },
        {
            input.RGBAAt(1, 1),
        },
    }
    actual, _ := patterns.LoadSpiral(input, mask)
    t.Logf("input: %v", input.Pix)
    t.Logf("expected: %v", expected)
    t.Logf("actual: %v", actual)
    // Compare equality of each element in each slice
    for slice := 0; slice < len(expected); slice++ {
        for pixel := 0; pixel < len(expected[slice]); pixel++ {
            colorExpected := expected[slice][pixel]
            colorActual := (*actual)[slice][pixel].ToColor()
            if colorActual != colorExpected {
                t.Errorf("Pixel %d of slice %d is inequal. Expected %v, got %v", pixel, slice, colorExpected, colorActual)
            }
        }
    }
}
