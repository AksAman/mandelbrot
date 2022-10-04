package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/AksAman/mandelbrot/mandelbrot"
	"github.com/AksAman/mandelbrot/utils"
	"github.com/disintegration/imaging"
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
	zoom       = flag.Float64("zoom", mandelbrot.DefaultConfig.Zoom, "Zoom of the image")
	hueOffset  = flag.Float64("hue", mandelbrot.DefaultConfig.HueOffset, "Hue offset of the image")
	offsetX    = flag.Float64("offsetX", mandelbrot.DefaultConfig.OffsetX, "Offset X of the image")
	offsetY    = flag.Float64("offsetY", mandelbrot.DefaultConfig.OffsetY, "Offset Y of the image")
	jpgQuality = flag.Int("quality", 100, "JPG Quality")
)

func main() {
	flag.Parse()

	// loop()
	tStart := time.Now()
	config := mandelbrot.Config{
		Width:         *width,
		Height:        *height,
		Threshold:     *threshold,
		Workers:       *workers,
		Scale:         *scale,
		Mode:          mandelbrot.Mode(*mode),
		MaxIterations: *iterations,
		Zoom:          *zoom,
		HueOffset:     *hueOffset,
		OffsetX:       *offsetX,
		OffsetY:       *offsetY,
		Smooth:        true,
	}

	img, err := mandelbrot.Create(config)
	if err != nil {
		log.Fatal(err)
	}
	finalConfig := img.(*mandelbrot.Mandelbrot).Config
	img = imaging.AdjustContrast(img, 2)

	tTaken := time.Since(tStart)

	log.Printf("Time taken to create image: %s\n", tTaken)

	tStart = time.Now()

	err = SaveImage(img, GetFilenameWithFlags(*out, finalConfig))
	if err != nil {
		log.Fatal(err)
	}
	tTaken = time.Since(tStart)
	log.Printf("Time taken to save image: %s\n", tTaken)
}

func loop() {
	config := mandelbrot.Config{
		Width:         *width,
		Height:        *height,
		Threshold:     *threshold,
		Workers:       *workers,
		Scale:         *scale,
		Mode:          mandelbrot.Mode(*mode),
		MaxIterations: *iterations,
		Zoom:          *zoom,
		HueOffset:     *hueOffset,
		OffsetX:       0,
		OffsetY:       0,
		Smooth:        true,
	}
	ext := filepath.Ext(*out)
	fnameWithoutExt := strings.Split(*out, ext)[0]

	workers := 8
	type workerJob struct {
		i, j   int
		config mandelbrot.Config
	}
	jobChan := make(chan workerJob)

	wg := &sync.WaitGroup{}
	wg.Add(workers)

	// create workers with jobs waiting on channel
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for job := range jobChan {
				img, err := mandelbrot.Create(job.config)
				if err != nil {
					log.Println(job.i, job.j, err)
					continue
				}

				sum := utils.SumOfAllPixels(img)
				if sum == 0 {
					log.Println(job.i, job.j, "sum is 0")
					continue
				}
				// log.Println(job.i, job.j, "sum is", sum)

				img = imaging.AdjustContrast(img, 2)

				fname := fmt.Sprintf("%s_%d_%d.%v", fnameWithoutExt, job.i, job.j, ext)
				fname = GetFilenameWithFlags(fname, config)
				err = SaveImage(img, fname)
				if err != nil {
					log.Println(job.i, job.j, err)
					continue
				}
			}
		}()
	}

	for i := -*width; i < *width; i++ {
		for j := -*height; j < *height; j++ {
			config.OffsetX = float64(i) / float64(*width)
			config.OffsetY = float64(j) / float64(*height)
			jobChan <- workerJob{i, j, config}
		}
	}
	close(jobChan)
	wg.Wait()
}

func GetFilenameWithFlags(filename string, config mandelbrot.Config) string {
	ext := filepath.Ext(filename)
	filenameWithoutExt := strings.Split(filename, ext)[0]
	return filenameWithoutExt + "#" + getFlags(config) + ext
}

func getFlags(config mandelbrot.Config) string {
	return strings.Join([]string{
		"i=" + strconv.Itoa(config.MaxIterations),
		"t=" + fmt.Sprintf("%v", config.Threshold),
		"z=" + fmt.Sprintf("%v", config.Zoom),
		// "h=" + fmt.Sprintf("%v", *hueOffset),
		"x=" + fmt.Sprintf("%v", config.OffsetX),
		"y=" + fmt.Sprintf("%v", config.OffsetY),
	}, "_")
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
		err = jpeg.Encode(f, img, &jpeg.Options{Quality: *jpgQuality})
	case ".gif":
		err = gif.Encode(f, img, &gif.Options{NumColors: 256})
	default:
		err = errors.New("Unsupported image format: " + ext)
	}

	return err
}
