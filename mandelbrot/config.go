package mandelbrot

type SetScale struct {
	min float64
	max float64
}

var DefaultXScale = SetScale{
	min: -2.5,
	max: 1,
}

var DefaultYScale = SetScale{
	min: -1,
	max: 1,
}

type Mode string

const (
	Sequential Mode = "seq"
	Pixel      Mode = "pixel"
	Row        Mode = "row"
	Parallel   Mode = "workers"
)

type Config struct {
	Width, Height int
	Threshold     float64
	MaxIterations int
	XScale        *SetScale
	YScale        *SetScale
	Workers       int
	Scale         int
	Mode          Mode
}

var DefaultConfig = Config{
	Width:         700,
	Height:        400,
	Threshold:     4.0,
	MaxIterations: 1000,
	Workers:       4,
	XScale:        &DefaultXScale,
	YScale:        &DefaultYScale,
	Mode:          Sequential,
}

func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return DefaultConfig
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if cfg.Width == 0 {
		cfg.Width = DefaultConfig.Width
	}

	if cfg.Height == 0 {
		cfg.Height = DefaultConfig.Height
	}

	if cfg.Threshold == 0 {
		cfg.Threshold = DefaultConfig.Threshold
	}

	if cfg.MaxIterations == 0 {
		cfg.MaxIterations = DefaultConfig.MaxIterations
	}
	if cfg.Workers == 0 {
		cfg.Workers = DefaultConfig.Workers
	}

	if cfg.XScale == nil {
		cfg.XScale = DefaultConfig.XScale
	}

	if cfg.YScale == nil {
		cfg.YScale = DefaultConfig.YScale
	}

	if cfg.Mode == "" {
		cfg.Mode = DefaultConfig.Mode
	}

	if cfg.Scale == 0 {
		cfg.Scale = 1
	}

	cfg.Width *= cfg.Scale
	cfg.Height *= cfg.Scale
	return cfg
}
