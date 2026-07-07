package main

import (
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/EdgarOrtegaRamirez/pixelforge/internal/analyzer"
	"github.com/EdgarOrtegaRamirez/pixelforge/internal/loader"
	"github.com/EdgarOrtegaRamirez/pixelforge/internal/models"
	"github.com/EdgarOrtegaRamirez/pixelforge/internal/output"
	"github.com/EdgarOrtegaRamirez/pixelforge/internal/processor"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	verbose bool
	format  string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "pixelforge",
		Short: "Image processing and analysis toolkit",
		Long:  "A comprehensive CLI tool for image processing, analysis, and manipulation",
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "text", "output format (text, json)")

	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(convertCmd)
	rootCmd.AddCommand(resizeCmd)
	rootCmd.AddCommand(cropCmd)
	rootCmd.AddCommand(rotateCmd)
	rootCmd.AddCommand(thumbnailCmd)
	rootCmd.AddCommand(colorsCmd)
	rootCmd.AddCommand(compareCmd)
	rootCmd.AddCommand(effectCmd)
	rootCmd.AddCommand(stripCmd)
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var infoCmd = &cobra.Command{
	Use:   "info <image>",
	Short: "Display image metadata and information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", path)
		}

		info, err := loader.GetInfo(path)
		if err != nil {
			return err
		}

		out, err := output.FormatInfo(info, format)
		if err != nil {
			return err
		}
		fmt.Print(out)
		return nil
	},
}

var convertCmd = &cobra.Command{
	Use:   "convert <input> <output>",
	Short: "Convert image to a different format",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]
		outputPath := args[1]

		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			return fmt.Errorf("input file not found: %s", inputPath)
		}

		outputFormat := loader.GetFormat(outputPath)
		if outputFormat == "unknown" {
			return fmt.Errorf("cannot determine output format from extension: %s", outputPath)
		}

		img, err := loader.Load(inputPath)
		if err != nil {
			return err
		}

		quality, _ := cmd.Flags().GetInt("quality")
		if err := loader.SaveWithQuality(img, outputPath, outputFormat, quality); err != nil {
			return err
		}

		green := color.New(color.FgGreen)
		fmt.Printf("%s Converted %s → %s\n", green.Sprint("✓"), inputPath, outputPath)
		return nil
	},
}

var resizeCmd = &cobra.Command{
	Use:   "resize <input> <output>",
	Short: "Resize an image",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]
		outputPath := args[1]

		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			return fmt.Errorf("input file not found: %s", inputPath)
		}

		width, _ := cmd.Flags().GetInt("width")
		height, _ := cmd.Flags().GetInt("height")
		filter, _ := cmd.Flags().GetString("filter")
		keepAspect, _ := cmd.Flags().GetBool("aspect")

		img, err := loader.Load(inputPath)
		if err != nil {
			return err
		}

		result, err := processor.Resize(img, models.ResizeOptions{
			Width:      width,
			Height:     height,
			Filter:     filter,
			KeepAspect: keepAspect,
		})
		if err != nil {
			return err
		}

		outputFormat := loader.GetFormat(outputPath)
		if outputFormat == "unknown" {
			outputFormat = loader.GetFormat(inputPath)
		}

		if err := loader.Save(result, outputPath, outputFormat); err != nil {
			return err
		}

		green := color.New(color.FgGreen)
		fmt.Printf("%s Resized %s → %s\n", green.Sprint("✓"), inputPath, outputPath)
		return nil
	},
}

var cropCmd = &cobra.Command{
	Use:   "crop <input> <output>",
	Short: "Crop an image",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]
		outputPath := args[1]

		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			return fmt.Errorf("input file not found: %s", inputPath)
		}

		x, _ := cmd.Flags().GetInt("x")
		y, _ := cmd.Flags().GetInt("y")
		width, _ := cmd.Flags().GetInt("width")
		height, _ := cmd.Flags().GetInt("height")

		img, err := loader.Load(inputPath)
		if err != nil {
			return err
		}

		result, err := processor.Crop(img, models.CropOptions{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		})
		if err != nil {
			return err
		}

		outputFormat := loader.GetFormat(outputPath)
		if outputFormat == "unknown" {
			outputFormat = loader.GetFormat(inputPath)
		}

		if err := loader.Save(result, outputPath, outputFormat); err != nil {
			return err
		}

		green := color.New(color.FgGreen)
		fmt.Printf("%s Cropped %s → %s\n", green.Sprint("✓"), inputPath, outputPath)
		return nil
	},
}

var rotateCmd = &cobra.Command{
	Use:   "rotate <input> <output>",
	Short: "Rotate an image",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]
		outputPath := args[1]

		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			return fmt.Errorf("input file not found: %s", inputPath)
		}

		angle, _ := cmd.Flags().GetFloat64("angle")
		fillColor, _ := cmd.Flags().GetString("fill")

		img, err := loader.Load(inputPath)
		if err != nil {
			return err
		}

		result, err := processor.Rotate(img, models.RotateOptions{
			Angle:     angle,
			FillColor: fillColor,
		})
		if err != nil {
			return err
		}

		outputFormat := loader.GetFormat(outputPath)
		if outputFormat == "unknown" {
			outputFormat = loader.GetFormat(inputPath)
		}

		if err := loader.Save(result, outputPath, outputFormat); err != nil {
			return err
		}

		green := color.New(color.FgGreen)
		fmt.Printf("%s Rotated %s → %s\n", green.Sprint("✓"), inputPath, outputPath)
		return nil
	},
}

