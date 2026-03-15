// Package imaging provides a drop-in replacement for github.com/disintegration/imaging.
//
// To migrate, change your import path from:
//
//	import "github.com/disintegration/imaging"
//
// to:
//
//	import "github.com/agentine/pixo/compat/imaging"
package imaging

import (
	"image"
	"image/color"
	"image/png"
	"io"

	"github.com/agentine/pixo"
)

// ResampleFilter is an interpolation filter for resizing.
type ResampleFilter int

const (
	// NearestNeighbor interpolation.
	NearestNeighbor ResampleFilter = iota
	// Box interpolation.
	Box
	// Linear interpolation.
	Linear
	// CatmullRom interpolation.
	CatmullRom
	// MitchellNetravali interpolation.
	MitchellNetravali
	// Lanczos interpolation.
	Lanczos
)

func toPixoFilter(f ResampleFilter) pixo.Filter {
	switch f {
	case NearestNeighbor:
		return pixo.NearestNeighbor
	case Box:
		return pixo.Box
	case Linear:
		return pixo.Linear
	case CatmullRom:
		return pixo.CatmullRom
	case MitchellNetravali:
		return pixo.MitchellNetravali
	case Lanczos:
		return pixo.Lanczos3
	default:
		return pixo.Lanczos3
	}
}

// Anchor is the anchor point for cropping.
type Anchor int

const (
	// Center anchor.
	Center Anchor = iota
	// TopLeft anchor.
	TopLeft
	// Top anchor.
	Top
	// TopRight anchor.
	TopRight
	// Left anchor.
	Left
	// Right anchor.
	Right
	// BottomLeft anchor.
	BottomLeft
	// Bottom anchor.
	Bottom
	// BottomRight anchor.
	BottomRight
)

func toPixoAnchor(a Anchor) pixo.Anchor {
	switch a {
	case Center:
		return pixo.Center
	case TopLeft:
		return pixo.TopLeft
	case Top:
		return pixo.Top
	case TopRight:
		return pixo.TopRight
	case Left:
		return pixo.Left
	case Right:
		return pixo.Right
	case BottomLeft:
		return pixo.BottomLeft
	case Bottom:
		return pixo.Bottom
	case BottomRight:
		return pixo.BottomRight
	default:
		return pixo.Center
	}
}

// Format is the image file format.
type Format = pixo.Format

const (
	JPEG = pixo.JPEG
	PNG  = pixo.PNG
	GIF  = pixo.GIF
	BMP  = pixo.BMP
	TIFF = pixo.TIFF
)

// Resize resizes img to the specified dimensions.
func Resize(img image.Image, width, height int, filter ResampleFilter) *image.NRGBA {
	return pixo.Resize(img, width, height, toPixoFilter(filter))
}

// Fit scales img to fit within the specified dimensions.
func Fit(img image.Image, width, height int, filter ResampleFilter) *image.NRGBA {
	return pixo.Fit(img, width, height, toPixoFilter(filter))
}

// Fill scales and crops img to fill the specified dimensions.
func Fill(img image.Image, width, height int, anchor Anchor, filter ResampleFilter) *image.NRGBA {
	return pixo.Fill(img, width, height, toPixoFilter(filter), toPixoAnchor(anchor))
}

// Thumbnail scales img to fill the specified dimensions, cropping from center.
func Thumbnail(img image.Image, width, height int, filter ResampleFilter) *image.NRGBA {
	return pixo.Thumbnail(img, width, height, toPixoFilter(filter))
}

// Crop crops img to the specified rectangle.
func Crop(img image.Image, rect image.Rectangle) *image.NRGBA {
	return pixo.Crop(img, rect)
}

// CropAnchor crops img to the specified dimensions using the anchor point.
func CropAnchor(img image.Image, width, height int, anchor Anchor) *image.NRGBA {
	return pixo.CropAnchor(img, width, height, toPixoAnchor(anchor))
}

// CropCenter crops img to the specified dimensions from center.
func CropCenter(img image.Image, width, height int) *image.NRGBA {
	return pixo.CropCenter(img, width, height)
}

// Rotate rotates img by the given angle with the specified background color.
func Rotate(img image.Image, angle float64, bgcolor color.Color) *image.NRGBA {
	return pixo.Rotate(img, angle, bgcolor)
}

// Rotate90 rotates img 90 degrees counter-clockwise.
func Rotate90(img image.Image) *image.NRGBA {
	return pixo.Rotate90(img)
}

