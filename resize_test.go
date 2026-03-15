package pixo

import (
	"image"
	"image/color"
	"testing"
)

func testImage(w, h int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.SetNRGBA(x, y, color.NRGBA{
				R: uint8(x * 255 / max(w-1, 1)),
				G: uint8(y * 255 / max(h-1, 1)),
				B: 128,
				A: 255,
			})
		}
	}
	return img
}

func TestResize(t *testing.T) {
	src := testImage(100, 50)

	tests := []struct {
		name   string
		w, h   int
		filter Filter
		wantW  int
		wantH  int
	}{
		{"downscale", 50, 25, Lanczos3, 50, 25},
		{"upscale", 200, 100, Linear, 200, 100},
		{"width only", 50, 0, CatmullRom, 50, 25},
		{"height only", 0, 100, NearestNeighbor, 200, 100},
		{"both zero", 0, 0, Box, 100, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resize(src, tt.w, tt.h, tt.filter)
			if result.Bounds().Dx() != tt.wantW || result.Bounds().Dy() != tt.wantH {
				t.Errorf("got %dx%d, want %dx%d", result.Bounds().Dx(), result.Bounds().Dy(), tt.wantW, tt.wantH)
			}
		})
	}
}

func TestResizeNil(t *testing.T) {
	result := Resize(nil, 100, 100, Lanczos3)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestFit(t *testing.T) {
	src := testImage(200, 100)

	result := Fit(src, 100, 100, Lanczos3)
	if result.Bounds().Dx() != 100 || result.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 100x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestFitNoUpscale(t *testing.T) {
	src := testImage(50, 25)

	result := Fit(src, 100, 100, Lanczos3)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 25 {
		t.Errorf("got %dx%d, want 50x25 (no upscale)", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestFill(t *testing.T) {
	src := testImage(200, 100)

	result := Fill(src, 100, 100, Lanczos3, Center)
	if result.Bounds().Dx() != 100 || result.Bounds().Dy() != 100 {
		t.Errorf("got %dx%d, want 100x100", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestThumbnail(t *testing.T) {
	src := testImage(200, 100)

	result := Thumbnail(src, 50, 50, Lanczos3)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 50x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestAllFilters(t *testing.T) {
	src := testImage(100, 100)
	filters := []Filter{NearestNeighbor, Box, Linear, CatmullRom, MitchellNetravali, Lanczos2, Lanczos3}

	for _, f := range filters {
		result := Resize(src, 50, 50, f)
		if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 50 {
			t.Errorf("filter %d: got %dx%d, want 50x50", f, result.Bounds().Dx(), result.Bounds().Dy())
		}
	}
}
