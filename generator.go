package main

import (
	"errors"
	"fmt"
	"github.com/ojrac/opensimplex-go"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	width, height, file, err := parseArgs(os.Args[1:])
	if err != nil {
		showHelp(os.Args[0])
		os.Exit(1)
	}

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = generateImage(width, height, f)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func generateImage(width int, height int, target io.Writer) error {
	rand.Seed(time.Now().UnixMilli())

	// smaller = smoother, larger = rougher
	roughness := (1 + rand.Float64()*2) / float64(width)

	// initialize all noises here
	rNoise := opensimplex.New(rand.Int63())
	gNoise := opensimplex.New(rand.Int63())
	bNoise := opensimplex.New(rand.Int63())
	aNoise := opensimplex.New(rand.Int63())

	// fill the image with data from the noises
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.NRGBA{
				R: uint8(math.Abs(rNoise.Eval2(float64(x)*roughness, float64(y)*roughness) * 255)),
				G: uint8(math.Abs(gNoise.Eval2(float64(x)*roughness, float64(y)*roughness) * 255)),
				B: uint8(math.Abs(bNoise.Eval2(float64(x)*roughness, float64(y)*roughness) * 255)),
				// we limit alpha to only 3 bits of randomness to keep it above a certain level
				A: uint8(math.Abs(aNoise.Eval2(float64(x)*roughness, float64(y)*roughness)*64)) + 191,
			})
		}
	}

	return png.Encode(target, img)
}

func parseArgs(args []string) (int, int, string, error) {
	var err error = nil
	var width, height int64
	var file string

	if (len(args) < 1) || (len(args) > 3) {
		return 0, 0, "", errors.New("wrong arg count")
	}

	width, err = strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return 0, 0, "", err
	} else if width <= 0 {
		return 0, 0, "", errors.New("width must be > 0")
	}

	// we default to a square resolution if nothing else was given
	height = width
	if len(args) > 1 {
		height, err = strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return 0, 0, "", err
		} else if height <= 0 {
			return 0, 0, "", errors.New("height must be > 0")
		}
	}

	// we default to a file in the current directory if no path was given
	file = fmt.Sprintf("./randomImage-%d-%d.png", width, height)
	if len(args) == 3 {
		file = args[2]
	}

	return int(width), int(height), file, nil
}

func showHelp(programName string) {
	fmt.Println("Invalid arguments...")
	fmt.Printf(" %s <width> <height> [outputFile]\n", programName)
}
