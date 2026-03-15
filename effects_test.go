package pixo

import (
	"testing"
)

func TestBlur(t *testing.T) {
	src := testImage(50, 50)
	result := Blur(src, 2.0)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 50x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestBlurZeroSigma(t *testing.T) {
	src := testImage(50, 50)
	result := Blur(src, 0)
	if result.Bounds().Dx() != 50 {
		t.Errorf("unexpected dimensions")
	}
}

func TestSharpen(t *testing.T) {
	src := testImage(50, 50)
	result := Sharpen(src, 1.0)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 50x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestGrayscale(t *testing.T) {
	src := testImage(50, 50)
	result := Grayscale(src)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 50x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
	// All channels should be equal in grayscale
	c := result.NRGBAAt(25, 25)
	if c.R != c.G || c.G != c.B {
		t.Errorf("expected equal channels: R=%d G=%d B=%d", c.R, c.G, c.B)
	}
}

func TestInvert(t *testing.T) {
	src := testImage(50, 50)
	result := Invert(src)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 50x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
	orig := src.NRGBAAt(0, 0)
	inv := result.NRGBAAt(0, 0)
	if inv.R != 255-orig.R || inv.G != 255-orig.G || inv.B != 255-orig.B {
		t.Errorf("invert mismatch: orig=%v, inv=%v", orig, inv)
	}
}

func TestConvolve3x3(t *testing.T) {
	src := testImage(50, 50)
	// Identity kernel
	kernel := [9]float64{0, 0, 0, 0, 1, 0, 0, 0, 0}
	result := Convolve3x3(src, kernel, nil)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 50x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestConvolve5x5(t *testing.T) {
	src := testImage(50, 50)
	// Identity kernel
	var kernel [25]float64
	kernel[12] = 1
	result := Convolve5x5(src, kernel, nil)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 50x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestConvolveWithOptions(t *testing.T) {
	src := testImage(50, 50)
	kernel := [9]float64{1, 1, 1, 1, 1, 1, 1, 1, 1}
	opts := &ConvolveOptions{Normalize: true, Abs: false, Bias: 0}
	result := Convolve3x3(src, kernel, opts)
	if result.Bounds().Dx() != 50 {
		t.Errorf("unexpected dimensions")
	}
}

func TestEffectsNil(t *testing.T) {
	Blur(nil, 1.0)
	Sharpen(nil, 1.0)
	Grayscale(nil)
	Invert(nil)
	Convolve3x3(nil, [9]float64{}, nil)
	Convolve5x5(nil, [25]float64{}, nil)
}
