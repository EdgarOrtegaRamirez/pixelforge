package analyzer

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sort"

	"github.com/EdgarOrtegaRamirez/pixelforge/internal/models"
)

// ExtractColors extracts the dominant colors from an image using k-means clustering
func ExtractColors(img image.Image, count int) (*models.ColorPalette, error) {
	if count <= 0 {
		count = 5
	}
	if count > 256 {
		count = 256
	}

	// Sample pixels
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Collect all pixels
	pixels := make([][3]float64, 0, width*height)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Convert from [0,65535] to [0,255]
			pixels = append(pixels, [3]float64{
				float64(r >> 8),
				float64(g >> 8),
				float64(b >> 8),
			})
		}
	}

	if len(pixels) == 0 {
		return nil, fmt.Errorf("image has no pixels")
	}

	// Run k-means clustering
	colors := kMeans(pixels, count, 20)

	// Count pixels in each cluster
	palette := &models.ColorPalette{
		Colors: make([]models.Color, 0, len(colors)),
	}

	for _, c := range colors {
		r := uint8(math.Round(c[0]))
		g := uint8(math.Round(c[1]))
		b := uint8(math.Round(c[2]))
		hex := fmt.Sprintf("#%02x%02x%02x", r, g, b)

		// Count pixels in this cluster
		count := 0
		for _, p := range pixels {
			if distance(p, c) < 30 { // threshold
				count++
			}
		}

		palette.Colors = append(palette.Colors, models.Color{
			R:     r,
			G:     g,
			B:     b,
			Count: count,
			Hex:   hex,
		})
	}

	// Sort by count (most dominant first)
	sort.Slice(palette.Colors, func(i, j int) bool {
		return palette.Colors[i].Count > palette.Colors[j].Count
	})

	return palette, nil
}

// Compare compares two images and returns similarity metrics
func Compare(img1, img2 image.Image) (*models.CompareResult, error) {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()

	// Check if dimensions match
	if bounds1.Dx() != bounds2.Dx() || bounds1.Dy() != bounds2.Dy() {
		return nil, fmt.Errorf("image dimensions do not match: %dx%d vs %dx%d",
			bounds1.Dx(), bounds1.Dy(), bounds2.Dx(), bounds2.Dy())
	}

	width := bounds1.Dx()
	height := bounds1.Dy()
	totalPixels := width * height

	if totalPixels == 0 {
		return &models.CompareResult{
			Identical: true,
			SSIM:      1.0,
			PSNR:      math.MaxFloat64,
		}, nil
	}

	// Calculate pixel-level difference
	var sumSquaredDiff float64
	pixelDiff := 0

	for y := bounds1.Min.Y; y < bounds1.Max.Y; y++ {
		for x := bounds1.Min.X; x < bounds1.Max.X; x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()

			// Convert to [0,255]
			r1f := float64(r1 >> 8)
			g1f := float64(g1 >> 8)
			b1f := float64(b1 >> 8)
			r2f := float64(r2 >> 8)
			g2f := float64(g2 >> 8)
			b2f := float64(b2 >> 8)

			// Calculate squared differences
			dr := r1f - r2f
			dg := g1f - g2f
			db := b1f - b2f

			sumSquaredDiff += dr*dr + dg*dg + db*db

			// Count differing pixels (threshold of 10 per channel)
			if math.Abs(dr) > 10 || math.Abs(dg) > 10 || math.Abs(db) > 10 {
				pixelDiff++
			}
		}
	}

	// Calculate PSNR (Peak Signal-to-Noise Ratio)
	mse := sumSquaredDiff / float64(totalPixels*3)
	var psnr float64
	if mse > 0 {
		psnr = 10 * math.Log10(255.0*255.0/mse)
	} else {
		psnr = math.MaxFloat64
	}

	// Calculate SSIM approximation (simplified)
	// For a proper SSIM, we'd need luminance, contrast, and structure comparisons
	ssim := 1.0 - (float64(pixelDiff) / float64(totalPixels))
	if ssim < 0 {
		ssim = 0
	}

	return &models.CompareResult{
		Identical:   pixelDiff == 0,
		SSIM:        ssim,
		PSNR:        psnr,
		PixelDiff:   pixelDiff,
		TotalPixels: totalPixels,
		DiffPercent: float64(pixelDiff) / float64(totalPixels) * 100,
	}, nil
}

