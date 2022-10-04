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
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/AksAman/mandelbrot/mandelbrot"
	"github.com/AksAman/mandelbrot/utils"
	"github.com/disintegration/imaging"
)

var (
	out        = flag.String("out", ".png", "Name of the output file with extension")
	iterations = flag.Int("iter", mandelbrot.DefaultConfig.MaxIterations, "Max Iterations")
	width      = flag.Int("width", mandelbrot.DefaultConfig.Width, "Width of the image")
	height     = flag.Int("height", mandelbrot.DefaultConfig.Height, "Height of the image")
	threshold  = flag.Float64("threshold", mandelbrot.DefaultConfig.Threshold, "Threshold for the mandelbrot set")
	workers    = flag.Int("workers", mandelbrot.DefaultConfig.Workers, "Number of workers to use")
	scale      = flag.Int("scale", 1, "Scale of the image")
	mode       = flag.String("mode", string(mandelbrot.DefaultConfig.Mode), "Mode of the image (options: seq, pixel, row, workers)")
	port       = flag.String("port", "8080", "Port to run the server on")
)

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", r.URL.Path)
	}))
	mux.HandleFunc("/mandelbrot", mandelbrotHandler)

	addr := fmt.Sprintf(":%s", *port)
	log.Printf("Server running on port %s", addr)
	http.ListenAndServe(addr, loggerMiddleware(mux))
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			if r.URL.Path != "/favicon.ico" {
				log.Printf("%s %s %s %s", r.Method, time.Since(start), r.RequestURI, r.Proto)
			}
		},
	)
}

func mandelbrotHandler(w http.ResponseWriter, r *http.Request) {

	width := utils.GetQueryParam(r, "width", *width)
	height := utils.GetQueryParam(r, "height", *height)
	iterations := utils.GetQueryParam(r, "iterations", *iterations)
	threshold := utils.GetQueryParam(r, "threshold", *threshold)
	workers := utils.GetQueryParam(r, "workers", *workers)
	scale := utils.GetQueryParam(r, "scale", *scale)
	mode := utils.GetQueryParam(r, "mode", *mode)
	out := utils.GetQueryParam(r, "out", *out)
	zoom := utils.GetQueryParam(r, "zoom", 1.)
	smooth := utils.GetQueryParam(r, "smooth", false)
	offsetX := utils.GetQueryParam(r, "offsetX", 0.)
	offsetY := utils.GetQueryParam(r, "offsetY", 0.)
	hue := utils.GetQueryParam(r, "hue", 0.)
	save := utils.GetQueryParam(r, "save", false)

	config := mandelbrot.Config{
		Width:         width,
		Height:        height,
		Threshold:     threshold,
		Workers:       workers,
		Scale:         scale,
		Mode:          mandelbrot.Mode(mode),
		MaxIterations: iterations,
		Zoom:          zoom,
		Smooth:        smooth,
		OffsetX:       offsetX,
		OffsetY:       offsetY,
		HueOffset:     hue,
	}

	img, err := mandelbrot.Create(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	img = imaging.AdjustContrast(img, 20)
	img = imaging.AdjustBrightness(img, 20)

	err = EncodeImage(w, img, out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if save {
		filename := GetFilenameWithFlags("./img/web-colored.jpg", config)
		err = SaveImage(img, filename)
		if err == nil {
			log.Println("Saved image to", filename)
		}

	}
}

func EncodeImage(w http.ResponseWriter, img image.Image, extension string) (err error) {
	switch extension {
	case ".png":
		err = png.Encode(w, img)
	case ".jpg", ".jpeg":
		err = jpeg.Encode(w, img, &jpeg.Options{Quality: 100})
	case ".gif":
		err = gif.Encode(w, img, &gif.Options{NumColors: 256})
	default:
		err = errors.New("Unsupported image format: " + extension)
	}

	return err
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
		err = jpeg.Encode(f, img, &jpeg.Options{Quality: 75})
	case ".gif":
		err = gif.Encode(f, img, &gif.Options{NumColors: 256})
	default:
		err = errors.New("Unsupported image format: " + ext)
	}

	return err
}
