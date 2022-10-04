package main

import (
	"errors"
	"flag"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/AksAman/mandelbrot/mandelbrot"
)

var (
	out        = flag.String("out", "mandelbrot.png", "Name of the output file with extension")
	iterations = flag.Int("iter", mandelbrot.DefaultConfig.MaxIterations, "Max Iterations")
	width      = flag.Int("width", mandelbrot.DefaultConfig.Width, "Width of the image")
	height     = flag.Int("height", mandelbrot.DefaultConfig.Height, "Height of the image")
	threshold  = flag.Float64("threshold", mandelbrot.DefaultConfig.Threshold, "Threshold for the mandelbrot set")
	workers    = flag.Int("workers", mandelbrot.DefaultConfig.Workers, "Number of workers to use")
	scale      = flag.Int("scale", mandelbrot.DefaultConfig.Scale, "Scale of the image")
	mode       = flag.String("mode", string(mandelbrot.DefaultConfig.Mode), "Mode of the image (options: seq, pixel, row, workers)")
)

func main() {
	flag.Parse()

	tStart := time.Now()
	config := mandelbrot.Config{
		Width:         *width,
		Height:        *height,
		Threshold:     *threshold,
		Workers:       *workers,
		Scale:         *scale,
		Mode:          mandelbrot.Mode(*mode),
		MaxIterations: *iterations,
	}

	img, err := mandelbrot.Create(config)
	if err != nil {
		log.Fatal(err)
	}
	tTaken := time.Since(tStart)

	log.Printf("Time taken to create image: %s\n", tTaken)

	tStart = time.Now()

	err = SaveImage(img, *out)
	if err != nil {
		log.Fatal(err)
	}
	tTaken = time.Since(tStart)
	log.Printf("Time taken to save image: %s\n", tTaken)
}

func SaveImage(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	switch ext := filepath.Ext(filename); ext {
	case ".png":
		err = png.Encode(f, img)
	case ".jpg", ".jpeg":
		err = jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
	case ".gif":
		err = gif.Encode(f, img, &gif.Options{NumColors: 256})
	default:
		err = errors.New("Unsupported image format: " + ext)
	}

	return err
}
