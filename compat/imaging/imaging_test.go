package imaging

import (
	"bytes"
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
	result := Resize(src, 50, 25, Lanczos)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 25 {
		t.Errorf("got %dx%d, want 50x25", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestFit(t *testing.T) {
	src := testImage(200, 100)
	result := Fit(src, 100, 100, Lanczos)
	if result.Bounds().Dx() != 100 || result.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 100x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestFill(t *testing.T) {
	src := testImage(200, 100)
	result := Fill(src, 100, 100, Center, Lanczos)
	if result.Bounds().Dx() != 100 || result.Bounds().Dy() != 100 {
		t.Errorf("got %dx%d, want 100x100", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestThumbnail(t *testing.T) {
	src := testImage(200, 100)
	result := Thumbnail(src, 50, 50, Lanczos)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 50x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestCropFunctions(t *testing.T) {
	src := testImage(100, 100)

	result := Crop(src, image.Rect(10, 10, 60, 60))
	if result.Bounds().Dx() != 50 {
		t.Errorf("Crop: got width %d, want 50", result.Bounds().Dx())
	}

	result = CropAnchor(src, 50, 50, TopLeft)
	if result.Bounds().Dx() != 50 {
		t.Errorf("CropAnchor: got width %d, want 50", result.Bounds().Dx())
	}

	result = CropCenter(src, 50, 50)
	if result.Bounds().Dx() != 50 {
		t.Errorf("CropCenter: got width %d, want 50", result.Bounds().Dx())
	}
}

func TestTransforms(t *testing.T) {
	src := testImage(100, 50)

	if Rotate90(src).Bounds().Dx() != 50 {
		t.Error("Rotate90 failed")
	}
	if Rotate180(src).Bounds().Dx() != 100 {
		t.Error("Rotate180 failed")
	}
	if Rotate270(src).Bounds().Dx() != 50 {
		t.Error("Rotate270 failed")
	}
	if FlipH(src).Bounds().Dx() != 100 {
		t.Error("FlipH failed")
	}
	if FlipV(src).Bounds().Dx() != 100 {
		t.Error("FlipV failed")
	}
	if Transpose(src).Bounds().Dx() != 50 {
		t.Error("Transpose failed")
	}
	if Transverse(src).Bounds().Dx() != 50 {
		t.Error("Transverse failed")
	}

	result := Rotate(src, 45, color.Black)
	if result.Bounds().Dx() <= 0 {
		t.Error("Rotate failed")
	}
}

func TestAdjustments(t *testing.T) {
	src := testImage(50, 50)

	AdjustBrightness(src, 50)
	AdjustContrast(src, 50)
	AdjustGamma(src, 2.0)
	AdjustSaturation(src, 50)
	AdjustSigmoid(src, 0.5, 5)
}

func TestEffects(t *testing.T) {
	src := testImage(50, 50)

	Blur(src, 1.0)
	Sharpen(src, 1.0)
	Grayscale(src)
	Invert(src)
}

func TestCompose(t *testing.T) {
	bg := testImage(100, 100)
	fg := testImage(50, 50)

	Paste(bg, fg, image.Pt(10, 10))
	PasteCenter(bg, fg)
	Overlay(bg, fg, image.Pt(10, 10), 0.5)
	OverlayCenter(bg, fg, 0.5)
}

func TestHistogram(t *testing.T) {
	src := testImage(10, 10)
	hist := Histogram(src)
	if hist[0][128] != 100 {
		t.Errorf("expected 100 at R=128, got %d", hist[0][128])
	}
}

func TestEncodeDecodePNG(t *testing.T) {
	src := testImage(50, 50)
	var buf bytes.Buffer
	if err := Encode(&buf, src, PNG); err != nil {
		t.Fatal(err)
	}
	img, err := Decode(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != 50 {
		t.Errorf("got width %d, want 50", img.Bounds().Dx())
	}
}

func TestAllFilters(t *testing.T) {
	src := testImage(100, 100)
	filters := []ResampleFilter{NearestNeighbor, Box, Linear, CatmullRom, MitchellNetravali, Lanczos}
	for _, f := range filters {
		result := Resize(src, 50, 50, f)
		if result.Bounds().Dx() != 50 {
			t.Errorf("filter %d: got width %d, want 50", f, result.Bounds().Dx())
		}
	}
}

func TestAllAnchors(t *testing.T) {
	src := testImage(100, 100)
	anchors := []Anchor{Center, TopLeft, Top, TopRight, Left, Right, BottomLeft, Bottom, BottomRight}
	for _, a := range anchors {
		result := Fill(src, 50, 50, a, Lanczos)
		if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 50 {
			t.Errorf("anchor %d: got %dx%d, want 50x50", a, result.Bounds().Dx(), result.Bounds().Dy())
		}
	}
}
