package mandelbrot

import (
	"errors"
	"image"
	"image/color"
	"log"
	"sync"
)

// Mandelbrot implements image.Image interface
type Mandelbrot struct {
	Config Config
	img    [][]color.RGBA
}

// interface methods

// ColorModel returns the Image's color model.
func (mandel *Mandelbrot) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (mandel *Mandelbrot) Bounds() image.Rectangle {
	return image.Rect(0, 0, mandel.Config.Width, mandel.Config.Height)
}

// At returns the color of the pixel at (x, y).
func (mandel *Mandelbrot) At(x, y int) color.Color {
	return mandel.img[x][y]
}

func initMandelbrot(config ...Config) *Mandelbrot {
	cfg := configDefault(config...)

	mandel := Mandelbrot{
		Config: cfg,
		img:    make([][]color.RGBA, cfg.Width),
	}

	for i := range mandel.img {
		mandel.img[i] = make([]color.RGBA, cfg.Height)
	}

	return &mandel
}

func Create(config ...Config) (image.Image, error) {
	mandel := initMandelbrot(config...)

	log.Printf("Using mode: %v\n", mandel.Config.Mode)
	switch mandel.Config.Mode {
	case Sequential:
		mandel.sequentialFill()
	case Pixel:
		mandel.fillUsingOneGoroutinePerPixel()
	case Row:
		mandel.fillUsingOneGoroutinePerRow()
	case Parallel:
		mandel.fillUsingWorkers()
	default:
		return nil, errors.New("invalid mode")
	}

	return mandel, nil
}

// --scale 1 --threshold 32 --iter 1000  0.57s user 0.17s system 112% cpu 0.660 total
// sequentialFill fills the image sequentially
func (mandel *Mandelbrot) sequentialFill() {
	for i, row := range mandel.img {
		for j := range row {
			mandel.fillPixel(i, j)
		}
	}
}

// --scale 1 --threshold 32 --iter 1000  1.12s user 0.27s system 247% cpu 0.564 total
// fillUsingOneGoroutinePerPixel one goroutine per pixel
func (mandel *Mandelbrot) fillUsingOneGoroutinePerPixel() {
	wg := &sync.WaitGroup{}
	wg.Add(mandel.Config.Width * mandel.Config.Height)
	for i, row := range mandel.img {
		for j := range row {
			go func(i, j int) {
				defer wg.Done()
				mandel.fillPixel(i, j)
			}(i, j)
		}
	}
	wg.Wait()
}

// --scale 1 --threshold 32 --iter 1000  0.76s user 0.15s system 235% cpu 0.384 total
// fillUsingOneGoroutinePerRow creates one goroutine for every row
func (mandel *Mandelbrot) fillUsingOneGoroutinePerRow() {
	wg := &sync.WaitGroup{}
	wg.Add(mandel.Config.Width)
	for i := range mandel.img {
		go func(i int) {
			defer wg.Done()
			for j := range mandel.img[i] {
				mandel.fillPixel(i, j)
			}
		}(i)
	}
	wg.Wait()
}

// --scale 1 --threshold 32 --iter 1000  1.30s user 0.23s system 179% cpu 0.856 total
// fillUsingWorkers uses fixed user defined count of goroutines to fill image
func (mandel *Mandelbrot) fillUsingWorkers() {
	workers := mandel.Config.Workers

	log.Printf("using %v workers\n", workers)

	type workerJob struct {
		i, j int
	}
	workerChan := make(chan workerJob)

	wg := &sync.WaitGroup{}
	wg.Add(workers)

	// create workers with jobs waiting on channel
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for job := range workerChan {
				mandel.fillPixel(job.i, job.j)
			}
		}()
	}

	// create jobs and send on channel
	for i, row := range mandel.img {
		for j := range row {
			mandel.fillPixel(i, j)
			workerChan <- workerJob{i, j}
		}
	}

	close(workerChan)
	wg.Wait()
}

func (mandel *Mandelbrot) fillPixel(px, py int) {
	// normalized x coord

	x0 := mapRange(float64(px), 0, float64(mandel.Config.Width), mandel.Config.XScale.min, mandel.Config.XScale.max)
	y0 := mapRange(float64(py), 0, float64(mandel.Config.Height), mandel.Config.YScale.min, mandel.Config.YScale.max)

	x, y := 0., 0.
	x2, y2 := 0., 0.
	for i := 0; x2+y2 <= mandel.Config.Threshold && i < mandel.Config.MaxIterations; i++ {
		y = (x+x)*y + y0
		x = x2 - y2 + x0

		x2 = x * x
		y2 = y * y
	}

	addColor(&mandel.img[px][py], x, y)

}

func addColor(c *color.RGBA, x, y float64) {
	n := byte(x * y)

	// TODO: Add user defined palletes

	c.R, c.G, c.B, c.A = n, n, n, 255
}
