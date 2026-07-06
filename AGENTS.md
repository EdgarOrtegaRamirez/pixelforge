# AGENTS.md

## Project Overview

PixelForge is a comprehensive image processing and analysis CLI toolkit written in Go. It provides commands for image transformation, color extraction, comparison, and metadata inspection.

## Architecture

```
pixelforge/
├── cmd/pixelforge/         # CLI entry point using Cobra
│   └── main.go
├── internal/
│   ├── models/             # Data structures and types
│   │   └── models.go
│   ├── loader/             # Image I/O and format detection
│   │   ├── loader.go       # Load, Save, GetInfo, GetFormat
│   │   └── loader_test.go
│   ├── processor/          # Image transformations
│   │   ├── processor.go    # Resize, Crop, Rotate, Effects
│   │   └── processor_test.go
│   ├── analyzer/           # Image analysis
│   │   ├── analyzer.go     # Color extraction, Compare, Histogram
│   │   └── analyzer_test.go
│   └── output/             # Output formatting
│       ├── output.go       # Text and JSON formatters
│       └── output_test.go
├── tests/                  # Test fixtures
├── go.mod
├── go.sum
└── README.md
```

## Key Components

### Models (`internal/models/models.go`)
- `ImageInfo` — Image metadata (dimensions, format, EXIF, etc.)
- `Color`, `ColorPalette` — Color extraction results
- `ResizeOptions`, `CropOptions`, `RotateOptions` — Transformation configs
- `CompareResult` — Image comparison results
- Filter and position constants

### Loader (`internal/loader/loader.go`)
- `Load(path)` — Load image from file (auto-detects format)
- `Save(img, path, format)` — Save image to file
- `SaveWithQuality(img, path, format, quality)` — Save with quality control
- `GetInfo(path)` — Extract image metadata
- `GetFormat(path)` — Detect format from extension

### Processor (`internal/processor/processor.go`)
- `Resize(img, opts)` — Resize with aspect ratio control
- `Crop(img, opts)` — Crop with coordinate clamping
- `Rotate(img, opts)` — Rotate with fill color
- `Thumbnail(img, size)` — Generate thumbnails
- `Flip`, `Grayscale`, `Blur`, `Sharpen` — Common operations
- `AdjustBrightness`, `AdjustContrast`, `AdjustSaturation`, `AdjustGamma` — Adjustments
- `Rotate90`, `Rotate180`, `Rotate270` — Quick rotations
- `Paste`, `Overlay` — Compositing

### Analyzer (`internal/analyzer/analyzer.go`)
- `ExtractColors(img, count)` — K-means++ color clustering
- `Compare(img1, img2)` — SSIM, PSNR, pixel diff
- `GetHistogram(img)` — RGB + luma histograms
- `GetAverageColor(img)` — Mean color calculation

### Output (`internal/output/output.go`)
- `FormatInfo`, `FormatColors`, `FormatCompare` — Multi-format output
- Text with colors (via `github.com/fatih/color`)
- JSON for machine consumption

## Building & Testing

```bash
# Build
go build -o pixelforge ./cmd/pixelforge/

# Test all
go test ./...

# Test specific package
go test ./internal/processor/... -v

# Vet
go vet ./...
```

## Key Algorithms

1. **K-means++ Clustering** — `analyzer.go:kMeans()` with smart initialization
2. **Euclidean Distance** — `analyzer.go:distance()` for color similarity
3. **SSIM Approximation** — `analyzer.go:Compare()` for structural similarity
4. **PSNR Calculation** — `analyzer.go:Compare()` for image quality
5. **Hex Color Parsing** — `processor.go:parseHexColor()` supports #RGB, #RRGGBB, #RRGGBBAA

## Adding New Effects

1. Add effect function in `processor.go`
2. Add case in `effectCmd` RunE in `cmd/pixelforge/main.go`
3. Add test in `processor_test.go`
4. Update README.md effects section

## Dependencies

- `github.com/disintegration/imaging` — Core image processing
- `github.com/spf13/cobra` — CLI framework
- `github.com/fatih/color` — Terminal colors
- `golang.org/x/image` — Additional format support (BMP, TIFF, WebP)

## Code Style

- Use `gofmt` for formatting
- Run `go vet` before committing
- Tests use standard `testing` package
- Table-driven tests where appropriate
