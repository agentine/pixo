package pixo

import (
	"image"
	"image/color"
	"image/draw"
)

// Paste pastes img onto background at the given position.
func Paste(background image.Image, img image.Image, pos image.Point) *image.NRGBA {
	if background == nil {
		return &image.NRGBA{}
	}

	dst := toNRGBA(background)
	// Make a copy so we don't modify the original
	bgBounds := dst.Bounds()
	result := image.NewNRGBA(bgBounds)
	draw.Draw(result, bgBounds, dst, bgBounds.Min, draw.Src)

	if img != nil {
		rect := image.Rect(pos.X, pos.Y, pos.X+img.Bounds().Dx(), pos.Y+img.Bounds().Dy())
		draw.Draw(result, rect, img, img.Bounds().Min, draw.Over)
	}
	return result
}

// PasteCenter pastes img onto the center of background.
func PasteCenter(background image.Image, img image.Image) *image.NRGBA {
	if background == nil {
		return &image.NRGBA{}
	}
	if img == nil {
		return toNRGBA(background)
	}

	bgBounds := background.Bounds()
	imgBounds := img.Bounds()
	x := (bgBounds.Dx() - imgBounds.Dx()) / 2
	y := (bgBounds.Dy() - imgBounds.Dy()) / 2
	return Paste(background, img, image.Pt(x+bgBounds.Min.X, y+bgBounds.Min.Y))
}

// Overlay draws img over background at pos with the given opacity [0, 1].
func Overlay(background image.Image, img image.Image, pos image.Point, opacity float64) *image.NRGBA {
	if background == nil {
		return &image.NRGBA{}
	}

	dst := toNRGBA(background)
	bgBounds := dst.Bounds()
	result := image.NewNRGBA(bgBounds)
	draw.Draw(result, bgBounds, dst, bgBounds.Min, draw.Src)

	if img == nil || opacity <= 0 {
		return result
	}

	if opacity > 1 {
		opacity = 1
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()

	for y := range srcBounds.Dy() {
		for x := range srcBounds.Dx() {
			dx := pos.X + x
			dy := pos.Y + y
			if dx < bgBounds.Min.X || dx >= bgBounds.Max.X || dy < bgBounds.Min.Y || dy >= bgBounds.Max.Y {
				continue
			}

			fg := src.NRGBAAt(srcBounds.Min.X+x, srcBounds.Min.Y+y)
			bg := result.NRGBAAt(dx, dy)

			fgA := float64(fg.A) / 255 * opacity
			bgA := float64(bg.A) / 255

			outA := fgA + bgA*(1-fgA)
			if outA == 0 {
				continue
			}

			result.SetNRGBA(dx, dy, color.NRGBA{
				R: clampUint8((float64(fg.R)*fgA + float64(bg.R)*bgA*(1-fgA)) / outA),
				G: clampUint8((float64(fg.G)*fgA + float64(bg.G)*bgA*(1-fgA)) / outA),
				B: clampUint8((float64(fg.B)*fgA + float64(bg.B)*bgA*(1-fgA)) / outA),
				A: clampUint8(outA * 255),
			})
		}
	}
	return result
}

// OverlayCenter draws img over the center of background with the given opacity.
func OverlayCenter(background image.Image, img image.Image, opacity float64) *image.NRGBA {
	if background == nil {
		return &image.NRGBA{}
	}
	if img == nil {
		return toNRGBA(background)
	}

	bgBounds := background.Bounds()
	imgBounds := img.Bounds()
	x := (bgBounds.Dx() - imgBounds.Dx()) / 2
	y := (bgBounds.Dy() - imgBounds.Dy()) / 2
	return Overlay(background, img, image.Pt(x+bgBounds.Min.X, y+bgBounds.Min.Y), opacity)
}
