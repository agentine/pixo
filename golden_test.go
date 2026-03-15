package pixo

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"testing"
)

// generateTestImage creates a gradient test image programmatically.
func generateTestImage(w, h int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.SetNRGBA(x, y, color.NRGBA{
				R: uint8(x * 255 / max(w-1, 1)),
				G: uint8(y * 255 / max(h-1, 1)),
				B: uint8((x + y) * 255 / max(w+h-2, 1)),
				A: 255,
			})
		}
	}
	return img
}

// compareImages compares two images pixel-by-pixel with tolerance.
func compareImages(a, b *image.NRGBA, tolerance int) bool {
	if a.Bounds() != b.Bounds() {
		return false
	}
	bounds := a.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			ca := a.NRGBAAt(x, y)
			cb := b.NRGBAAt(x, y)
			if absDiff(int(ca.R), int(cb.R)) > tolerance ||
				absDiff(int(ca.G), int(cb.G)) > tolerance ||
				absDiff(int(ca.B), int(cb.B)) > tolerance ||
				absDiff(int(ca.A), int(cb.A)) > tolerance {
				return false
			}
		}
	}
	return true
}

func absDiff(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

func saveGolden(t *testing.T, name string, img *image.NRGBA) {
	t.Helper()
	dir := filepath.Join("testdata", "golden")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, name+".png")
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
}

func loadGolden(t *testing.T, name string) *image.NRGBA {
	t.Helper()
	path := filepath.Join("testdata", "golden", name+".png")
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return toNRGBA(img)
}

func goldenTest(t *testing.T, name string, result *image.NRGBA, tolerance int) {
	t.Helper()
	golden := loadGolden(t, name)
	if golden == nil {
		// First run — save golden
		saveGolden(t, name, result)
		t.Logf("saved golden: %s", name)
		return
	}
	if !compareImages(result, golden, tolerance) {
		saveGolden(t, name+"_actual", result)
		t.Errorf("golden mismatch for %s", name)
	}
}

func TestGoldenResize(t *testing.T) {
	src := generateTestImage(100, 100)
	filters := map[string]Filter{
		"nearest":      NearestNeighbor,
		"box":          Box,
		"linear":       Linear,
		"catmullrom":   CatmullRom,
		"mitchell":     MitchellNetravali,
		"lanczos2":     Lanczos2,
		"lanczos3":     Lanczos3,
	}
	for name, f := range filters {
		t.Run(name, func(t *testing.T) {
			result := Resize(src, 50, 50, f)
			goldenTest(t, "resize_"+name, result, 1)
		})
	}
}

func TestGoldenTransforms(t *testing.T) {
	src := generateTestImage(80, 60)

	tests := map[string]*image.NRGBA{
		"rotate90":   Rotate90(src),
		"rotate180":  Rotate180(src),
		"rotate270":  Rotate270(src),
		"fliph":      FlipH(src),
		"flipv":      FlipV(src),
		"transpose":  Transpose(src),
		"transverse": Transverse(src),
	}
	for name, result := range tests {
		t.Run(name, func(t *testing.T) {
			goldenTest(t, "transform_"+name, result, 0)
		})
	}
}

func TestGoldenCrop(t *testing.T) {
	src := generateTestImage(100, 100)

	t.Run("center", func(t *testing.T) {
		result := CropCenter(src, 50, 50)
		goldenTest(t, "crop_center", result, 0)
	})

	t.Run("topleft", func(t *testing.T) {
		result := CropAnchor(src, 50, 50, TopLeft)
		goldenTest(t, "crop_topleft", result, 0)
	})
}

func TestGoldenAdjust(t *testing.T) {
	src := generateTestImage(50, 50)

	tests := map[string]*image.NRGBA{
		"brightness_up":   AdjustBrightness(src, 30),
		"brightness_down": AdjustBrightness(src, -30),
		"contrast_up":     AdjustContrast(src, 30),
		"gamma":           AdjustGamma(src, 1.5),
		"saturation":      AdjustSaturation(src, 50),
		"hue_shift":       AdjustHue(src, 90),
		"sigmoid":         AdjustSigmoid(src, 0.5, 5),
	}
	for name, result := range tests {
		t.Run(name, func(t *testing.T) {
			goldenTest(t, "adjust_"+name, result, 1)
		})
	}
}

