package processor

import (
	"image"
	"image/color"
	"testing"

	"github.com/EdgarOrtegaRamirez/pixelforge/internal/models"
)

func TestResize(t *testing.T) {
	img := createTestImage()

	tests := []struct {
		name    string
		opts    models.ResizeOptions
		wantW   int
		wantH   int
		wantErr bool
	}{
		{
			name:  "resize by width only",
			opts:  models.ResizeOptions{Width: 50, Height: 0, KeepAspect: true},
			wantW: 50,
			wantH: 50,
		},
		{
			name:  "resize by height only",
			opts:  models.ResizeOptions{Width: 0, Height: 50, KeepAspect: true},
			wantW: 50,
			wantH: 50,
		},
		{
			name:  "resize both dimensions with aspect",
			opts:  models.ResizeOptions{Width: 80, Height: 80, KeepAspect: true},
			wantW: 80,
			wantH: 80,
		},
		{
			name:  "resize both dimensions without aspect",
			opts:  models.ResizeOptions{Width: 60, Height: 40, KeepAspect: false},
			wantW: 60,
			wantH: 40,
		},
		{
			name:  "resize with different filters",
			opts:  models.ResizeOptions{Width: 50, Height: 50, Filter: models.FilterBilinear, KeepAspect: true},
			wantW: 50,
			wantH: 50,
		},
		{
			name:    "error: no dimensions",
			opts:    models.ResizeOptions{Width: 0, Height: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Resize(img, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			bounds := result.Bounds()
			if bounds.Dx() != tt.wantW || bounds.Dy() != tt.wantH {
				t.Errorf("Resize() dimensions = %dx%d, want %dx%d", bounds.Dx(), bounds.Dy(), tt.wantW, tt.wantH)
			}
		})
	}
}

func TestCrop(t *testing.T) {
	img := createTestImage()

	tests := []struct {
		name    string
		opts    models.CropOptions
		wantW   int
		wantH   int
		wantErr bool
	}{
		{
			name:  "crop from origin",
			opts:  models.CropOptions{X: 0, Y: 0, Width: 50, Height: 50},
			wantW: 50,
			wantH: 50,
		},
		{
			name:  "crop with offset",
			opts:  models.CropOptions{X: 25, Y: 25, Width: 50, Height: 50},
			wantW: 50,
			wantH: 50,
		},
		{
			name:  "crop with clamping",
			opts:  models.CropOptions{X: 80, Y: 80, Width: 50, Height: 50},
			wantW: 20,
			wantH: 20,
		},
		{
			name:    "error: negative coordinates",
			opts:    models.CropOptions{X: -10, Y: 0, Width: 50, Height: 50},
			wantErr: true,
		},
		{
			name:    "error: zero dimensions",
			opts:    models.CropOptions{X: 0, Y: 0, Width: 0, Height: 50},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Crop(img, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Crop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			bounds := result.Bounds()
			if bounds.Dx() != tt.wantW || bounds.Dy() != tt.wantH {
				t.Errorf("Crop() dimensions = %dx%d, want %dx%d", bounds.Dx(), bounds.Dy(), tt.wantW, tt.wantH)
			}
		})
	}
}

func TestRotate(t *testing.T) {
	img := createTestImage()

	tests := []struct {
		name    string
		opts    models.RotateOptions
		wantW   int
		wantH   int
		wantErr bool
	}{
		{
			name:  "rotate 90 degrees",
			opts:  models.RotateOptions{Angle: 90},
			wantW: 100,
			wantH: 100,
		},
		{
			name:  "rotate 45 degrees",
			opts:  models.RotateOptions{Angle: 45, FillColor: "#ffffff"},
			wantW: 142, // approximately 100*sqrt(2)
			wantH: 142,
		},
		{
			name:    "rotate with invalid color",
			opts:    models.RotateOptions{Angle: 45, FillColor: "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Rotate(img, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rotate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			bounds := result.Bounds()
			// Allow some tolerance for rotated dimensions
			if bounds.Dx() < tt.wantW-5 || bounds.Dx() > tt.wantW+5 {
				t.Errorf("Rotate() width = %d, want ~%d", bounds.Dx(), tt.wantW)
			}
			if bounds.Dy() < tt.wantH-5 || bounds.Dy() > tt.wantH+5 {
				t.Errorf("Rotate() height = %d, want ~%d", bounds.Dy(), tt.wantH)
			}
		})
	}
}

func TestThumbnail(t *testing.T) {
	img := createTestImage()

	result, err := Thumbnail(img, 50)
	if err != nil {
		t.Fatalf("Thumbnail() error = %v", err)
	}

	bounds := result.Bounds()
	if bounds.Dx() != 50 || bounds.Dy() != 50 {
		t.Errorf("Thumbnail() dimensions = %dx%d, want 50x50", bounds.Dx(), bounds.Dy())
	}

	// Test error case
	_, err = Thumbnail(img, 0)
	if err == nil {
		t.Error("Thumbnail() expected error for size 0")
	}
}

func TestGrayscale(t *testing.T) {
	img := createTestImage()
	result := Grayscale(img)

	if result == nil {
		t.Fatal("Grayscale() returned nil")
	}

	bounds := result.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("Grayscale() dimensions = %dx%d, want 100x100", bounds.Dx(), bounds.Dy())
	}
}

