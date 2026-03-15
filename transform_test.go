package pixo

import (
	"image/color"
	"testing"
)

func TestRotate90(t *testing.T) {
	src := testImage(100, 50)
	result := Rotate90(src)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 100 {
		t.Errorf("got %dx%d, want 50x100", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestRotate180(t *testing.T) {
	src := testImage(100, 50)
	result := Rotate180(src)
	if result.Bounds().Dx() != 100 || result.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 100x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestRotate270(t *testing.T) {
	src := testImage(100, 50)
	result := Rotate270(src)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 100 {
		t.Errorf("got %dx%d, want 50x100", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestRotateArbitrary(t *testing.T) {
	src := testImage(100, 100)
	result := Rotate(src, 45, color.Black)
	// Rotated image should be larger
	if result.Bounds().Dx() <= 100 || result.Bounds().Dy() <= 100 {
		t.Errorf("expected dimensions > 100, got %dx%d", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestRotate0(t *testing.T) {
	src := testImage(100, 50)
	result := Rotate(src, 0, color.Black)
	if result.Bounds().Dx() != 100 || result.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 100x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestFlipH(t *testing.T) {
	src := testImage(100, 50)
	result := FlipH(src)
	if result.Bounds().Dx() != 100 || result.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 100x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
	// Check that pixel at (0,0) is now at (99,0)
	srcPx := src.NRGBAAt(0, 0)
	dstPx := result.NRGBAAt(99, 0)
	if srcPx != dstPx {
		t.Errorf("FlipH pixel mismatch: src(0,0)=%v, dst(99,0)=%v", srcPx, dstPx)
	}
}

func TestFlipV(t *testing.T) {
	src := testImage(100, 50)
	result := FlipV(src)
	if result.Bounds().Dx() != 100 || result.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 100x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
	srcPx := src.NRGBAAt(0, 0)
	dstPx := result.NRGBAAt(0, 49)
	if srcPx != dstPx {
		t.Errorf("FlipV pixel mismatch: src(0,0)=%v, dst(0,49)=%v", srcPx, dstPx)
	}
}

func TestTranspose(t *testing.T) {
	src := testImage(100, 50)
	result := Transpose(src)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 100 {
		t.Errorf("got %dx%d, want 50x100", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestTransverse(t *testing.T) {
	src := testImage(100, 50)
	result := Transverse(src)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 100 {
		t.Errorf("got %dx%d, want 50x100", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestTransformNil(t *testing.T) {
	fns := []struct {
		name string
		fn   func()
	}{
		{"Rotate90", func() { Rotate90(nil) }},
		{"Rotate180", func() { Rotate180(nil) }},
		{"Rotate270", func() { Rotate270(nil) }},
		{"FlipH", func() { FlipH(nil) }},
		{"FlipV", func() { FlipV(nil) }},
		{"Transpose", func() { Transpose(nil) }},
		{"Transverse", func() { Transverse(nil) }},
		{"Rotate", func() { Rotate(nil, 45, color.Black) }},
	}
	for _, tt := range fns {
		t.Run(tt.name, func(t *testing.T) {
			tt.fn() // should not panic
		})
	}
}