func TestGoldenEffects(t *testing.T) {
	src := generateTestImage(50, 50)

	tests := map[string]*image.NRGBA{
		"blur":      Blur(src, 2.0),
		"sharpen":   Sharpen(src, 1.0),
		"grayscale": Grayscale(src),
		"invert":    Invert(src),
	}
	for name, result := range tests {
		t.Run(name, func(t *testing.T) {
			goldenTest(t, "effect_"+name, result, 1)
		})
	}
}

func TestGoldenCompose(t *testing.T) {
	bg := generateTestImage(100, 100)
	fg := generateTestImage(50, 50)

	t.Run("paste", func(t *testing.T) {
		result := Paste(bg, fg, image.Pt(25, 25))
		goldenTest(t, "compose_paste", result, 0)
	})

	t.Run("overlay", func(t *testing.T) {
		result := Overlay(bg, fg, image.Pt(25, 25), 0.5)
		goldenTest(t, "compose_overlay", result, 1)
	})
}

func TestGoldenIORoundTrip(t *testing.T) {
	src := generateTestImage(50, 50)

	t.Run("png", func(t *testing.T) {
		var buf bytes.Buffer
		if err := Encode(&buf, src, PNG); err != nil {
			t.Fatal(err)
		}
		img, err := Decode(&buf)
		if err != nil {
			t.Fatal(err)
		}
		result := toNRGBA(img)
		if !compareImages(src, result, 0) {
			t.Error("PNG round-trip pixel mismatch")
		}
	})

	t.Run("jpeg", func(t *testing.T) {
		var buf bytes.Buffer
		if err := Encode(&buf, src, JPEG, JPEGQuality(100)); err != nil {
			t.Fatal(err)
		}
		img, err := Decode(&buf)
		if err != nil {
			t.Fatal(err)
		}
		result := toNRGBA(img)
		// JPEG is lossy — even at quality 100
		if !compareImages(src, result, 10) {
			t.Error("JPEG round-trip pixel mismatch (tolerance 10)")
		}
	})
}

// Edge case tests

func TestResizeZeroDimensions(t *testing.T) {
	result := Resize(nil, 50, 50, Lanczos3)
	if result.Bounds().Dx() != 0 || result.Bounds().Dy() != 0 {
		t.Errorf("nil input should return empty image")
	}
}

func TestResizeEmptyImage(t *testing.T) {
	empty := image.NewNRGBA(image.Rect(0, 0, 0, 0))
	result := Resize(empty, 50, 50, Lanczos3)
	if result.Bounds().Dx() != 0 || result.Bounds().Dy() != 0 {
		t.Errorf("empty input should return empty image")
	}
}

func TestCropLargerThanImage(t *testing.T) {
	src := testImage(50, 50)
	result := CropCenter(src, 100, 100)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 50 {
		t.Errorf("crop larger than image: got %dx%d, want 50x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestAdjustExtremeValues(t *testing.T) {
	src := testImage(10, 10)
	// Extreme brightness
	result := AdjustBrightness(src, 100)
	for y := range 10 {
		for x := range 10 {
			c := result.NRGBAAt(x, y)
			if c.R != 255 || c.G != 255 || c.B != 255 {
				// With +100% brightness shift of 255, all values should clamp to 255
			}
			if c.A != 255 {
				t.Errorf("alpha should be preserved: got %d", c.A)
			}
		}
	}

	// Full desaturation
	result = AdjustSaturation(src, -100)
	for y := range 10 {
		for x := range 10 {
			c := result.NRGBAAt(x, y)
			// Grayscale means R==G==B (within rounding)
			if absDiff(int(c.R), int(c.G)) > 1 || absDiff(int(c.G), int(c.B)) > 1 {
				t.Errorf("desaturated pixel should be grayscale: R=%d G=%d B=%d", c.R, c.G, c.B)
			}
		}
	}
}

func TestBlurPreservesAlpha(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 10, 10))
	for y := range 10 {
		for x := range 10 {
			src.SetNRGBA(x, y, color.NRGBA{R: 128, G: 128, B: 128, A: 255})
		}
	}
	result := Blur(src, 1.0)
	for y := 2; y < 8; y++ {
		for x := 2; x < 8; x++ {
			c := result.NRGBAAt(x, y)
			if c.A != 255 {
				t.Errorf("blur should preserve alpha for uniform image: got %d at (%d,%d)", c.A, x, y)
			}
		}
	}
}