func TestInvert(t *testing.T) {
	img := createTestImage()
	result := Invert(img)

	if result == nil {
		t.Fatal("Invert() returned nil")
	}

	// Check that colors are inverted
	bounds := result.Bounds()
	r, g, b, _ := result.At(bounds.Min.X, bounds.Min.Y).RGBA()
	_ = b // unused but part of the color decomposition

	// Original was red (255,0,0), inverted should be cyan (0,255,255)
	if r > 100 {
		t.Errorf("Invert() red channel too high: %d", r)
	}
	if g < 60000 {
		t.Errorf("Invert() green channel too low: %d", g)
	}
}

func TestFlip(t *testing.T) {
	img := createTestImage()

	// Horizontal flip
	resultH := Flip(img, true)
	if resultH == nil {
		t.Fatal("Flip(horizontal) returned nil")
	}

	// Vertical flip
	resultV := Flip(img, false)
	if resultV == nil {
		t.Fatal("Flip(vertical) returned nil")
	}
}

func TestAdjustBrightness(t *testing.T) {
	img := createTestImage()

	// Brighten
	result := AdjustBrightness(img, 0.5)
	if result == nil {
		t.Fatal("AdjustBrightness() returned nil")
	}

	// Darken
	result = AdjustBrightness(img, -0.5)
	if result == nil {
		t.Fatal("AdjustBrightness() returned nil")
	}
}

func TestAdjustContrast(t *testing.T) {
	img := createTestImage()

	result := AdjustContrast(img, 0.5)
	if result == nil {
		t.Fatal("AdjustContrast() returned nil")
	}
}

func TestAdjustSaturation(t *testing.T) {
	img := createTestImage()

	result := AdjustSaturation(img, 0.5)
	if result == nil {
		t.Fatal("AdjustSaturation() returned nil")
	}
}

func TestBlur(t *testing.T) {
	img := createTestImage()

	result := Blur(img, 2.0)
	if result == nil {
		t.Fatal("Blur() returned nil")
	}
}

func TestSharpen(t *testing.T) {
	img := createTestImage()

	result := Sharpen(img, 1.0)
	if result == nil {
		t.Fatal("Sharpen() returned nil")
	}
}

func TestAdjustGamma(t *testing.T) {
	img := createTestImage()

	result := AdjustGamma(img, 2.2)
	if result == nil {
		t.Fatal("AdjustGamma() returned nil")
	}

	// Test invalid gamma (should use default)
	result = AdjustGamma(img, -1.0)
	if result == nil {
		t.Fatal("AdjustGamma() returned nil for negative gamma")
	}
}

