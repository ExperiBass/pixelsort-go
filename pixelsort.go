package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"pixelsort_go/shared"
	"pixelsort_go/sorters"
	"pixelsort_go/types"
	"slices"
	"strings"

	"github.com/remeh/sizedwaitgroup"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "`image`(s) to sort",
				Required: true,
				Action: func(ctx *cli.Context, v []string) error {
					if strings.HasSuffix(v[0], "jpg") || strings.HasSuffix(v[0], "jpeg") {
						image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
					} else if strings.HasSuffix(v[0], "png") {
						image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
						return nil
					} else {
						return fmt.Errorf("input is an invalid image (supported: jpg, png)")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "out",
				Value:    "out.png",
				Aliases:  []string{"o"},
				Usage:    "`file` to output to",
				Required: false,
			},
			&cli.StringFlag{
				Name:    "mask",
				Aliases: []string{"m"},
				Usage:   "b&w `mask` to lay over the image; white is skipped",
			},
			&cli.StringFlag{
				Name:    "sorter",
				Value:   "row",
				Aliases: []string{"s"},
				Usage:   "sorting `algorithm` to use",
			},
			&cli.StringFlag{
				Name:    "comparator",
				Value:   "lightness",
				Aliases: []string{"c"},
				Usage:   "comparison `function` to use",
			},
			&cli.Float64Flag{
				Name:    "lower_threshold",
				Value:   0.0,
				Aliases: []string{"l"},
				Usage:   "pixels below this `threshold` won't be sorted",
				Action: func(ctx *cli.Context, v float64) error {
					if v < 0.0 || v > 1.0 {
						return fmt.Errorf("lower_threshold is out of range [0.0-1.0]")
					}
					return nil
				},
			},
			&cli.Float64Flag{
				Name:    "upper_threshold",
				Value:   1.0,
				Aliases: []string{"u"},
				Usage:   "pixels above this `threshold` won't be sorted",
				Action: func(ctx *cli.Context, v float64) error {
					if v < 0.0 || v > 1.0 {
						return fmt.Errorf("upper_threshold is out of range [0.0-1.0]")
					}
					return nil
				},
			},
			&cli.IntFlag{
				Name:    "quality",
				Value:   100,
				Aliases: []string{"q"},
				Usage:   "(jpeg only) the `quality` of the output image.",
				Action: func(ctx *cli.Context, v int) error {
					if v < 20 || v > 100 {
						return fmt.Errorf("quality is out of range [20-100]")
					}
					return nil
				},
			},
			&cli.IntFlag{
				Name:    "section_length",
				Value:   60,
				Aliases: []string{"L"},
				Usage:   "The base `length` of each slice",
			},
			&cli.BoolFlag{
				Name:  "reverse",
				Value: false,
				Usage: "reverse sort, or not.",
			},
			&cli.Float64Flag{
				Name:  "randomness",
				Value: 1.0,
				Usage: "used to determine which [row]s to skip and how wild [wave] edges should be",
				Action: func(ctx *cli.Context, v float64) error {
					if v < 0.0 || v > 1.0 {
						return fmt.Errorf("randomness is out of range [0.0-1.0]")
					}
					return nil
				},
			},
			&cli.IntFlag{
				Name:    "threads",
				Value:   1,
				Aliases: []string{"t"},
				Usage:   "Sort images in parallel across `N` threads",
			},
		},
		Name:  "pixelsort_go",
		Usage: "Visual decimation.",
		Action: func(ctx *cli.Context) error {
			inputs := ctx.StringSlice("input")
			threadCount := ctx.Int("threads")
			shared.Config.Mask = ctx.String("mask")
			shared.Config.Sorter = ctx.String("sorter")
			shared.Config.Comparator = ctx.String("comparator")
			shared.Config.Quality = ctx.Int("quality")
			shared.Config.SectionLength = ctx.Int("section_length")
			shared.Config.Thresholds.Lower = float32(ctx.Float64("lower_threshold"))
			shared.Config.Thresholds.Upper = float32(ctx.Float64("upper_threshold"))
			shared.Config.Reverse = ctx.Bool("reverse")
			shared.Config.Randomness = float32(ctx.Float64("randomness"))

			/// lets make it noisy
			infoString := fmt.Sprintf("Sorting %d images with a config of %+v.", len(inputs), shared.Config)
			println(infoString)
			if len(inputs) == 1 {
				return sortingTime(inputs[0], ctx.String(("out")), shared.Config.Mask)
			}

			// multiple imgs
			// sort em first so frames dont get jumbled
			slices.SortFunc(inputs, func(a, b string) int {
				return strings.Compare(a, b)
			})
			inputLen := len(inputs)
			// create workgroup
			wg := sizedwaitgroup.New(threadCount)

			splitFileName := strings.Split(inputs[0], ".")
			fileSuffix := splitFileName[len(splitFileName)-1]
			for i := 0; i < inputLen; i++ {
				wg.Add()
				go func(i int) {
					defer wg.Done()

					in := inputs[i]
					out := fmt.Sprintf("frame%04d.%s", i, fileSuffix)
					println(fmt.Sprintf("Sorting image %d (%q -> %q)...", i+1, in, out))
					err := sortingTime(in, out, shared.Config.Mask)
					if err != nil {
						cli.Exit(fmt.Sprintf("Error occured during sort of image %d (%q): %q", i+1, in, err), 1)
					}
				}(i)
			}
			wg.Wait()
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func sortingTime(input, output, maskpath string) error {
	file, err := os.Open(input)
	if err != nil {
		cli.Exit("Input image could not be opened", 1)
	}
	defer file.Close()

	rawImg, _, err := image.Decode(file)
	if err != nil {
		cli.Exit("Input image could not be decoded", 1)
	}

	// convert to rgba
	b := rawImg.Bounds()
	imgDims := image.Rect(0, 0, b.Dx(), b.Dy())
	img := image.NewRGBA(imgDims)
	mask := image.NewRGBA(imgDims)
	draw.Draw(img, img.Bounds(), rawImg, b.Min, draw.Src)
	rawImg = nil

	// TODO: figure out how to load mask once
	if maskpath != "" {
		maskFile, err := os.Open(maskpath)

		if err != nil {
			cli.Exit("Mask image could not be opened", 1)
		}
		defer maskFile.Close()
		rawMask, _, err := image.Decode(maskFile)
		if err != nil {
			cli.Exit("Mask image could not be decoded", 1)
		}
		draw.Draw(mask, mask.Bounds(), rawMask, b.Min, draw.Src)
		rawMask = nil
	}

	dims := img.Bounds().Max

	// split image into rows
	rows := make([][]types.PixelWithMask, dims.Y)
	for y := 0; y < dims.Y; y++ {
		row := make([]types.PixelWithMask, dims.X)

		for x := 0; x < dims.X; x++ {
			pixel := img.At(x, y).(color.RGBA)
			masked := mask.At(x, y).(color.RGBA).R
			wrapped := types.PixelWithMask{R: pixel.R, G: pixel.G, B: pixel.B, A: pixel.A, Mask: masked}
			row[x] = wrapped
		}
		rows[y] = row
	}

	/// more whitespace
	/// im not gonna rant again
	/// just
	/// *sigh*

	// pass the rows to the sorter
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		sorters.Sort(row)
	}

	/// cant believe i have to do this so the stupid extension doesnt fucking trim my lines
	/// like fuck dude i just want some fucking whitespace, its not that big of a deal

	// now write
	outputImg := image.NewRGBA(imgDims)
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		for j := 0; j < len(row); j++ {
			currPixWithMask := row[j]
			pixel := color.RGBA{currPixWithMask.R, currPixWithMask.G, currPixWithMask.B, currPixWithMask.A}
			outputImg.SetRGBA(j, i, pixel)
		}
	}

	f, err := os.Create(output)
	if err != nil {
		cli.Exit("Could not create output file", 1)
	}
	if strings.HasSuffix(input, "jpg") || strings.HasSuffix(input, "jpeg") {
		options := jpeg.Options{
			Quality: shared.Config.Quality,
		}
		jpeg.Encode(f, outputImg.SubImage(imgDims), &options)
	} else {
		png.Encode(f, outputImg)
	}
	return nil
}
