package resize

import (
	"image"
	"image/color"
	"testing"
)

func testImage(w, h int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.SetNRGBA(x, y, color.NRGBA{R: 128, G: 64, B: 32, A: 255})
		}
	}
	return img
}

func TestResize(t *testing.T) {
	src := testImage(100, 50)
	result := Resize(50, 25, src, Lanczos3)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 25 {
		t.Errorf("got %dx%d, want 50x25", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestResizePreserveAspect(t *testing.T) {
	src := testImage(100, 50)
	result := Resize(50, 0, src, Bilinear)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 25 {
		t.Errorf("got %dx%d, want 50x25", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestThumbnail(t *testing.T) {
	src := testImage(200, 100)
	result := Thumbnail(100, 100, src, Lanczos3)
	if result.Bounds().Dx() != 100 || result.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 100x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestAllInterpolations(t *testing.T) {
	src := testImage(100, 100)
	interps := []InterpolationFunction{NearestNeighbor, Bilinear, Bicubic, MitchellNetravali, Lanczos2, Lanczos3}
	for _, interp := range interps {
		result := Resize(50, 50, src, interp)
		if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 50 {
			t.Errorf("interp %d: got %dx%d, want 50x50", interp, result.Bounds().Dx(), result.Bounds().Dy())
		}
	}
}
