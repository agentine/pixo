// Package resize provides a drop-in replacement for github.com/nfnt/resize.
//
// To migrate, change your import path from:
//
//	import "github.com/nfnt/resize"
//
// to:
//
//	import "github.com/agentine/pixo/compat/resize"
package resize

import (
	"image"

	"github.com/agentine/pixo"
)

// InterpolationFunction is the type for interpolation algorithms.
type InterpolationFunction int

const (
	// NearestNeighbor interpolation.
	NearestNeighbor InterpolationFunction = iota
	// Bilinear interpolation.
	Bilinear
	// Bicubic interpolation.
	Bicubic
	// MitchellNetravali interpolation.
	MitchellNetravali
	// Lanczos2 interpolation.
	Lanczos2
	// Lanczos3 interpolation.
	Lanczos3
)

func toPixoFilter(interp InterpolationFunction) pixo.Filter {
	switch interp {
	case NearestNeighbor:
		return pixo.NearestNeighbor
	case Bilinear:
		return pixo.Linear
	case Bicubic:
		return pixo.CatmullRom
	case MitchellNetravali:
		return pixo.MitchellNetravali
	case Lanczos2:
		return pixo.Lanczos2
	case Lanczos3:
		return pixo.Lanczos3
	default:
		return pixo.Lanczos3
	}
}

// Resize scales an image to the given dimensions using the specified interpolation.
// Width or height of 0 means preserving aspect ratio.
func Resize(width, height uint, img image.Image, interp InterpolationFunction) image.Image {
	return pixo.Resize(img, int(width), int(height), toPixoFilter(interp))
}

// Thumbnail downscales an image preserving aspect ratio to fit within maxWidth x maxHeight.
func Thumbnail(maxWidth, maxHeight uint, img image.Image, interp InterpolationFunction) image.Image {
	return pixo.Fit(img, int(maxWidth), int(maxHeight), toPixoFilter(interp))
}
