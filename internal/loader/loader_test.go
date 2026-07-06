package loader

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestGetFormat(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"image.jpg", "jpeg"},
		{"image.jpeg", "jpeg"},
		{"image.png", "png"},
		{"image.gif", "gif"},
		{"image.bmp", "bmp"},
		{"image.tiff", "tiff"},
		{"image.tif", "tiff"},
		{"image.webp", "webp"},
		{"image.txt", "unknown"},
		{"image", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := GetFormat(tt.path)
			if result != tt.expected {
				t.Errorf("GetFormat(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	// Create a test PNG file
	img := createTestImage()
	path := filepath.Join(t.TempDir(), "test.png")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}

	// Test loading
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded == nil {
		t.Fatal("Load() returned nil image")
	}

	bounds := loaded.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("Load() dimensions = %dx%d, want 100x100", bounds.Dx(), bounds.Dy())
	}
}

func TestLoadNonexistent(t *testing.T) {
	_, err := Load("/nonexistent/image.png")
	if err == nil {
		t.Error("Load() expected error for nonexistent file")
	}
}

func TestGetInfo(t *testing.T) {
	img := createTestImage()
	path := filepath.Join(t.TempDir(), "test.png")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}

	info, err := GetInfo(path)
	if err != nil {
		t.Fatalf("GetInfo() error = %v", err)
	}

	if info.Format != "png" {
		t.Errorf("GetInfo() Format = %q, want %q", info.Format, "png")
	}

	if info.Width != 100 || info.Height != 100 {
		t.Errorf("GetInfo() dimensions = %dx%d, want 100x100", info.Width, info.Height)
	}

	if info.ColorModel != "RGBA" {
		t.Errorf("GetInfo() ColorModel = %q, want %q", info.ColorModel, "RGBA")
	}
}

func TestGetInfoNonexistent(t *testing.T) {
	_, err := GetInfo("/nonexistent/image.png")
	if err == nil {
		t.Error("GetInfo() expected error for nonexistent file")
	}
}

func createTestImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// Fill with colored quadrants
	colors := []color.Color{
		color.RGBA{255, 0, 0, 255},   // Red
		color.RGBA{0, 255, 0, 255},   // Green
		color.RGBA{0, 0, 255, 255},   // Blue
		color.RGBA{255, 255, 0, 255}, // Yellow
	}

	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			var c color.Color
			if x < 50 && y < 50 {
				c = colors[0]
			} else if x >= 50 && y < 50 {
				c = colors[1]
			} else if x < 50 && y >= 50 {
				c = colors[2]
			} else {
				c = colors[3]
			}
			img.Set(x, y, c)
		}
	}

	return img
}