var thumbnailCmd = &cobra.Command{
	Use:   "thumbnail <input> <output>",
	Short: "Generate a thumbnail",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]
		outputPath := args[1]

		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			return fmt.Errorf("input file not found: %s", inputPath)
		}

		size, _ := cmd.Flags().GetInt("size")

		img, err := loader.Load(inputPath)
		if err != nil {
			return err
		}

		result, err := processor.Thumbnail(img, size)
		if err != nil {
			return err
		}

		outputFormat := loader.GetFormat(outputPath)
		if outputFormat == "unknown" {
			outputFormat = loader.GetFormat(inputPath)
		}

		if err := loader.Save(result, outputPath, outputFormat); err != nil {
			return err
		}

		green := color.New(color.FgGreen)
		fmt.Printf("%s Generated thumbnail %s → %s\n", green.Sprint("✓"), inputPath, outputPath)
		return nil
	},
}

var colorsCmd = &cobra.Command{
	Use:   "colors <image>",
	Short: "Extract dominant colors from an image",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", path)
		}

		count, _ := cmd.Flags().GetInt("count")

		img, err := loader.Load(path)
		if err != nil {
			return err
		}

		palette, err := analyzer.ExtractColors(img, count)
		if err != nil {
			return err
		}

		out, err := output.FormatColors(palette, format)
		if err != nil {
			return err
		}
		fmt.Print(out)
		return nil
	},
}

var compareCmd = &cobra.Command{
	Use:   "compare <image1> <image2>",
	Short: "Compare two images for differences",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		path1 := args[0]
		path2 := args[1]

		if _, err := os.Stat(path1); os.IsNotExist(err) {
			return fmt.Errorf("first image not found: %s", path1)
		}
		if _, err := os.Stat(path2); os.IsNotExist(err) {
			return fmt.Errorf("second image not found: %s", path2)
		}

		img1, err := loader.Load(path1)
		if err != nil {
			return err
		}

		img2, err := loader.Load(path2)
		if err != nil {
			return err
		}

		result, err := analyzer.Compare(img1, img2)
		if err != nil {
			return err
		}

		out, err := output.FormatCompare(result, format)
		if err != nil {
			return err
		}
		fmt.Print(out)
		return nil
	},
}

var effectCmd = &cobra.Command{
	Use:   "effect <input> <output> <effect-name>",
	Short: "Apply visual effects to an image",
	Long: `Available effects:
  grayscale    - Convert to grayscale
  invert       - Invert colors
  blur         - Apply Gaussian blur (use --sigma)
  sharpen      - Sharpen image (use --sigma)
  brightness   - Adjust brightness (use --factor, -1.0 to 1.0)
  contrast     - Adjust contrast (use --factor, -1.0 to 1.0)
  saturation   - Adjust color saturation (use --factor, -1.0 to 1.0)
  gamma        - Apply gamma correction (use --gamma)
  flip-h       - Flip horizontally
  flip-v       - Flip vertically
  rotate90     - Rotate 90° clockwise
  rotate180    - Rotate 180°
  rotate270    - Rotate 270° clockwise`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]
		outputPath := args[1]
		effectName := args[2]

		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			return fmt.Errorf("input file not found: %s", inputPath)
		}

		img, err := loader.Load(inputPath)
		if err != nil {
			return err
		}

		var result image.Image

		switch strings.ToLower(effectName) {
		case "grayscale":
			result = processor.Grayscale(img)
		case "invert":
			result = processor.Invert(img)
		case "blur":
			sigma, _ := cmd.Flags().GetFloat64("sigma")
			result = processor.Blur(img, sigma)
		case "sharpen":
			sigma, _ := cmd.Flags().GetFloat64("sigma")
			result = processor.Sharpen(img, sigma)
		case "brightness":
			factor, _ := cmd.Flags().GetFloat64("factor")
			result = processor.AdjustBrightness(img, factor)
		case "contrast":
			factor, _ := cmd.Flags().GetFloat64("factor")
			result = processor.AdjustContrast(img, factor)
		case "saturation":
			factor, _ := cmd.Flags().GetFloat64("factor")
			result = processor.AdjustSaturation(img, factor)
		case "gamma":
			gamma, _ := cmd.Flags().GetFloat64("gamma")
			result = processor.AdjustGamma(img, gamma)
		case "flip-h":
			result = processor.Flip(img, true)
		case "flip-v":
			result = processor.Flip(img, false)
		case "rotate90":
			result = processor.Rotate90(img)
		case "rotate180":
			result = processor.Rotate180(img)
		case "rotate270":
			result = processor.Rotate270(img)
		default:
			return fmt.Errorf("unknown effect: %s", effectName)
		}

		outputFormat := loader.GetFormat(outputPath)
		if outputFormat == "unknown" {
			outputFormat = loader.GetFormat(inputPath)
		}

		if err := loader.Save(result, outputPath, outputFormat); err != nil {
			return err
		}

		green := color.New(color.FgGreen)
		fmt.Printf("%s Applied %s to %s → %s\n", green.Sprint("✓"), effectName, inputPath, outputPath)
		return nil
	},
}

var stripCmd = &cobra.Command{
	Use:   "strip <input> <output>",
	Short: "Remove metadata from an image",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]
		outputPath := args[1]

		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			return fmt.Errorf("input file not found: %s", inputPath)
		}

		img, err := loader.Load(inputPath)
		if err != nil {
			return err
		}

		outputFormat := loader.GetFormat(outputPath)
		if outputFormat == "unknown" {
			outputFormat = loader.GetFormat(inputPath)
		}

		// Re-encoding without metadata effectively strips it
		if err := loader.Save(img, outputPath, outputFormat); err != nil {
			return err
		}

		green := color.New(color.FgGreen)
		fmt.Printf("%s Stripped metadata %s → %s\n", green.Sprint("✓"), inputPath, outputPath)
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("pixelforge v%s\n", version)
		fmt.Printf("Module: github.com/EdgarOrtegaRamirez/pixelforge\n")
	},
}
