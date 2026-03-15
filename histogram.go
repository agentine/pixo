package pixo

import (
	"image"
)

// Histogram computes per-channel (R, G, B, A) histogram of img.
// Returns [4][256]int where index 0=R, 1=G, 2=B, 3=A.
func Histogram(img image.Image) [4][256]int {
	var hist [4][256]int
	if img == nil {
		return hist
	}

	src := toNRGBA(img)
	srcBounds := src.Bounds()
	w := srcBounds.Dx()
	h := srcBounds.Dy()

	for y := range h {
		for x := range w {
			c := src.NRGBAAt(srcBounds.Min.X+x, srcBounds.Min.Y+y)
			hist[0][c.R]++
			hist[1][c.G]++
			hist[2][c.B]++
			hist[3][c.A]++
		}
	}
	return hist
}
