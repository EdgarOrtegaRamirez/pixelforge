# PixelForge

A comprehensive image processing and analysis toolkit for the command line.

## Features

- **Info** — Display image metadata (dimensions, format, color model, file size)
- **Convert** — Convert between image formats (PNG, JPEG, GIF, BMP, TIFF, WebP)
- **Resize** — Resize images with multiple filter options (Nearest, Bilinear, Bicubic, Lanczos)
- **Crop** — Crop images with precise coordinate control
- **Rotate** — Rotate images by any angle with configurable fill color
- **Thumbnail** — Generate thumbnails with aspect ratio preservation
- **Colors** — Extract dominant colors using k-means clustering
- **Compare** — Compare two images for similarity (SSIM, PSNR, pixel difference)
- **Effects** — Apply visual effects (blur, sharpen, brightness, contrast, saturation, gamma, grayscale, invert, flip, rotate)
- **Strip** — Remove EXIF and metadata for privacy
- **Version** — Print version information

## Installation

```bash
go install github.com/EdgarOrtegaRamirez/pixelforge/cmd/pixelforge@latest
```

Or build from source:

```bash
git clone https://github.com/EdgarOrtegaRamirez/pixelforge.git
cd pixelforge
go build -o pixelforge ./cmd/pixelforge/
```

## Quick Start

```bash
# Get image info
pixelforge info photo.jpg

# Convert PNG to JPEG
pixelforge convert image.png image.jpg

# Resize to 800px width (keep aspect ratio)
pixelforge resize photo.jpg resized.jpg --width=800 --aspect

# Crop a region
pixelforge crop photo.jpg cropped.jpg --x=100 --y=100 --width=400 --height=300

# Rotate 45 degrees
pixelforge rotate photo.jpg rotated.jpg --angle=45

# Generate thumbnail
pixelforge thumbnail photo.jpg thumb.jpg --size=200

# Extract dominant colors
pixelforge colors photo.jpg --count=5

# Compare two images
pixelforge compare image1.jpg image2.jpg

# Apply effects
pixelforge effect photo.jpg blur.jpg blur --sigma=2.0
pixelforge effect photo.jpg gray.jpg grayscale
pixelforge effect photo.jpg bright.jpg brightness --factor=0.3

# Remove metadata
pixelforge strip private.jpg clean.jpg
```

## Commands

### `pixelforge info <image>`

Display comprehensive image metadata.

```bash
pixelforge info photo.png
pixelforge info photo.png -f json  # JSON output
```

### `pixelforge convert <input> <output>`

Convert between image formats.

```bash
pixelforge convert image.png image.jpg
pixelforge convert image.jpg image.webp --quality=90
```

### `pixelforge resize <input> <output>`

Resize images with various options.

```bash
pixelforge resize photo.jpg resized.jpg --width=800 --height=600
pixelforge resize photo.jpg resized.jpg --width=800 --aspect  # Keep aspect ratio
pixelforge resize photo.jpg resized.jpg --height=400 --filter=bilinear
```

**Filters:**
- `nearest` — Nearest neighbor (fastest)
- `bilinear` — Bilinear interpolation
- `bicubic` — Bicubic interpolation
- `lanczos` — Lanczos resampling (default, highest quality)

### `pixelforge crop <input> <output>`

Crop a region from an image.

```bash
pixelforge crop photo.jpg cropped.jpg --x=100 --y=50 --width=400 --height=300
```

### `pixelforge rotate <input> <output>`

Rotate an image by any angle.

```bash
pixelforge rotate photo.jpg rotated.jpg --angle=90
pixelforge rotate photo.jpg rotated.jpg --angle=45 --fill="#ff0000"
```

### `pixelforge thumbnail <input> <output>`

Generate a square thumbnail.

```bash
pixelforge thumbnail photo.jpg thumb.jpg --size=150
```

### `pixelforge colors <image>`

Extract dominant colors using k-means clustering.

```bash
pixelforge colors photo.jpg --count=8
pixelforge colors photo.jpg -f json  # JSON output
```

### `pixelforge compare <image1> <image2>`

Compare two images for similarity.

```bash
pixelforge compare original.jpg modified.jpg
pixelforge compare a.png b.png -f json
```

**Output includes:**
- Identical status
- SSIM (Structural Similarity Index)
- PSNR (Peak Signal-to-Noise Ratio)
- Pixel difference count and percentage

### `pixelforge effect <input> <output> <effect>`

Apply visual effects to images.

```bash
pixelforge effect photo.jpg blur.jpg blur --sigma=3.0
pixelforge effect photo.jpg sharp.jpg sharpen --sigma=1.5
pixelforge effect photo.jpg bright.jpg brightness --factor=0.4
pixelforge effect photo.jpg contrast.jpg contrast --factor=-0.3
pixelforge effect photo.jpg sat.jpg saturation --factor=0.5
pixelforge effect photo.jpg gamma.jpg gamma --gamma=2.2
pixelforge effect photo.jpg gray.jpg grayscale
pixelforge effect photo.jpg inv.jpg invert
pixelforge effect photo.jpg fh.jpg flip-h
pixelforge effect photo.jpg fv.jpg flip-v
pixelforge effect photo.jpg r90.jpg rotate90
pixelforge effect photo.jpg r180.jpg rotate180
pixelforge effect photo.jpg r270.jpg rotate270
```

### `pixelforge strip <input> <output>`

Remove all metadata (EXIF, IPTC, etc.) from an image.

```bash
pixelforge strip private.jpg clean.jpg
```

### `pixelforge version`

Print version information.

```bash
pixelforge version
```

## Global Flags

- `-f, --format` — Output format: `text` (default) or `json`
- `-v, --verbose` — Verbose output

## Supported Formats

| Format | Extension(s) | Read | Write |
|--------|--------------|------|-------|
| JPEG   | .jpg, .jpeg  | ✅   | ✅    |
| PNG    | .png         | ✅   | ✅    |
| GIF    | .gif         | ✅   | ✅    |
| BMP    | .bmp         | ✅   | ❌    |
| TIFF   | .tiff, .tif  | ✅   | ❌    |
| WebP   | .webp        | ✅   | ❌    |

## Architecture

```
pixelforge/
├── cmd/pixelforge/       # CLI entry point (Cobra)
├── internal/
│   ├── models/           # Data structures
│   ├── loader/           # Image loading, saving, format detection
│   ├── processor/        # Image transformations (resize, crop, rotate, effects)
│   ├── analyzer/         # Color extraction, image comparison, histograms
│   └── output/           # Output formatting (text, JSON)
└── tests/                # Test fixtures
```

## Key Algorithms

1. **K-means++ Clustering** — For dominant color extraction with smart initialization
2. **Euclidean Distance** — For color similarity in clustering and comparison
3. **SSIM Approximation** — Structural similarity measurement between images
4. **PSNR Calculation** — Peak signal-to-noise ratio for image quality
5. **Lanczos Resampling** — High-quality image resizing filter

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run specific package
go test ./internal/processor/... -v
```

## License

MIT License
