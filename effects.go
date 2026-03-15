package pixo

import (
	"image"
	"image/color"
	"math"
)

// Blur applies a Gaussian blur to img with the given sigma.
func Blur(img image.Image, sigma float64) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}
	if sigma <= 0 {
		return toNRGBA(img)
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	w := srcBounds.Dx()
	h := srcBounds.Dy()

	// Build 1D Gaussian kernel
	radius := int(math.Ceil(sigma * 3))
	if radius < 1 {
		radius = 1
	}
	size := 2*radius + 1
	kernel := make([]float64, size)
	var sum float64
	for i := range size {
		x := float64(i - radius)
		kernel[i] = math.Exp(-(x * x) / (2 * sigma * sigma))
		sum += kernel[i]
	}
	for i := range kernel {
		kernel[i] /= sum
	}

	// Horizontal pass
	tmp := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			var r, g, b, a float64
			for k := range size {
				sx := x + k - radius
				if sx < 0 {
					sx = 0
				} else if sx >= w {
					sx = w - 1
				}
				c := src.NRGBAAt(srcBounds.Min.X+sx, srcBounds.Min.Y+y)
				r += float64(c.R) * kernel[k]
				g += float64(c.G) * kernel[k]
				b += float64(c.B) * kernel[k]
				a += float64(c.A) * kernel[k]
			}
			tmp.SetNRGBA(x, y, color.NRGBA{
				R: clampUint8(r), G: clampUint8(g), B: clampUint8(b), A: clampUint8(a),
			})
		}
	}

	// Vertical pass
	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			var r, g, b, a float64
			for k := range size {
				sy := y + k - radius
				if sy < 0 {
					sy = 0
				} else if sy >= h {
					sy = h - 1
				}
				c := tmp.NRGBAAt(x, sy)
				r += float64(c.R) * kernel[k]
				g += float64(c.G) * kernel[k]
				b += float64(c.B) * kernel[k]
				a += float64(c.A) * kernel[k]
			}
			dst.SetNRGBA(x, y, color.NRGBA{
				R: clampUint8(r), G: clampUint8(g), B: clampUint8(b), A: clampUint8(a),
			})
		}
	}
	return dst
}

// Sharpen sharpens img using unsharp masking with the given sigma.
func Sharpen(img image.Image, sigma float64) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}
	if sigma <= 0 {
		return toNRGBA(img)
	}

	src := toNRGBA(img)
	blurred := Blur(img, sigma)

	srcBounds := src.Bounds()
	w := srcBounds.Dx()
	h := srcBounds.Dy()

	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			cs := src.NRGBAAt(srcBounds.Min.X+x, srcBounds.Min.Y+y)
			cb := blurred.NRGBAAt(x, y)
			dst.SetNRGBA(x, y, color.NRGBA{
				R: clampUint8(float64(cs.R)*2 - float64(cb.R)),
				G: clampUint8(float64(cs.G)*2 - float64(cb.G)),
				B: clampUint8(float64(cs.B)*2 - float64(cb.B)),
				A: cs.A,
			})
		}
	}
	return dst
}

// Grayscale converts img to grayscale.
func Grayscale(img image.Image) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	w := srcBounds.Dx()
	h := srcBounds.Dy()

	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			c := src.NRGBAAt(srcBounds.Min.X+x, srcBounds.Min.Y+y)
			gray := clampUint8(0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B))
			dst.SetNRGBA(x, y, color.NRGBA{R: gray, G: gray, B: gray, A: c.A})
		}
	}
	return dst
}

// Invert inverts the colors of img.
func Invert(img image.Image) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	w := srcBounds.Dx()
	h := srcBounds.Dy()

	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			c := src.NRGBAAt(srcBounds.Min.X+x, srcBounds.Min.Y+y)
			dst.SetNRGBA(x, y, color.NRGBA{
				R: 255 - c.R, G: 255 - c.G, B: 255 - c.B, A: c.A,
			})
		}
	}
	return dst
}

// Convolve3x3 applies a 3x3 convolution kernel to img.
func Convolve3x3(img image.Image, kernel [9]float64, options *ConvolveOptions) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	w := srcBounds.Dx()
	h := srcBounds.Dy()

	k := kernel
	if options != nil && options.Normalize {
		var sum float64
		for _, v := range k {
			sum += v
		}
		if sum != 0 {
			for i := range k {
				k[i] /= sum
			}
		}
	}

	bias := 0
	abs := false
	if options != nil {
		bias = options.Bias
		abs = options.Abs
	}

	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			var r, g, b float64
			idx := 0
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					sx := x + kx
					sy := y + ky
					if sx < 0 {
						sx = 0
					} else if sx >= w {
						sx = w - 1
					}
					if sy < 0 {
						sy = 0
					} else if sy >= h {
						sy = h - 1
					}
					c := src.NRGBAAt(srcBounds.Min.X+sx, srcBounds.Min.Y+sy)
					r += float64(c.R) * k[idx]
					g += float64(c.G) * k[idx]
					b += float64(c.B) * k[idx]
					idx++
				}
			}
			if abs {
				r = math.Abs(r)
				g = math.Abs(g)
				b = math.Abs(b)
			}
			orig := src.NRGBAAt(srcBounds.Min.X+x, srcBounds.Min.Y+y)
			dst.SetNRGBA(x, y, color.NRGBA{
				R: clampUint8(r + float64(bias)),
				G: clampUint8(g + float64(bias)),
				B: clampUint8(b + float64(bias)),
				A: orig.A,
			})
		}
	}
	return dst
}

// Convolve5x5 applies a 5x5 convolution kernel to img.
func Convolve5x5(img image.Image, kernel [25]float64, options *ConvolveOptions) *image.NRGBA {
	if img == nil {
		return &image.NRGBA{}
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	w := srcBounds.Dx()
	h := srcBounds.Dy()

	k := kernel
	if options != nil && options.Normalize {
		var sum float64
		for _, v := range k {
			sum += v
		}
		if sum != 0 {
			for i := range k {
				k[i] /= sum
			}
		}
	}

	bias := 0
	abs := false
	if options != nil {
		bias = options.Bias
		abs = options.Abs
	}

	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			var r, g, b float64
			idx := 0
			for ky := -2; ky <= 2; ky++ {
				for kx := -2; kx <= 2; kx++ {
					sx := x + kx
					sy := y + ky
					if sx < 0 {
						sx = 0
					} else if sx >= w {
						sx = w - 1
					}
					if sy < 0 {
						sy = 0
					} else if sy >= h {
						sy = h - 1
					}
					c := src.NRGBAAt(srcBounds.Min.X+sx, srcBounds.Min.Y+sy)
					r += float64(c.R) * k[idx]
					g += float64(c.G) * k[idx]
					b += float64(c.B) * k[idx]
					idx++
				}
			}
			if abs {
				r = math.Abs(r)
				g = math.Abs(g)
				b = math.Abs(b)
			}
			orig := src.NRGBAAt(srcBounds.Min.X+x, srcBounds.Min.Y+y)
			dst.SetNRGBA(x, y, color.NRGBA{
				R: clampUint8(r + float64(bias)),
				G: clampUint8(g + float64(bias)),
				B: clampUint8(b + float64(bias)),
				A: orig.A,
			})
		}
	}
	return dst
}