func TestRotate90(t *testing.T) {
	img := createTestImage()
	result := Rotate90(img)

	if result == nil {
		t.Fatal("Rotate90() returned nil")
	}

	bounds := result.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("Rotate90() dimensions = %dx%d, want 100x100", bounds.Dx(), bounds.Dy())
	}
}

func TestRotate180(t *testing.T) {
	img := createTestImage()
	result := Rotate180(img)

	if result == nil {
		t.Fatal("Rotate180() returned nil")
	}
}

func TestRotate270(t *testing.T) {
	img := createTestImage()
	result := Rotate270(img)

	if result == nil {
		t.Fatal("Rotate270() returned nil")
	}
}

func TestPaste(t *testing.T) {
	bg := createTestImage()
	fg := image.NewRGBA(image.Rect(0, 0, 50, 50))

	// Fill foreground with blue
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			fg.Set(x, y, color.RGBA{0, 0, 255, 255})
		}
	}

	result := Paste(bg, fg, 25, 25)
	if result == nil {
		t.Fatal("Paste() returned nil")
	}
}

func TestOverlay(t *testing.T) {
	bg := createTestImage()
	overlayImg := image.NewRGBA(image.Rect(0, 0, 50, 50))

	// Fill overlay with semi-transparent red
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			overlayImg.Set(x, y, color.RGBA{255, 0, 0, 128})
		}
	}

	result := Overlay(bg, overlayImg, 25, 25, 0.5)
	if result == nil {
		t.Fatal("Overlay() returned nil")
	}
}

func TestParseHexColor(t *testing.T) {
	tests := []struct {
		hex     string
		wantR   uint8
		wantG   uint8
		wantB   uint8
		wantA   uint8
		wantErr bool
	}{
		{"#ff0000", 255, 0, 0, 255, false},
		{"#00ff00", 0, 255, 0, 255, false},
		{"#0000ff", 0, 0, 255, 255, false},
		{"#ffffff", 255, 255, 255, 255, false},
		{"#000000", 0, 0, 0, 255, false},
		{"#f00", 255, 0, 0, 255, false},      // Short form
		{"#ff000080", 255, 0, 0, 128, false}, // With alpha
		{"", 0, 0, 0, 0, true},               // Empty
		{"invalid", 0, 0, 0, 0, true},        // Invalid
		{"#gg0000", 0, 0, 0, 0, true},        // Invalid hex
	}

	for _, tt := range tests {
		t.Run(tt.hex, func(t *testing.T) {
			result, err := parseHexColor(tt.hex)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHexColor(%q) error = %v, wantErr %v", tt.hex, err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			if result.R != tt.wantR || result.G != tt.wantG || result.B != tt.wantB || result.A != tt.wantA {
				t.Errorf("parseHexColor(%q) = rgba(%d,%d,%d,%d), want rgba(%d,%d,%d,%d)",
					tt.hex, result.R, result.G, result.B, result.A, tt.wantR, tt.wantG, tt.wantB, tt.wantA)
			}
		})
	}
}

func TestGetFilter(t *testing.T) {
	tests := []struct {
		name string
	}{
		{models.FilterNearest},
		{models.FilterBilinear},
		{models.FilterBicubic},
		{models.FilterLanczos},
		{""}, // default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just ensure it doesn't panic
			_ = getFilter(tt.name)
		})
	}
}

func createTestImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			var c color.Color
			if x < 50 && y < 50 {
				c = color.RGBA{255, 0, 0, 255}
			} else if x >= 50 && y < 50 {
				c = color.RGBA{0, 255, 0, 255}
			} else if x < 50 && y >= 50 {
				c = color.RGBA{0, 0, 255, 255}
			} else {
				c = color.RGBA{255, 255, 0, 255}
			}
			img.Set(x, y, c)
		}
	}

	return img
}
