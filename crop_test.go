package pixo

import (
	"image"
	"testing"
)

func TestCrop(t *testing.T) {
	src := testImage(100, 100)
	result := Crop(src, image.Rect(10, 10, 60, 60))
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 50x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestCropOutOfBounds(t *testing.T) {
	src := testImage(100, 100)
	result := Crop(src, image.Rect(50, 50, 200, 200))
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 50x50 (clamped)", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestCropAnchor(t *testing.T) {
	src := testImage(100, 100)

	anchors := []Anchor{TopLeft, Top, TopRight, Left, Center, Right, BottomLeft, Bottom, BottomRight}
	for _, a := range anchors {
		result := CropAnchor(src, 50, 50, a)
		if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 50 {
			t.Errorf("anchor %d: got %dx%d, want 50x50", a, result.Bounds().Dx(), result.Bounds().Dy())
		}
	}
}

func TestCropCenter(t *testing.T) {
	src := testImage(100, 100)
	result := CropCenter(src, 50, 50)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 50x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestCropNil(t *testing.T) {
	result := Crop(nil, image.Rect(0, 0, 50, 50))
	if result == nil {
		t.Fatal("expected non-nil")
	}
}
