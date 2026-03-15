package pixo

import (
	"image"
	"image/color"
	"image/draw"
)

// toNRGBA converts any image.Image to *image.NRGBA.
func toNRGBA(img image.Image) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}
	if nrgba, ok := img.(*image.NRGBA); ok {
		return nrgba
	}
	bounds := img.Bounds()
	dst := image.NewNRGBA(bounds)
	draw.Draw(dst, bounds, img, bounds.Min, draw.Src)
	return dst
}

// newNRGBA creates a new NRGBA image filled with the given color.
func newNRGBA(width, height int, c color.Color) *image.NRGBA {
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	if c != nil {
		draw.Draw(dst, dst.Bounds(), &image.Uniform{C: c}, image.Point{}, draw.Src)
	}
	return dst
}
