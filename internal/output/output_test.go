package output

import (
	"strings"
	"testing"

	"github.com/EdgarOrtegaRamirez/pixelforge/internal/models"
)

func TestFormatInfo(t *testing.T) {
	info := &models.ImageInfo{
		Path:       "test.png",
		Format:     "png",
		Width:      100,
		Height:     100,
		ColorModel: "RGBA",
		HasAlpha:   false,
		FileSize:   1024,
	}

	// Test text format
	result, err := FormatInfo(info, "text")
	if err != nil {
		t.Fatalf("FormatInfo() error = %v", err)
	}

	if !strings.Contains(result, "test.png") {
		t.Error("FormatInfo() text format missing path")
	}
	if !strings.Contains(result, "100 × 100") {
		t.Error("FormatInfo() text format missing dimensions")
	}

	// Test JSON format
	result, err = FormatInfo(info, "json")
	if err != nil {
		t.Fatalf("FormatInfo() error = %v", err)
	}

	if !strings.Contains(result, "\"path\"") {
		t.Error("FormatInfo() JSON format missing path field")
	}

	// Test default format (empty string)
	result, err = FormatInfo(info, "")
	if err != nil {
		t.Fatalf("FormatInfo() error = %v", err)
	}

	if !strings.Contains(result, "test.png") {
		t.Error("FormatInfo() default format missing path")
	}

	// Test unsupported format
	_, err = FormatInfo(info, "xml")
	if err == nil {
		t.Error("FormatInfo() expected error for unsupported format")
	}
}

func TestFormatColors(t *testing.T) {
	palette := &models.ColorPalette{
		Colors: []models.Color{
			{R: 255, G: 0, B: 0, Count: 100, Hex: "#ff0000"},
			{R: 0, G: 255, B: 0, Count: 50, Hex: "#00ff00"},
			{R: 0, G: 0, B: 255, Count: 25, Hex: "#0000ff"},
		},
	}

	// Test text format
	result, err := FormatColors(palette, "text")
	if err != nil {
		t.Fatalf("FormatColors() error = %v", err)
	}

	if !strings.Contains(result, "#ff0000") {
		t.Error("FormatColors() text format missing hex color")
	}
	if !strings.Contains(result, "#00ff00") {
		t.Error("FormatColors() text format missing second color")
	}

	// Test JSON format
	result, err = FormatColors(palette, "json")
	if err != nil {
		t.Fatalf("FormatColors() error = %v", err)
	}

	if !strings.Contains(result, "\"colors\"") {
		t.Error("FormatColors() JSON format missing colors field")
	}

	// Test unsupported format
	_, err = FormatColors(palette, "xml")
	if err == nil {
		t.Error("FormatColors() expected error for unsupported format")
	}
}

func TestFormatCompare(t *testing.T) {
	result := &models.CompareResult{
		Identical:   false,
		SSIM:        0.85,
		PSNR:        30.5,
		PixelDiff:   1000,
		TotalPixels: 10000,
		DiffPercent: 10.0,
	}

	// Test text format
	text, err := FormatCompare(result, "text")
	if err != nil {
		t.Fatalf("FormatCompare() error = %v", err)
	}

	if !strings.Contains(text, "Different") {
		t.Error("FormatCompare() text format missing status")
	}
	if !strings.Contains(text, "1000") {
		t.Error("FormatCompare() text format missing pixel diff")
	}

	// Test JSON format
	jsonStr, err := FormatCompare(result, "json")
	if err != nil {
		t.Fatalf("FormatCompare() error = %v", err)
	}

	if !strings.Contains(jsonStr, "\"ssim\"") {
		t.Error("FormatCompare() JSON format missing ssim field")
	}

	// Test unsupported format
	_, err = FormatCompare(result, "xml")
	if err == nil {
		t.Error("FormatCompare() expected error for unsupported format")
	}
}

func TestFormatCompareIdentical(t *testing.T) {
	result := &models.CompareResult{
		Identical:   true,
		SSIM:        1.0,
		PSNR:        999999,
		PixelDiff:   0,
		TotalPixels: 10000,
		DiffPercent: 0,
	}

	text, err := FormatCompare(result, "text")
	if err != nil {
		t.Fatalf("FormatCompare() error = %v", err)
	}

	if !strings.Contains(text, "Identical") {
		t.Error("FormatCompare() text format missing identical status")
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.00 KB"},
		{1536, "1.50 KB"},
		{1048576, "1.00 MB"},
		{1073741824, "1.00 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatFileSize(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatFileSize(%d) = %q, want %q", tt.bytes, result, tt.expected)
			}
		})
	}
}
