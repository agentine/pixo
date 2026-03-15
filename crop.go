package pixo

import (
	"image"
	"image/draw"
)

// Crop crops img to the specified rectangle.
func Crop(img image.Image, rect image.Rectangle) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	rect = rect.Intersect(img.Bounds())
	if rect.Empty() {
		return &image.NRGBA{}
	}

	dst := image.NewNRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	draw.Draw(dst, dst.Bounds(), img, rect.Min, draw.Src)
	return dst
}

// CropAnchor crops img to the specified dimensions using the anchor point.
func CropAnchor(img image.Image, width, height int, anchor Anchor) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	if width <= 0 || height <= 0 || srcW == 0 || srcH == 0 {
		return &image.NRGBA{}
	}

	if width > srcW {
		width = srcW
	}
	if height > srcH {
		height = srcH
	}

	var x, y int
	switch anchor {
	case TopLeft:
		x, y = 0, 0
	case Top:
		x, y = (srcW-width)/2, 0
	case TopRight:
		x, y = srcW-width, 0
	case Left:
		x, y = 0, (srcH-height)/2
	case Center:
		x, y = (srcW-width)/2, (srcH-height)/2
	case Right:
		x, y = srcW-width, (srcH-height)/2
	case BottomLeft:
		x, y = 0, srcH-height
	case Bottom:
		x, y = (srcW-width)/2, srcH-height
	case BottomRight:
		x, y = srcW-width, srcH-height
	default:
		x, y = (srcW-width)/2, (srcH-height)/2
	}

	rect := image.Rect(x+srcBounds.Min.X, y+srcBounds.Min.Y, x+srcBounds.Min.X+width, y+srcBounds.Min.Y+height)
	return Crop(img, rect)
}

// CropCenter crops img to the specified dimensions from the center.
func CropCenter(img image.Image, width, height int) *image.NRGBA {
	return CropAnchor(img, width, height, Center)
}
