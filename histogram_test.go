package pixo

import (
	"image"
	"image/color"
	"testing"
)

func TestHistogram(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 10, 10))
	for y := range 10 {
		for x := range 10 {
			img.SetNRGBA(x, y, color.NRGBA{R: 128, G: 64, B: 32, A: 255})
		}
	}
	hist := Histogram(img)
	if hist[0][128] != 100 {
		t.Errorf("R channel: expected 100 at 128, got %d", hist[0][128])
	}
	if hist[1][64] != 100 {
		t.Errorf("G channel: expected 100 at 64, got %d", hist[1][64])
	}
	if hist[2][32] != 100 {
		t.Errorf("B channel: expected 100 at 32, got %d", hist[2][32])
	}
	if hist[3][255] != 100 {
		t.Errorf("A channel: expected 100 at 255, got %d", hist[3][255])
	}
}

func TestHistogramNil(t *testing.T) {
	hist := Histogram(nil)
	for ch := range 4 {
		for i := range 256 {
			if hist[ch][i] != 0 {
				t.Fatalf("expected all zeros for nil image")
			}
		}
	}
}
