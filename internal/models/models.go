package models

import (
	"image"
	"time"
)

// ImageInfo contains metadata about an image
type ImageInfo struct {
	Path       string            `json:"path"`
	Format     string            `json:"format"`
	Width      int               `json:"width"`
	Height     int               `json:"height"`
	ColorModel string            `json:"color_model"`
	BitDepth   int               `json:"bit_depth,omitempty"`
	HasAlpha   bool              `json:"has_alpha"`
	FileSize   int64             `json:"file_size"`
	ModifiedAt time.Time         `json:"modified_at,omitempty"`
	EXIF       map[string]string `json:"exif,omitempty"`
	DPI        [2]float64        `json:"dpi,omitempty"`
	Animated   bool              `json:"animated,omitempty"`
	FrameCount int               `json:"frame_count,omitempty"`
}

// Color represents an RGB color
type Color struct {
	R     uint8  `json:"r"`
	G     uint8  `json:"g"`
	B     uint8  `json:"b"`
	Count int    `json:"count,omitempty"`
	Hex   string `json:"hex"`
}

// ColorPalette represents a set of dominant colors
type ColorPalette struct {
	Colors []Color `json:"colors"`
}

// ResizeOptions configures image resizing
type ResizeOptions struct {
	Width      int
	Height     int
	Filter     string // nearest, bilinear, bicubic, lanczos
	KeepAspect bool
}

// CropOptions configures image cropping
type CropOptions struct {
	X      int
	Y      int
	Width  int
	Height int
}

// RotateOptions configures image rotation
type RotateOptions struct {
	Angle     float64
	FillColor string // hex color for background
}

// OptimizeOptions configures image optimization
type OptimizeOptions struct {
	Quality   int  // 1-100 for JPEG
	Lossless  bool // for PNG/WebP
	StripMeta bool
}

// WatermarkOptions configures watermarking
type WatermarkOptions struct {
	Text      string
	ImagePath string
	Position  string // top-left, top-right, bottom-left, bottom-right, center
	Opacity   float64
	FontSize  int
}

// CompareResult contains image comparison results
type CompareResult struct {
	Identical     bool    `json:"identical"`
	SSIM          float64 `json:"ssim"`
	PSNR          float64 `json:"psnr"`
	PixelDiff     int     `json:"pixel_diff"`
	TotalPixels   int     `json:"total_pixels"`
	DiffPercent   float64 `json:"diff_percent"`
	DiffImagePath string  `json:"diff_image_path,omitempty"`
}

// BatchOperation defines a batch processing operation
type BatchOperation struct {
	Type      string // resize, convert, optimize, strip, watermark
	Options   interface{}
	InputDir  string
	OutputDir string
	Recursive bool
	Pattern   string // glob pattern for filtering
}

// Format constants
const (
	FormatJPEG = "jpeg"
	FormatPNG  = "png"
	FormatGIF  = "gif"
	FormatBMP  = "bmp"
	FormatTIFF = "tiff"
	FormatWEBP = "webp"
)

// Filter constants
const (
	FilterNearest  = "nearest"
	FilterBilinear = "bilinear"
	FilterBicubic  = "bicubic"
	FilterLanczos  = "lanczos"
)

// Position constants
const (
	PositionTopLeft     = "top-left"
	PositionTopRight    = "top-right"
	PositionBottomLeft  = "bottom-left"
	PositionBottomRight = "bottom-right"
	PositionCenter      = "center"
)

// Helper to convert image.Image dimensions
func GetImageBounds(img image.Image) (int, int) {
	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy()
}
