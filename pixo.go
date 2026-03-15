// Package pixo provides image processing functions.
//
// All functions accept image.Image and return *image.NRGBA.
package pixo

// Filter is a resampling filter for image resizing.
type Filter int

const (
	// NearestNeighbor is a nearest-neighbor interpolation filter.
	NearestNeighbor Filter = iota
	// Box is a box (averaging) interpolation filter.
	Box
	// Linear is a bilinear interpolation filter.
	Linear
	// CatmullRom is a Catmull-Rom (bicubic) interpolation filter.
	CatmullRom
	// MitchellNetravali is a Mitchell-Netravali (bicubic) interpolation filter.
	MitchellNetravali
	// Lanczos2 is a Lanczos resampling filter with a=2.
	Lanczos2
	// Lanczos3 is a Lanczos resampling filter with a=3.
	Lanczos3
)

// Anchor is the anchor point for cropping and filling operations.
type Anchor int

const (
	// Center is the center anchor point.
	Center Anchor = iota
	// TopLeft is the top-left anchor point.
	TopLeft
	// Top is the top-center anchor point.
	Top
	// TopRight is the top-right anchor point.
	TopRight
	// Left is the center-left anchor point.
	Left
	// Right is the center-right anchor point.
	Right
	// BottomLeft is the bottom-left anchor point.
	BottomLeft
	// Bottom is the bottom-center anchor point.
	Bottom
	// BottomRight is the bottom-right anchor point.
	BottomRight
)

// ConvolveOptions specifies options for convolution operations.
type ConvolveOptions struct {
	// Normalize indicates whether the kernel should be normalized.
	Normalize bool
	// Abs indicates whether to take the absolute value of results.
	Abs bool
	// Bias is added to each channel value after convolution.
	Bias int
}
