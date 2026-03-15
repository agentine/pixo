package pixo

import (
	"image"
	"image/color"
	"math"
)

// clampUint8 clamps a float64 to [0, 255] and returns uint8.
func clampUint8(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(math.Round(v))
}

// AdjustBrightness adjusts the brightness of img by the given percentage [-100, 100].
func AdjustBrightness(img image.Image, percentage float64) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	w := srcBounds.Dx()
	h := srcBounds.Dy()

	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	shift := 255 * percentage / 100

	for y := range h {
		for x := range w {
			c := src.NRGBAAt(srcBounds.Min.X+x, srcBounds.Min.Y+y)
			dst.SetNRGBA(x, y, color.NRGBA{
				R: clampUint8(float64(c.R) + shift),
				G: clampUint8(float64(c.G) + shift),
				B: clampUint8(float64(c.B) + shift),
				A: c.A,
			})
		}
	}
	return dst
}

// AdjustContrast adjusts the contrast of img by the given percentage [-100, 100].
func AdjustContrast(img image.Image, percentage float64) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	w := srcBounds.Dx()
	h := srcBounds.Dy()

	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	factor := (100 + percentage) / 100
	factor = factor * factor

	for y := range h {
		for x := range w {
			c := src.NRGBAAt(srcBounds.Min.X+x, srcBounds.Min.Y+y)
			dst.SetNRGBA(x, y, color.NRGBA{
				R: clampUint8((float64(c.R)/255-0.5)*factor*255 + 128),
				G: clampUint8((float64(c.G)/255-0.5)*factor*255 + 128),
				B: clampUint8((float64(c.B)/255-0.5)*factor*255 + 128),
				A: c.A,
			})
		}
	}
	return dst
}

// AdjustGamma adjusts the gamma of img. Gamma must be > 0.
// Gamma < 1 brightens, gamma > 1 darkens.
func AdjustGamma(img image.Image, gamma float64) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}
	if gamma <= 0 {
		return toNRGBA(img)
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	w := srcBounds.Dx()
	h := srcBounds.Dy()

	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	invGamma := 1.0 / gamma

	// Precompute lookup table
	var lut [256]uint8
	for i := range 256 {
		lut[i] = clampUint8(math.Pow(float64(i)/255, invGamma) * 255)
	}

	for y := range h {
		for x := range w {
			c := src.NRGBAAt(srcBounds.Min.X+x, srcBounds.Min.Y+y)
			dst.SetNRGBA(x, y, color.NRGBA{
				R: lut[c.R],
				G: lut[c.G],
				B: lut[c.B],
				A: c.A,
			})
		}
	}
	return dst
}

// rgbToHSL converts RGB to HSL. h in [0,360), s,l in [0,1].
func rgbToHSL(r, g, b uint8) (h, s, l float64) {
	rf := float64(r) / 255
	gf := float64(g) / 255
	bf := float64(b) / 255

	maxC := math.Max(rf, math.Max(gf, bf))
	minC := math.Min(rf, math.Min(gf, bf))
	l = (maxC + minC) / 2

	if maxC == minC {
		return 0, 0, l
	}

	d := maxC - minC
	if l > 0.5 {
		s = d / (2 - maxC - minC)
	} else {
		s = d / (maxC + minC)
	}

	switch maxC {
	case rf:
		h = (gf - bf) / d
		if gf < bf {
			h += 6
		}
	case gf:
		h = (bf-rf)/d + 2
	case bf:
		h = (rf-gf)/d + 4
	}
	h *= 60
	return
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2 {
		return q
	}
	if t < 2.0/3 {
		return p + (q-p)*(2.0/3-t)*6
	}
	return p
}

// hslToRGB converts HSL to RGB. h in [0,360), s,l in [0,1].
func hslToRGB(h, s, l float64) (r, g, b uint8) {
	if s == 0 {
		v := clampUint8(l * 255)
		return v, v, v
	}
	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q
	hNorm := h / 360
	r = clampUint8(hueToRGB(p, q, hNorm+1.0/3) * 255)
	g = clampUint8(hueToRGB(p, q, hNorm) * 255)
	b = clampUint8(hueToRGB(p, q, hNorm-1.0/3) * 255)
	return
}

// AdjustSaturation adjusts the saturation of img by the given percentage [-100, 100].
func AdjustSaturation(img image.Image, percentage float64) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	w := srcBounds.Dx()
	h := srcBounds.Dy()

	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	factor := 1 + percentage/100

	for y := range h {
		for x := range w {
			c := src.NRGBAAt(srcBounds.Min.X+x, srcBounds.Min.Y+y)
			hue, sat, lum := rgbToHSL(c.R, c.G, c.B)
			sat = math.Max(0, math.Min(1, sat*factor))
			nr, ng, nb := hslToRGB(hue, sat, lum)
			dst.SetNRGBA(x, y, color.NRGBA{R: nr, G: ng, B: nb, A: c.A})
		}
	}
	return dst
}

// AdjustSigmoid adjusts the contrast of img using a sigmoidal function.
// Midpoint is the center of the sigmoid (typically 0.5), factor controls steepness.
func AdjustSigmoid(img image.Image, midpoint, factor float64) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	w := srcBounds.Dx()
	h := srcBounds.Dy()

	dst := image.NewNRGBA(image.Rect(0, 0, w, h))

	sigmoid := func(x float64) float64 {
		return 1.0 / (1.0 + math.Exp(-factor*(x-midpoint)))
	}
	s0 := sigmoid(0)
	s1 := sigmoid(1)

	// Precompute lookup table
	var lut [256]uint8
	for i := range 256 {
		v := float64(i) / 255
		v = (sigmoid(v) - s0) / (s1 - s0)
		lut[i] = clampUint8(v * 255)
	}

	for y := range h {
		for x := range w {
			c := src.NRGBAAt(srcBounds.Min.X+x, srcBounds.Min.Y+y)
			dst.SetNRGBA(x, y, color.NRGBA{
				R: lut[c.R],
				G: lut[c.G],
				B: lut[c.B],
				A: c.A,
			})
		}
	}
	return dst
}

// AdjustHue shifts the hue of img by the given degrees [-180, 180].
func AdjustHue(img image.Image, shift int) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	w := srcBounds.Dx()
	h := srcBounds.Dy()

	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	shiftF := float64(shift)

	for y := range h {
		for x := range w {
			c := src.NRGBAAt(srcBounds.Min.X+x, srcBounds.Min.Y+y)
			hue, sat, lum := rgbToHSL(c.R, c.G, c.B)
			hue = math.Mod(hue+shiftF, 360)
			if hue < 0 {
				hue += 360
			}
			nr, ng, nb := hslToRGB(hue, sat, lum)
			dst.SetNRGBA(x, y, color.NRGBA{R: nr, G: ng, B: nb, A: c.A})
		}
	}
	return dst
}
