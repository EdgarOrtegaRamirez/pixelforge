package analyzer

import (
	"image"
	"image/color"
	"testing"
)

func TestExtractColors(t *testing.T) {
	img := createTestImage()

	palette, err := ExtractColors(img, 4)
	if err != nil {
		t.Fatalf("ExtractColors() error = %v", err)
	}

	if len(palette.Colors) != 4 {
		t.Errorf("ExtractColors() returned %d colors, want 4", len(palette.Colors))
	}

	// Colors should be sorted by count
	for i := 1; i < len(palette.Colors); i++ {
		if palette.Colors[i].Count > palette.Colors[i-1].Count {
			t.Errorf("Colors not sorted by count at index %d", i)
		}
	}
}

func TestExtractColorsMoreThanPixels(t *testing.T) {
	img := createTestImage()

	// Request more colors than there are unique colors
	// The function caps at 256, and k-means may return more clusters
	palette, err := ExtractColors(img, 1000)
	if err != nil {
		t.Fatalf("ExtractColors() error = %v", err)
	}

	// Should return at most 256 colors (capped by the function)
	if len(palette.Colors) > 256 {
		t.Errorf("ExtractColors() returned %d colors, expected at most 256", len(palette.Colors))
	}
}

func TestExtractColorsInvalidCount(t *testing.T) {
	img := createTestImage()

	// Count of 0 should use default
	palette, err := ExtractColors(img, 0)
	if err != nil {
		t.Fatalf("ExtractColors() error = %v", err)
	}

	if len(palette.Colors) != 5 {
		t.Errorf("ExtractColors() with count=0 returned %d colors, want 5", len(palette.Colors))
	}
}

func TestCompare(t *testing.T) {
	img1 := createTestImage()
	img2 := createTestImage2()

	result, err := Compare(img1, img2)
	if err != nil {
		t.Fatalf("Compare() error = %v", err)
	}

	if result.Identical {
		t.Error("Compare() identical = true, want false")
	}

	if result.SSIM >= 1.0 {
		t.Errorf("Compare() SSIM = %f, want < 1.0", result.SSIM)
	}

	if result.PixelDiff <= 0 {
		t.Errorf("Compare() PixelDiff = %d, want > 0", result.PixelDiff)
	}

	if result.DiffPercent <= 0 {
		t.Errorf("Compare() DiffPercent = %f, want > 0", result.DiffPercent)
	}
}

func TestCompareIdentical(t *testing.T) {
	img1 := createTestImage()
	img2 := createTestImage()

	result, err := Compare(img1, img2)
	if err != nil {
		t.Fatalf("Compare() error = %v", err)
	}

	if !result.Identical {
		t.Error("Compare() identical = false, want true")
	}

	if result.SSIM != 1.0 {
		t.Errorf("Compare() SSIM = %f, want 1.0", result.SSIM)
	}

	if result.PixelDiff != 0 {
		t.Errorf("Compare() PixelDiff = %d, want 0", result.PixelDiff)
	}
}

func TestCompareDifferentDimensions(t *testing.T) {
	img1 := createTestImage()
	img2 := image.NewRGBA(image.Rect(0, 0, 50, 50))

	_, err := Compare(img1, img2)
	if err == nil {
		t.Error("Compare() expected error for different dimensions")
	}
}

func TestGetHistogram(t *testing.T) {
	img := createTestImage()

	histogram := GetHistogram(img)

	// Check that all channels are present
	for _, channel := range []string{"red", "green", "blue", "luma"} {
		if _, ok := histogram[channel]; !ok {
			t.Errorf("GetHistogram() missing channel: %s", channel)
		}
	}

	// Red channel should have values
	totalRed := 0
	for _, v := range histogram["red"] {
		totalRed += v
	}
	if totalRed != 10000 { // 100x100 image
		t.Errorf("GetHistogram() red total = %d, want 10000", totalRed)
	}
}

func TestGetAverageColor(t *testing.T) {
	// Create image with known average
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{100, 150, 200, 255})
		}
	}

	avg := GetAverageColor(img)

	// Average should be close to the set color
	if avg.R != 100 || avg.G != 150 || avg.B != 200 {
		t.Errorf("GetAverageColor() = rgba(%d,%d,%d), want rgba(100,150,200)", avg.R, avg.G, avg.B)
	}
}

func TestGetAverageColorEmpty(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))

	avg := GetAverageColor(img)

	// Should return black for empty image
	if avg.R != 0 || avg.G != 0 || avg.B != 0 {
		t.Errorf("GetAverageColor() for empty image = rgba(%d,%d,%d), want rgba(0,0,0)", avg.R, avg.G, avg.B)
	}
}

func TestKMeans(t *testing.T) {
	// Create simple pixel data
	pixels := [][3]float64{
		{0, 0, 0},       // Black
		{0, 0, 0},       // Black
		{0, 0, 0},       // Black
		{255, 255, 255}, // White
		{255, 255, 255}, // White
		{255, 255, 255}, // White
	}

	centroids := kMeans(pixels, 2, 10)

	if len(centroids) != 2 {
		t.Errorf("kMeans() returned %d centroids, want 2", len(centroids))
	}

	// Centroids should be close to black and white
	foundBlack := false
	foundWhite := false
	for _, c := range centroids {
		if c[0] < 50 && c[1] < 50 && c[2] < 50 {
			foundBlack = true
		}
		if c[0] > 200 && c[1] > 200 && c[2] > 200 {
			foundWhite = true
		}
	}

	if !foundBlack || !foundWhite {
		t.Error("kMeans() did not find expected black and white centroids")
	}
}

func TestDistance(t *testing.T) {
	tests := []struct {
		a, b [3]float64
		want float64
	}{
		{[3]float64{0, 0, 0}, [3]float64{0, 0, 0}, 0},
		{[3]float64{0, 0, 0}, [3]float64{1, 0, 0}, 1},
		{[3]float64{0, 0, 0}, [3]float64{0, 1, 0}, 1},
		{[3]float64{0, 0, 0}, [3]float64{0, 0, 1}, 1},
		{[3]float64{0, 0, 0}, [3]float64{1, 1, 1}, 1.732}, // sqrt(3)
	}

	for _, tt := range tests {
		result := distance(tt.a, tt.b)
		if result < tt.want-0.01 || result > tt.want+0.01 {
			t.Errorf("distance(%v, %v) = %f, want %f", tt.a, tt.b, result, tt.want)
		}
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

func createTestImage2() *image.RGBA {
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
				// Different color
				c = color.RGBA{255, 0, 255, 255}
			}
			img.Set(x, y, c)
		}
	}

	return img
}
