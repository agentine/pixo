package pixo

import (
	"image"
	"image/color"
	"math"
)

// Rotate rotates img by the given angle (in degrees, counter-clockwise) with the specified background color.
func Rotate(img image.Image, angle float64, bgcolor color.Color) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	angle = math.Mod(angle, 360)
	if angle < 0 {
		angle += 360
	}

	if angle == 0 {
		return toNRGBA(img)
	}
	if angle == 90 {
		return Rotate90(img)
	}
	if angle == 180 {
		return Rotate180(img)
	}
	if angle == 270 {
		return Rotate270(img)
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	rad := angle * math.Pi / 180
	sinA := math.Abs(math.Sin(rad))
	cosA := math.Abs(math.Cos(rad))

	dstW := int(math.Ceil(float64(srcW)*cosA + float64(srcH)*sinA))
	dstH := int(math.Ceil(float64(srcW)*sinA + float64(srcH)*cosA))

	dst := newNRGBA(dstW, dstH, bgcolor)

	srcCX := float64(srcW) / 2
	srcCY := float64(srcH) / 2
	dstCX := float64(dstW) / 2
	dstCY := float64(dstH) / 2

	sinR := math.Sin(-rad)
	cosR := math.Cos(-rad)

	for y := range dstH {
		for x := range dstW {
			dx := float64(x) - dstCX + 0.5
			dy := float64(y) - dstCY + 0.5

			sx := cosR*dx - sinR*dy + srcCX
			sy := sinR*dx + cosR*dy + srcCY

			ix := int(sx)
			iy := int(sy)

			if ix >= 0 && ix < srcW && iy >= 0 && iy < srcH {
				dst.SetNRGBA(x, y, src.NRGBAAt(ix+srcBounds.Min.X, iy+srcBounds.Min.Y))
			}
		}
	}

	return dst
}

// Rotate90 rotates img 90 degrees counter-clockwise.
func Rotate90(img image.Image) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	dst := image.NewNRGBA(image.Rect(0, 0, srcH, srcW))
	for y := range srcW {
		for x := range srcH {
			dst.SetNRGBA(x, y, src.NRGBAAt(srcBounds.Min.X+srcW-1-y, srcBounds.Min.Y+x))
		}
	}
	return dst
}

// Rotate180 rotates img 180 degrees.
func Rotate180(img image.Image) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	dst := image.NewNRGBA(image.Rect(0, 0, srcW, srcH))
	for y := range srcH {
		for x := range srcW {
			dst.SetNRGBA(x, y, src.NRGBAAt(srcBounds.Min.X+srcW-1-x, srcBounds.Min.Y+srcH-1-y))
		}
	}
	return dst
}

// Rotate270 rotates img 270 degrees counter-clockwise (90 degrees clockwise).
func Rotate270(img image.Image) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	dst := image.NewNRGBA(image.Rect(0, 0, srcH, srcW))
	for y := range srcW {
		for x := range srcH {
			dst.SetNRGBA(x, y, src.NRGBAAt(srcBounds.Min.X+y, srcBounds.Min.Y+srcH-1-x))
		}
	}
	return dst
}

// FlipH flips img horizontally.
func FlipH(img image.Image) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	dst := image.NewNRGBA(image.Rect(0, 0, srcW, srcH))
	for y := range srcH {
		for x := range srcW {
			dst.SetNRGBA(x, y, src.NRGBAAt(srcBounds.Min.X+srcW-1-x, srcBounds.Min.Y+y))
		}
	}
	return dst
}

// FlipV flips img vertically.
func FlipV(img image.Image) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	dst := image.NewNRGBA(image.Rect(0, 0, srcW, srcH))
	for y := range srcH {
		for x := range srcW {
			dst.SetNRGBA(x, y, src.NRGBAAt(srcBounds.Min.X+x, srcBounds.Min.Y+srcH-1-y))
		}
	}
	return dst
}

// Transpose flips img along the top-left to bottom-right diagonal.
func Transpose(img image.Image) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	dst := image.NewNRGBA(image.Rect(0, 0, srcH, srcW))
	for y := range srcW {
		for x := range srcH {
			dst.SetNRGBA(x, y, src.NRGBAAt(srcBounds.Min.X+y, srcBounds.Min.Y+x))
		}
	}
	return dst
}

// Transverse flips img along the top-right to bottom-left diagonal.
func Transverse(img image.Image) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	dst := image.NewNRGBA(image.Rect(0, 0, srcH, srcW))
	for y := range srcW {
		for x := range srcH {
			dst.SetNRGBA(x, y, src.NRGBAAt(srcBounds.Min.X+srcW-1-y, srcBounds.Min.Y+srcH-1-x))
		}
	}
	return dst
}
