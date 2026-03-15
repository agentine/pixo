package pixo

import (
	"image"
	"testing"
)

func TestPaste(t *testing.T) {
	bg := testImage(100, 100)
	fg := testImage(50, 50)
	result := Paste(bg, fg, image.Pt(10, 10))
	if result.Bounds().Dx() != 100 || result.Bounds().Dy() != 100 {
		t.Errorf("got %dx%d, want 100x100", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestPasteCenter(t *testing.T) {
	bg := testImage(100, 100)
	fg := testImage(50, 50)
	result := PasteCenter(bg, fg)
	if result.Bounds().Dx() != 100 || result.Bounds().Dy() != 100 {
		t.Errorf("got %dx%d, want 100x100", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestOverlay(t *testing.T) {
	bg := testImage(100, 100)
	fg := testImage(50, 50)
	result := Overlay(bg, fg, image.Pt(10, 10), 0.5)
	if result.Bounds().Dx() != 100 || result.Bounds().Dy() != 100 {
		t.Errorf("got %dx%d, want 100x100", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestOverlayCenter(t *testing.T) {
	bg := testImage(100, 100)
	fg := testImage(50, 50)
	result := OverlayCenter(bg, fg, 0.5)
	if result.Bounds().Dx() != 100 || result.Bounds().Dy() != 100 {
		t.Errorf("got %dx%d, want 100x100", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestComposeNil(t *testing.T) {
	fg := testImage(50, 50)
	Paste(nil, fg, image.Pt(0, 0))
	PasteCenter(nil, fg)
	Overlay(nil, fg, image.Pt(0, 0), 0.5)
	OverlayCenter(nil, fg, 0.5)
}
