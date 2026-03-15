package pixo

import (
	"image"
	"math"

	"golang.org/x/image/draw"
)

// filterToInterpolator maps a Filter to a golang.org/x/image/draw interpolator.
func filterToInterpolator(f Filter) draw.Interpolator {
	switch f {
	case NearestNeighbor:
		return draw.NearestNeighbor
	case Box:
		return draw.ApproxBiLinear
	case Linear:
		return draw.BiLinear
	case CatmullRom:
		return draw.CatmullRom
	case MitchellNetravali:
		return draw.CatmullRom // Go's x/image doesn't have MN, CatmullRom is closest
	case Lanczos2:
		return draw.CatmullRom
	case Lanczos3:
		return draw.CatmullRom
	default:
		return draw.CatmullRom
	}
}

// Resize resizes img to the specified width and height using the given filter.
// If width or height is 0, the value is calculated to preserve the aspect ratio.
// If both are 0, the original image is returned as *image.NRGBA.
func Resize(img image.Image, width, height int, filter Filter) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	if srcW == 0 || srcH == 0 {
		return &image.NRGBA{}
	}

	if width == 0 && height == 0 {
		return toNRGBA(img)
	}

	if width == 0 {
		width = int(math.Round(float64(height) * float64(srcW) / float64(srcH)))
		if width == 0 {
			width = 1
		}
	}

	if height == 0 {
		height = int(math.Round(float64(width) * float64(srcH) / float64(srcW)))
		if height == 0 {
			height = 1
		}
	}

	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	interp := filterToInterpolator(filter)
	interp.Scale(dst, dst.Bounds(), img, srcBounds, draw.Over, nil)
	return dst
}

// Fit scales img to fit within the specified dimensions while preserving aspect ratio.
// The resulting image will be at most width x height pixels.
func Fit(img image.Image, width, height int, filter Filter) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	if srcW == 0 || srcH == 0 || width == 0 || height == 0 {
		return &image.NRGBA{}
	}

	if srcW <= width && srcH <= height {
		return toNRGBA(img)
	}

	ratio := math.Min(float64(width)/float64(srcW), float64(height)/float64(srcH))
	newW := int(math.Round(float64(srcW) * ratio))
	newH := int(math.Round(float64(srcH) * ratio))

	if newW == 0 {
		newW = 1
	}
	if newH == 0 {
		newH = 1
	}

	return Resize(img, newW, newH, filter)
}

// Fill scales img to fill the specified dimensions, cropping the excess using the anchor point.
// The resulting image will be exactly width x height pixels.
func Fill(img image.Image, width, height int, filter Filter, anchor Anchor) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	if srcW == 0 || srcH == 0 || width == 0 || height == 0 {
		return &image.NRGBA{}
	}

	ratio := math.Max(float64(width)/float64(srcW), float64(height)/float64(srcH))
	newW := int(math.Round(float64(srcW) * ratio))
	newH := int(math.Round(float64(srcH) * ratio))

	if newW == 0 {
		newW = 1
	}
	if newH == 0 {
		newH = 1
	}

	resized := Resize(img, newW, newH, filter)
	return CropAnchor(resized, width, height, anchor)
}

// Thumbnail scales img to fill the specified dimensions, then crops to exact size.
// It uses Center anchor and is equivalent to Fill with Center anchor.
func Thumbnail(img image.Image, width, height int, filter Filter) *image.NRGBA {
	return Fill(img, width, height, filter, Center)
}
