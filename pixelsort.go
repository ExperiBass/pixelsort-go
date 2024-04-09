package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"pixelsort_go/intervals"
	"pixelsort_go/patterns"
	"pixelsort_go/shared"
	"slices"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/remeh/sizedwaitgroup"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:                   "pixelsort_go",
		Usage:                  "Organize pixels.",
		Version:                "v1.0.0",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "`image`(s) to sort, or a dir full of images",
				Required: true,
				Action: func(ctx *cli.Context, v []string) error {
					if len(v) < 1 {
						return cli.Exit("No inputs given", 1)
					}
					/// ... is this needed?
					// image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
					// image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "pattern",
				Aliases: []string{"p"},
				Usage:   "`pattern` loader to use",
				Value:   "row",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "`file` to output to",
			},
			&cli.StringFlag{
				Name:    "mask",
				Aliases: []string{"m"},
				Usage:   "b&w `mask` to determine which pixels to touch; white is skipped",
			},
			&cli.StringFlag{
				Name:    "interval",
				Value:   "row",
				Aliases: []string{"I"},
				// TODO: print valid intervals and comparators
				Usage: "interval `func`tion to use",
			},
			&cli.StringFlag{
				Name:    "comparator",
				Value:   "lightness",
				Aliases: []string{"c"},
				Usage:   "comparison `func`tion to use",
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
			&cli.Float64Flag{
				Name:    "angle",
				Value:   0.0,
				Aliases: []string{"a"},
				Usage:   "rotate the image by `deg`rees, pos or neg",
				// TODO: clamp to -360 - 360
			},
			&cli.IntFlag{
				Name:    "section_length",
				Value:   60,
				Aliases: []string{"L"},
				Usage:   "The base `len`gth of each slice",
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
		Commands: []*cli.Command{},
		Action: func(ctx *cli.Context) error {
			inputs := ctx.StringSlice("input")
			output := ctx.String("output")
			mask := ctx.String("mask")
			masks := make([]string, 0)
			shared.Config.Pattern = ctx.String("pattern")
			shared.Config.Interval = ctx.String("interval")
			shared.Config.Comparator = ctx.String("comparator")
			shared.Config.Thresholds.Lower = float32(ctx.Float64("lower_threshold"))
			shared.Config.Thresholds.Upper = float32(ctx.Float64("upper_threshold"))
			shared.Config.SectionLength = ctx.Int("section_length")
			shared.Config.Reverse = ctx.Bool("reverse")
			shared.Config.Randomness = float32(ctx.Float64("randomness"))
			shared.Config.Angle = ctx.Float64("angle")
			threadCount := ctx.Int("threads")

			// profiling
			/*masked := "unmasked"
			if mask != "" || len(masks) > 0 {
				masked = "masked"
			}
			profileFile, err := os.Create(fmt.Sprintf("cpuprofile-%s-%s-%s.prof", shared.Config.Comparator, shared.Config.Sorter, masked))
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(profileFile)
			defer pprof.StopCPUProfile()*/

			/// this can be done better but im lazy and braindead
			/// MAYBE: accept multiple dirs? pop them and append contents?
			inputLen := len(inputs)
			if inputLen == 1 {
				input := inputs[0]

				inputfile, err := os.Open(input)
				if err != nil {
					return cli.Exit("Input could not be opened", 1)
				}
				defer inputfile.Close()

				inputStat, err := inputfile.Stat()
				if err != nil {
					return cli.Exit(fmt.Sprintf("Error getting input file stats: %s", err), 1)
				}
				inputfile = nil

				if inputStat.IsDir() {
					res, err := readDirImages(input)
					if err != nil {
						return err
					}
					inputs = res
					inputLen = len(inputs)
				}
			}

			// masking
			if mask != "" {
				maskfile, err := os.Open(mask)
				if err != nil {
					return cli.Exit("Mask could not be opened", 1)
				}
				defer maskfile.Close()
				maskStat, err := maskfile.Stat()
				if err != nil {
					return cli.Exit(fmt.Sprintf("Error getting mask stats: %s", err), 1)
				}
				maskfile = nil
				if maskStat.IsDir() {
					res, err := readDirImages(mask)
					if err != nil {
						return err
					}
					masks = res
				} else {
					masks = append(masks, mask)
				}
			}

			maskLen := len(masks)
			if maskLen == 0 {
				// empty string, will be ignored by sorting
				masks = append(masks, "")
				maskLen = len(masks)
			}
			infoString := fmt.Sprintf("Sorting %d images with a config of %+v.", inputLen, shared.Config)
			println(infoString)

			// multiple imgs
			// sort em first so frames dont get jumbled
			slices.SortFunc(inputs, func(a, b string) int {
				return strings.Compare(a, b)
			})
			slices.SortFunc(masks, func(a, b string) int {
				return strings.Compare(a, b)
			})

			// create workgroup
			wg := sizedwaitgroup.New(threadCount)

			for i := 0; i < inputLen; i++ {
				wg.Add()
				go func(i int) {
					defer wg.Done()

					in := inputs[i]
					splitFileName := strings.Split(inputs[0], ".")
					fileSuffix := splitFileName[len(splitFileName)-1] // MAYBE: just use .png and .jpg
					out := output
					if inputLen > 1 {
						out = fmt.Sprintf("frame%04d.%s", i, fileSuffix)
					} else if out == "" {
						out = fmt.Sprintf("sorted.%s", fileSuffix)
					}
					println(fmt.Sprintf("Loading image %d (%q -> %q)...", i+1, in, out))
					maskIdx := min(i, maskLen-1)
					err := sortingTime(in, out, masks[maskIdx])
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

func readDirImages(input string) ([]string, error) {
	files, err := os.ReadDir(input)
	if err != nil {
		return nil, cli.Exit("Couldn't read input dir", 1)
	}
	inputs := make([]string, len(files)) // allocate enough space
	for idx, file := range files {
		if !file.IsDir() && file.Type().IsRegular() {
			name := file.Name()
			if strings.HasSuffix(name, "jpg") || strings.HasSuffix(name, "jpeg") || strings.HasSuffix(name, "png") {
				inputs[idx] = fmt.Sprintf("%s/%s", input, name)
			}
		}
	}
	// remove empty elms
	inputs = slices.DeleteFunc(inputs, func(s string) bool {
		return len(strings.TrimSpace(s)) == 0
	})
	return inputs, nil
}

func sortingTime(input, output, maskpath string) error {
	file, err := os.Open(input)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Input \"%s\" could not be opened", input), 1)
	}
	defer file.Close()

	rawImg, _, err := image.Decode(file)
	if err != nil {
		return cli.Exit(fmt.Sprintf("Input \"%s\" could not be decoded", input), 1)
	}

	/// RO TA TE
	/// god why do i have to do thissssssswddenwfiosbduglzx er agdxbv
	/// this is used in the writing step cause `imaging` doesnt have a option to
	/// auto-crop transparency
	originalDims := rawImg.Bounds()
	if math.Mod(shared.Config.Angle, 360) != 0 {
		rawImg = (*image.RGBA)(imaging.Rotate(rawImg, shared.Config.Angle, color.Transparent))
	}

	/// convert to rgba
	b := rawImg.Bounds()
	sortingDims := image.Rect(0, 0, b.Dx(), b.Dy())
	img := image.NewRGBA(sortingDims)
	mask := image.NewRGBA(sortingDims)
	draw.Draw(img, img.Bounds(), rawImg, b.Min, draw.Src)
	rawImg = nil

	/// MAYBE: figure out how to load mask once if used multiple times
	if maskpath != "" {
		maskFile, err := os.Open(maskpath)
		if err != nil {
			return cli.Exit(fmt.Sprintf("Mask \"%s\" could not be opened", maskpath), 1)
		}
		defer maskFile.Close()

		rawMask, _, err := image.Decode(maskFile)
		if err != nil {
			return cli.Exit(fmt.Sprintf("Mask \"%s\" could not be decoded", maskpath), 1)
		}

		/// RO TA TE (again)
		if math.Mod(shared.Config.Angle, 180) != 0 {
			rawMask = imaging.Rotate(rawMask, float64(shared.Config.Angle), color.Transparent)
		}

		draw.Draw(mask, mask.Bounds(), rawMask, b.Min, draw.Src)
		rawMask = nil
	}

	/// load stretches
	stretches := patterns.Loader[fmt.Sprintf("%sload", shared.Config.Pattern)](*img, *mask)

	/// more whitespace
	/// im not gonna rant again
	/// just
	/// *sigh*

	println(fmt.Sprintf("Sorting %s...", input))
	/// pass the rows to the sorter
	start := time.Now()
	for i := 0; i < len(stretches); i++ {
		row := stretches[i]
		intervals.Sort(row)
	}
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println(output, "elapsed:", elapsed.Truncate(time.Millisecond).String())

	/// cant believe i have to do this so the stupid extension doesnt fucking trim my lines
	/// like fuck dude i just want some fucking whitespace, its not that big of a deal

	/// now write
	outputImg := patterns.Saver[fmt.Sprintf("%ssave", shared.Config.Pattern)](stretches, img.Bounds())

	/// ET AT OR
	if math.Mod(shared.Config.Angle, 360) != 0 {
		outputImg = (*image.RGBA)(imaging.Rotate(outputImg, -shared.Config.Angle, color.Transparent))
		/// gotta crop the invisible pixels
		if math.Mod(shared.Config.Angle, 90) != 0 {
			outputImg = (*image.RGBA)(imaging.CropCenter(outputImg, originalDims.Dx(), originalDims.Dy()))
		}
	}

	println(fmt.Sprintf("Writing %s...", output))
	f, err := os.Create(output)
	if err != nil {
		return cli.Exit("Could not create output file", 1)
	}

	/// spit the result out
	/// MAYBE: arbitrary extensions?
	if strings.HasSuffix(input, "jpg") || strings.HasSuffix(input, "jpeg") {
		options := jpeg.Options{
			Quality: 100,
		}
		jpeg.Encode(f, outputImg.SubImage(outputImg.Rect), &options)
	} else {
		png.Encode(f, outputImg)
	}
	return nil
}
