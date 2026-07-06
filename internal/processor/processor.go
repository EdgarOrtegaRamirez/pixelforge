package processor

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/EdgarOrtegaRamirez/pixelforge/internal/models"
	"github.com/disintegration/imaging"
)

// Resize resizes an image to the specified dimensions
func Resize(img image.Image, opts models.ResizeOptions) (image.Image, error) {
	if opts.Width <= 0 && opts.Height <= 0 {
		return nil, fmt.Errorf("at least one dimension must be specified")
	}

	filter := getFilter(opts.Filter)

	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	dstW := opts.Width
	dstH := opts.Height

	if opts.KeepAspect {
		if dstW <= 0 {
			ratio := float64(dstH) / float64(srcH)
			dstW = int(float64(srcW) * ratio)
		} else if dstH <= 0 {
			ratio := float64(dstW) / float64(srcW)
			dstH = int(float64(srcH) * ratio)
		} else {
			// Both specified, fit within bounds
			ratioW := float64(dstW) / float64(srcW)
			ratioH := float64(dstH) / float64(srcH)
			ratio := math.Min(ratioW, ratioH)
			dstW = int(float64(srcW) * ratio)
			dstH = int(float64(srcH) * ratio)
		}
	} else {
		if dstW <= 0 {
			dstW = srcW
		}
		if dstH <= 0 {
			dstH = srcH
		}
	}

	if dstW <= 0 || dstH <= 0 {
		return nil, fmt.Errorf("computed dimensions too small: %dx%d", dstW, dstH)
	}

	return imaging.Resize(img, dstW, dstH, filter), nil
}

// Crop crops an image to the specified rectangle
func Crop(img image.Image, opts models.CropOptions) (image.Image, error) {
	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	x := opts.X
	y := opts.Y
	w := opts.Width
	h := opts.Height

	if x < 0 || y < 0 {
		return nil, fmt.Errorf("crop coordinates must be non-negative")
	}
	if w <= 0 || h <= 0 {
		return nil, fmt.Errorf("crop dimensions must be positive")
	}
	if x+w > srcW {
		w = srcW - x
	}
	if y+h > srcH {
		h = srcH - y
	}
	if w <= 0 || h <= 0 {
		return nil, fmt.Errorf("crop region is outside image bounds")
	}

	rect := image.Rect(x, y, x+w, y+h)
	return imaging.Crop(img, rect), nil
}

// Rotate rotates an image by the specified angle (in degrees)
func Rotate(img image.Image, opts models.RotateOptions) (image.Image, error) {
	// Parse fill color
	bgColor := color.RGBA{255, 255, 255, 255} // default white
	if opts.FillColor != "" {
		c, err := parseHexColor(opts.FillColor)
		if err != nil {
			return nil, fmt.Errorf("invalid fill color: %w", err)
		}
		bgColor = c
	}

	return imaging.Rotate(img, opts.Angle, bgColor), nil
}

// Thumbnail generates a thumbnail of the image
func Thumbnail(img image.Image, size int) (image.Image, error) {
	if size <= 0 {
		return nil, fmt.Errorf("thumbnail size must be positive")
	}
	return imaging.Thumbnail(img, size, size, imaging.Lanczos), nil
}

// Flip flips an image horizontally or vertically
func Flip(img image.Image, horizontal bool) image.Image {
	if horizontal {
		return imaging.FlipH(img)
	}
	return imaging.FlipV(img)
}

// Grayscale converts an image to grayscale
func Grayscale(img image.Image) image.Image {
	return imaging.Grayscale(img)
}

// Blur applies Gaussian blur to an image
func Blur(img image.Image, sigma float64) image.Image {
	if sigma < 0 {
		sigma = 0
	}
	return imaging.Blur(img, sigma)
}

// Sharpen applies sharpening to an image
func Sharpen(img image.Image, sigma float64) image.Image {
	if sigma < 0 {
		sigma = 0
	}
	return imaging.Sharpen(img, sigma)
}

// AdjustBrightness adjusts image brightness (-1.0 to 1.0)
func AdjustBrightness(img image.Image, factor float64) image.Image {
	return imaging.AdjustBrightness(img, factor)
}

// AdjustContrast adjusts image contrast (-1.0 to 1.0)
func AdjustContrast(img image.Image, factor float64) image.Image {
	return imaging.AdjustContrast(img, factor)
}

// AdjustSaturation adjusts color saturation (-1.0 to 1.0)
func AdjustSaturation(img image.Image, factor float64) image.Image {
	return imaging.AdjustSaturation(img, factor)
}

// Invert inverts the colors of an image
func Invert(img image.Image) image.Image {
	return imaging.Invert(img)
}

// AdjustGamma applies gamma correction
func AdjustGamma(img image.Image, gamma float64) image.Image {
	if gamma <= 0 {
		gamma = 1.0
	}
	return imaging.AdjustGamma(img, gamma)
}

// Rotate90 rotates an image 90 degrees clockwise
func Rotate90(img image.Image) image.Image {
	return imaging.Rotate90(img)
}

// Rotate180 rotates an image 180 degrees
func Rotate180(img image.Image) image.Image {
	return imaging.Rotate180(img)
}

// Rotate270 rotates an image 270 degrees clockwise
func Rotate270(img image.Image) image.Image {
	return imaging.Rotate270(img)
}

// Paste pastes one image onto another at the specified position
func Paste(dst, src image.Image, x, y int) image.Image {
	return imaging.Paste(dst, src, image.Pt(x, y))
}

// Overlay overlays an image with transparency onto another
func Overlay(bg, overlayImg image.Image, x, y int, opacity float64) image.Image {
	if opacity < 0 {
		opacity = 0
	}
	if opacity > 1 {
		opacity = 1
	}

	return imaging.Overlay(bg, overlayImg, image.Pt(x, y), opacity)
}

// getFilter returns the appropriate imaging filter
func getFilter(name string) imaging.ResampleFilter {
	switch name {
	case models.FilterNearest:
		return imaging.NearestNeighbor
	case models.FilterBilinear:
		return imaging.Linear
	case models.FilterBicubic:
		return imaging.CatmullRom
	case models.FilterLanczos:
		return imaging.Lanczos
	default:
		return imaging.Lanczos
	}
}

// parseHexColor parses a hex color string (#RGB, #RRGGBB, #RRGGBBAA)
func parseHexColor(hex string) (color.RGBA, error) {
	if len(hex) == 0 {
		return color.RGBA{}, fmt.Errorf("empty color string")
	}

	if hex[0] == '#' {
		hex = hex[1:]
	}

	var r, g, b, a uint8
	a = 255

	switch len(hex) {
	case 3:
		_, err := fmt.Sscanf(hex, "%1x%1x%1x", &r, &g, &b)
		if err != nil {
			return color.RGBA{}, err
		}
		r = r * 17
		g = g * 17
		b = b * 17
	case 6:
		_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
		if err != nil {
			return color.RGBA{}, err
		}
	case 8:
		_, err := fmt.Sscanf(hex, "%02x%02x%02x%02x", &r, &g, &b, &a)
		if err != nil {
			return color.RGBA{}, err
		}
	default:
		return color.RGBA{}, fmt.Errorf("invalid hex color length: %d", len(hex))
	}

	return color.RGBA{r, g, b, a}, nil
}
