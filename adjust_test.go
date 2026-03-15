package pixo

import (
	"testing"
)

func TestAdjustBrightness(t *testing.T) {
	src := testImage(50, 50)
	result := AdjustBrightness(src, 50)
	if result.Bounds().Dx() != 50 {
		t.Errorf("unexpected dimensions")
	}
	// Brighter image should have higher pixel values
	orig := src.NRGBAAt(25, 25)
	mod := result.NRGBAAt(25, 25)
	if mod.R <= orig.R && orig.R < 128 {
		t.Errorf("expected brighter: orig R=%d, mod R=%d", orig.R, mod.R)
	}
}

func TestAdjustContrast(t *testing.T) {
	src := testImage(50, 50)
	result := AdjustContrast(src, 50)
	if result.Bounds().Dx() != 50 {
		t.Errorf("unexpected dimensions")
	}
}

func TestAdjustGamma(t *testing.T) {
	src := testImage(50, 50)
	result := AdjustGamma(src, 2.0)
	if result.Bounds().Dx() != 50 {
		t.Errorf("unexpected dimensions")
	}
}

func TestAdjustGammaInvalid(t *testing.T) {
	src := testImage(50, 50)
	result := AdjustGamma(src, 0)
	if result.Bounds().Dx() != 50 {
		t.Errorf("unexpected dimensions")
	}
}

func TestAdjustSaturation(t *testing.T) {
	src := testImage(50, 50)
	result := AdjustSaturation(src, -50)
	if result.Bounds().Dx() != 50 {
		t.Errorf("unexpected dimensions")
	}
}

func TestAdjustSigmoid(t *testing.T) {
	src := testImage(50, 50)
	result := AdjustSigmoid(src, 0.5, 5)
	if result.Bounds().Dx() != 50 {
		t.Errorf("unexpected dimensions")
	}
}

func TestAdjustHue(t *testing.T) {
	src := testImage(50, 50)
	result := AdjustHue(src, 90)
	if result.Bounds().Dx() != 50 {
		t.Errorf("unexpected dimensions")
	}
}

func TestAdjustNil(t *testing.T) {
	AdjustBrightness(nil, 50)
	AdjustContrast(nil, 50)
	AdjustGamma(nil, 2.0)
	AdjustSaturation(nil, 50)
	AdjustSigmoid(nil, 0.5, 5)
	AdjustHue(nil, 90)
}