func TestConvolveIdentity(t *testing.T) {
	src := testImage(20, 20)
	// Identity kernel
	kernel := [9]float64{0, 0, 0, 0, 1, 0, 0, 0, 0}
	result := Convolve3x3(src, kernel, nil)

	// Interior pixels should be identical
	for y := 1; y < 19; y++ {
		for x := 1; x < 19; x++ {
			cs := src.NRGBAAt(x, y)
			cr := result.NRGBAAt(x, y)
			if cs != cr {
				t.Errorf("identity kernel mismatch at (%d,%d): %v != %v", x, y, cs, cr)
			}
		}
	}
}

func TestRotate360(t *testing.T) {
	src := testImage(50, 50)
	result := Rotate(src, 360, color.Black)
	if result.Bounds().Dx() != 50 || result.Bounds().Dy() != 50 {
		t.Errorf("360 rotation: got %dx%d, want 50x50", result.Bounds().Dx(), result.Bounds().Dy())
	}
}

func TestInvertTwice(t *testing.T) {
	src := testImage(20, 20)
	result := Invert(Invert(src))
	if !compareImages(src, result, 0) {
		t.Error("double invert should return original")
	}
}

func TestFlipHTwice(t *testing.T) {
	src := testImage(20, 20)
	result := FlipH(FlipH(src))
	if !compareImages(src, result, 0) {
		t.Error("double FlipH should return original")
	}
}

func TestFlipVTwice(t *testing.T) {
	src := testImage(20, 20)
	result := FlipV(FlipV(src))
	if !compareImages(src, result, 0) {
		t.Error("double FlipV should return original")
	}
}

func TestRotate180Twice(t *testing.T) {
	src := testImage(20, 20)
	result := Rotate180(Rotate180(src))
	if !compareImages(src, result, 0) {
		t.Error("double Rotate180 should return original")
	}
}

func TestGrayscaleIdempotent(t *testing.T) {
	src := testImage(20, 20)
	once := Grayscale(src)
	twice := Grayscale(once)
	if !compareImages(once, twice, 0) {
		t.Error("grayscale should be idempotent")
	}
}

func TestFillAllAnchors(t *testing.T) {
	src := generateTestImage(200, 100)
	anchors := []Anchor{TopLeft, Top, TopRight, Left, Center, Right, BottomLeft, Bottom, BottomRight}
	for _, a := range anchors {
		result := Fill(src, 80, 80, Lanczos3, a)
		if result.Bounds().Dx() != 80 || result.Bounds().Dy() != 80 {
			t.Errorf("anchor %d: got %dx%d, want 80x80", a, result.Bounds().Dx(), result.Bounds().Dy())
		}
	}
}

func TestHistogramConsistency(t *testing.T) {
	src := generateTestImage(100, 100)
	hist := Histogram(src)

	// Total count per channel should equal pixel count
	total := 100 * 100
	for ch := range 4 {
		var sum int
		for i := range 256 {
			sum += hist[ch][i]
		}
		if sum != total {
			t.Errorf("channel %d: sum=%d, want %d", ch, sum, total)
		}
	}
}

func TestResizePreservesAspectRatio(t *testing.T) {
	src := generateTestImage(200, 100)

	// Width only
	result := Resize(src, 100, 0, Lanczos3)
	ratio := float64(result.Bounds().Dx()) / float64(result.Bounds().Dy())
	expected := 2.0
	if math.Abs(ratio-expected) > 0.1 {
		t.Errorf("aspect ratio: got %.2f, want ~%.2f", ratio, expected)
	}

	// Height only
	result = Resize(src, 0, 50, Lanczos3)
	ratio = float64(result.Bounds().Dx()) / float64(result.Bounds().Dy())
	if math.Abs(ratio-expected) > 0.1 {
		t.Errorf("aspect ratio: got %.2f, want ~%.2f", ratio, expected)
	}
}
