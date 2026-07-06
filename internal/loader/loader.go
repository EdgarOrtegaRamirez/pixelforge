package loader

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/EdgarOrtegaRamirez/pixelforge/internal/models"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

// Load loads an image from the specified path
func Load(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer f.Close()

	img, format, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image (format: %s): %w", format, err)
	}

	return img, nil
}

// GetFormat returns the image format from file extension
func GetFormat(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return models.FormatJPEG
	case ".png":
		return models.FormatPNG
	case ".gif":
		return models.FormatGIF
	case ".bmp":
		return models.FormatBMP
	case ".tiff", ".tif":
		return models.FormatTIFF
	case ".webp":
		return models.FormatWEBP
	default:
		return "unknown"
	}
}

// GetInfo extracts metadata from an image
func GetInfo(path string) (*models.ImageInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat image: %w", err)
	}

	img, format, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Detect alpha channel
	hasAlpha := false
	switch img.(type) {
	case *image.NRGBA:
		hasAlpha = true
	case *image.NRGBA64:
		hasAlpha = true
	case *image.Alpha:
		hasAlpha = true
	case *image.Alpha16:
		hasAlpha = true
	case *image.Paletted:
		// Check if palette has transparent entries
		pal := img.(*image.Paletted).Palette
		for _, c := range pal {
			_, _, _, a := c.RGBA()
			if a < 0xFFFF {
				hasAlpha = true
				break
			}
		}
	}

	// Detect color model
	colorModel := "Unknown"
	cm := img.ColorModel()
	switch cm {
	case color.RGBAModel:
		colorModel = "RGBA"
	case color.RGBA64Model:
		colorModel = "RGBA64"
	case color.NRGBAModel:
		colorModel = "NRGBA"
	case color.GrayModel:
		colorModel = "Grayscale"
	case color.Gray16Model:
		colorModel = "Grayscale16"
	case color.CMYKModel:
		colorModel = "CMYK"
	default:
		// Check if it's a palette
		if _, ok := cm.(color.Palette); ok {
			colorModel = "Indexed"
		} else {
			colorModel = "Other"
		}
	}

	// Check for GIF animation
	animated := false
	frameCount := 1
	if format == "gif" {
		f2, err := os.Open(path)
		if err == nil {
			defer f2.Close()
			g, err := gif.DecodeAll(f2)
			if err == nil && len(g.Image) > 1 {
				animated = true
				frameCount = len(g.Image)
			}
		}
	}

	info := &models.ImageInfo{
		Path:       path,
		Format:     format,
		Width:      width,
		Height:     height,
		ColorModel: colorModel,
		HasAlpha:   hasAlpha,
		FileSize:   stat.Size(),
		ModifiedAt: stat.ModTime(),
		Animated:   animated,
		FrameCount: frameCount,
	}

	return info, nil
}

// Save saves an image to the specified path in the given format
func Save(img image.Image, path string, format string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	// Set file modification time
	os.Chtimes(path, time.Now(), time.Now())

	switch strings.ToLower(format) {
	case models.FormatJPEG:
		return jpeg.Encode(f, img, &jpeg.Options{Quality: 95})
	case models.FormatPNG:
		encoder := &png.Encoder{
			CompressionLevel: png.BestCompression,
		}
		return encoder.Encode(f, img)
	case models.FormatGIF:
		return gif.Encode(f, img, nil)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// SaveWithQuality saves an image with configurable quality
func SaveWithQuality(img image.Image, path string, format string, quality int) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	switch strings.ToLower(format) {
	case models.FormatJPEG:
		if quality < 1 {
			quality = 1
		}
		if quality > 100 {
			quality = 100
		}
		return jpeg.Encode(f, img, &jpeg.Options{Quality: quality})
	case models.FormatPNG:
		encoder := &png.Encoder{
			CompressionLevel: png.BestCompression,
		}
		return encoder.Encode(f, img)
	case models.FormatGIF:
		return gif.Encode(f, img, nil)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}
