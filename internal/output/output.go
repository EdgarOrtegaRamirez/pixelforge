package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/EdgarOrtegaRamirez/pixelforge/internal/models"
	"github.com/fatih/color"
)

// FormatInfo formats image info for display
func FormatInfo(info *models.ImageInfo, format string) (string, error) {
	switch format {
	case "json":
		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return "", err
		}
		return string(data), nil
	case "text", "":
		return formatInfoText(info), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// FormatColors formats a color palette for display
func FormatColors(palette *models.ColorPalette, format string) (string, error) {
	switch format {
	case "json":
		data, err := json.MarshalIndent(palette, "", "  ")
		if err != nil {
			return "", err
		}
		return string(data), nil
	case "text", "":
		return formatColorsText(palette), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// FormatCompare formats comparison results for display
func FormatCompare(result *models.CompareResult, format string) (string, error) {
	switch format {
	case "json":
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return "", err
		}
		return string(data), nil
	case "text", "":
		return formatCompareText(result), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func formatInfoText(info *models.ImageInfo) string {
	var sb strings.Builder

	titleColor := color.New(color.FgCyan, color.Bold)
	labelColor := color.New(color.FgYellow)
	valueColor := color.New(color.FgWhite)

	titleColor.Fprint(&sb, "Image Information\n")
	sb.WriteString(strings.Repeat("─", 40) + "\n")

	printField := func(label, value string) {
		labelColor.Fprintf(&sb, "  %-15s", label)
		valueColor.Fprintf(&sb, " %s\n", value)
	}

	printField("Path", info.Path)
	printField("Format", strings.ToUpper(info.Format))
	printField("Dimensions", fmt.Sprintf("%d × %d", info.Width, info.Height))
	printField("Color Model", info.ColorModel)

	if info.HasAlpha {
		printField("Alpha", "Yes")
	} else {
		printField("Alpha", "No")
	}

	printField("File Size", formatFileSize(info.FileSize))

	if info.Animated {
		printField("Animated", fmt.Sprintf("Yes (%d frames)", info.FrameCount))
	}

	if len(info.EXIF) > 0 {
		sb.WriteString("\n")
		titleColor.Fprint(&sb, "EXIF Data\n")
		sb.WriteString(strings.Repeat("─", 40) + "\n")
		for k, v := range info.EXIF {
			printField(k, v)
		}
	}

	return sb.String()
}

func formatColorsText(palette *models.ColorPalette) string {
	var sb strings.Builder

	titleColor := color.New(color.FgCyan, color.Bold)
	sb.WriteString(titleColor.Sprint("Dominant Colors\n"))
	sb.WriteString(strings.Repeat("─", 40) + "\n")

	for i, c := range palette.Colors {
		// Create colored block
		block := color.RGB(int(c.R), int(c.G), int(c.B)).Sprint("██")
		fmt.Fprintf(&sb, "  %2d. %s %s (rgb(%d,%d,%d))", i+1, block, c.Hex, c.R, c.G, c.B)
		if c.Count > 0 {
			fmt.Fprintf(&sb, " [%d pixels]", c.Count)
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func formatCompareText(result *models.CompareResult) string {
	var sb strings.Builder

	titleColor := color.New(color.FgCyan, color.Bold)
	labelColor := color.New(color.FgYellow)
	greenColor := color.New(color.FgGreen)
	redColor := color.New(color.FgRed)

	titleColor.Fprint(&sb, "Image Comparison\n")
	sb.WriteString(strings.Repeat("─", 40) + "\n")

	printField := func(label, value string) {
		labelColor.Fprintf(&sb, "  %-15s", label)
		fmt.Fprintf(&sb, " %s\n", value)
	}

	if result.Identical {
		printField("Status", greenColor.Sprint("✓ Identical"))
	} else {
		printField("Status", redColor.Sprint("✗ Different"))
	}

	printField("Pixel Diff", fmt.Sprintf("%d / %d", result.PixelDiff, result.TotalPixels))
	printField("Diff %", fmt.Sprintf("%.2f%%", result.DiffPercent))
	printField("SSIM", fmt.Sprintf("%.4f", result.SSIM))

	if result.PSNR < 1000 {
		printField("PSNR", fmt.Sprintf("%.2f dB", result.PSNR))
	} else {
		printField("PSNR", "∞ dB")
	}

	return sb.String()
}

func formatFileSize(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