// Rotate180 rotates img 180 degrees.
func Rotate180(img image.Image) *image.NRGBA {
	return pixo.Rotate180(img)
}

// Rotate270 rotates img 270 degrees counter-clockwise.
func Rotate270(img image.Image) *image.NRGBA {
	return pixo.Rotate270(img)
}

// FlipH flips img horizontally.
func FlipH(img image.Image) *image.NRGBA {
	return pixo.FlipH(img)
}

// FlipV flips img vertically.
func FlipV(img image.Image) *image.NRGBA {
	return pixo.FlipV(img)
}

// Transpose flips img diagonally.
func Transpose(img image.Image) *image.NRGBA {
	return pixo.Transpose(img)
}

// Transverse flips img along the other diagonal.
func Transverse(img image.Image) *image.NRGBA {
	return pixo.Transverse(img)
}

// AdjustBrightness adjusts brightness by the given percentage [-100, 100].
func AdjustBrightness(img image.Image, percentage float64) *image.NRGBA {
	return pixo.AdjustBrightness(img, percentage)
}

// AdjustContrast adjusts contrast by the given percentage [-100, 100].
func AdjustContrast(img image.Image, percentage float64) *image.NRGBA {
	return pixo.AdjustContrast(img, percentage)
}

// AdjustGamma adjusts gamma. Gamma must be > 0.
func AdjustGamma(img image.Image, gamma float64) *image.NRGBA {
	return pixo.AdjustGamma(img, gamma)
}

// AdjustSaturation adjusts saturation by the given percentage [-100, 100].
func AdjustSaturation(img image.Image, percentage float64) *image.NRGBA {
	return pixo.AdjustSaturation(img, percentage)
}

// AdjustSigmoid adjusts contrast using a sigmoidal function.
func AdjustSigmoid(img image.Image, midpoint, factor float64) *image.NRGBA {
	return pixo.AdjustSigmoid(img, midpoint, factor)
}

// Blur applies a Gaussian blur with the given sigma.
func Blur(img image.Image, sigma float64) *image.NRGBA {
	return pixo.Blur(img, sigma)
}

// Sharpen sharpens img using unsharp masking.
func Sharpen(img image.Image, sigma float64) *image.NRGBA {
	return pixo.Sharpen(img, sigma)
}

// Grayscale converts img to grayscale.
func Grayscale(img image.Image) *image.NRGBA {
	return pixo.Grayscale(img)
}

// Invert inverts the colors of img.
func Invert(img image.Image) *image.NRGBA {
	return pixo.Invert(img)
}

// Paste pastes img onto background at pos.
func Paste(background image.Image, img image.Image, pos image.Point) *image.NRGBA {
	return pixo.Paste(background, img, pos)
}

// PasteCenter pastes img onto the center of background.
func PasteCenter(background image.Image, img image.Image) *image.NRGBA {
	return pixo.PasteCenter(background, img)
}

// Overlay draws img over background with opacity.
func Overlay(background image.Image, img image.Image, pos image.Point, opacity float64) *image.NRGBA {
	return pixo.Overlay(background, img, pos, opacity)
}

// OverlayCenter draws img over the center of background with opacity.
func OverlayCenter(background image.Image, img image.Image, opacity float64) *image.NRGBA {
	return pixo.OverlayCenter(background, img, opacity)
}

// Histogram returns per-channel histogram.
func Histogram(img image.Image) [4][256]int {
	return pixo.Histogram(img)
}

// Open opens and decodes an image file.
func Open(filename string) (image.Image, error) {
	return pixo.Open(filename)
}

// Save saves img to the given file.
func Save(img image.Image, filename string) error {
	return pixo.Save(img, filename)
}

// Decode decodes an image from a reader.
func Decode(r io.Reader) (image.Image, error) {
	return pixo.Decode(r)
}

// Encode encodes an image to a writer.
func Encode(w io.Writer, img image.Image, format Format, opts ...pixo.EncodeOption) error {
	return pixo.Encode(w, img, format, opts...)
}

// JPEGQuality returns an EncodeOption that sets JPEG quality.
func JPEGQuality(q int) pixo.EncodeOption {
	return pixo.JPEGQuality(q)
}

// PNGCompressionLevel returns an EncodeOption that sets PNG compression.
func PNGCompressionLevel(l png.CompressionLevel) pixo.EncodeOption {
	return pixo.PNGCompressionLevel(l)
}