// GetHistogram returns the color histogram of an image
func GetHistogram(img image.Image) map[string][256]int {
	bounds := img.Bounds()
	red := [256]int{}
	green := [256]int{}
	blue := [256]int{}
	luma := [256]int{}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// Convert to [0,255]
			r8 := r >> 8
			g8 := g >> 8
			b8 := b >> 8

			red[r8]++
			green[g8]++
			blue[b8]++

			// Luma (BT.601)
			l := uint8((0.299*float64(r8) + 0.587*float64(g8) + 0.114*float64(b8)))
			luma[l]++
		}
	}

	return map[string][256]int{
		"red":   red,
		"green": green,
		"blue":  blue,
		"luma":  luma,
	}
}

// GetAverageColor returns the average color of an image
func GetAverageColor(img image.Image) color.RGBA {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	totalPixels := width * height

	if totalPixels == 0 {
		return color.RGBA{0, 0, 0, 255}
	}

	var sumR, sumG, sumB float64

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			sumR += float64(r >> 8)
			sumG += float64(g >> 8)
			sumB += float64(b >> 8)
		}
	}

	return color.RGBA{
		uint8(math.Round(sumR / float64(totalPixels))),
		uint8(math.Round(sumG / float64(totalPixels))),
		uint8(math.Round(sumB / float64(totalPixels))),
		255,
	}
}

// kMeans performs k-means clustering on pixel data
func kMeans(pixels [][3]float64, k, maxIter int) [][3]float64 {
	n := len(pixels)
	if n == 0 {
		return nil
	}
	if k > n {
		k = n
	}

	// Initialize centroids using k-means++
	centroids := make([][3]float64, k)
	centroids[0] = pixels[0]

	for i := 1; i < k; i++ {
		// Calculate distances to nearest centroid
		dists := make([]float64, n)
		var totalDist float64
		for j, p := range pixels {
			minDist := math.MaxFloat64
			for c := 0; c < i; c++ {
				d := distance(p, centroids[c])
				if d < minDist {
					minDist = d
				}
			}
			dists[j] = minDist
			totalDist += minDist
		}

		// Choose next centroid probabilistically
		r := math.Mod(float64(i*12345+n*67890), totalDist)
		cumDist := 0.0
		for j, d := range dists {
			cumDist += d
			if cumDist >= r {
				centroids[i] = pixels[j]
				break
			}
		}
	}

	// Iterate
	for iter := 0; iter < maxIter; iter++ {
		// Assign pixels to nearest centroid
		assignments := make([]int, n)
		for j, p := range pixels {
			minDist := math.MaxFloat64
			bestC := 0
			for c, centroid := range centroids {
				d := distance(p, centroid)
				if d < minDist {
					minDist = d
					bestC = c
				}
			}
			assignments[j] = bestC
		}

		// Update centroids
		newCentroids := make([][3]float64, k)
		counts := make([]int, k)
		for j, p := range pixels {
			c := assignments[j]
			newCentroids[c][0] += p[0]
			newCentroids[c][1] += p[1]
			newCentroids[c][2] += p[2]
			counts[c]++
		}

		converged := true
		for c := 0; c < k; c++ {
			if counts[c] > 0 {
				newCentroids[c][0] /= float64(counts[c])
				newCentroids[c][1] /= float64(counts[c])
				newCentroids[c][2] /= float64(counts[c])
			}
			if distance(newCentroids[c], centroids[c]) > 1 {
				converged = false
			}
		}
		centroids = newCentroids

		if converged {
			break
		}
	}

	return centroids
}

// distance calculates Euclidean distance between two RGB colors
func distance(a, b [3]float64) float64 {
	return math.Sqrt((a[0]-b[0])*(a[0]-b[0]) + (a[1]-b[1])*(a[1]-b[1]) + (a[2]-b[2])*(a[2]-b[2]))
}
