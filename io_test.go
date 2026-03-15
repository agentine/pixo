package pixo

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestEncodeDecodeJPEG(t *testing.T) {
	src := testImage(50, 50)
	var buf bytes.Buffer
	if err := Encode(&buf, src, JPEG, JPEGQuality(90)); err != nil {
		t.Fatal(err)
	}
	img, err := Decode(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != 50 || img.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 50x50", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestEncodeDecodePNG(t *testing.T) {
	src := testImage(50, 50)
	var buf bytes.Buffer
	if err := Encode(&buf, src, PNG, PNGCompressionLevel(png.BestSpeed)); err != nil {
		t.Fatal(err)
	}
	img, err := Decode(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != 50 || img.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 50x50", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestEncodeDecodeGIF(t *testing.T) {
	src := testImage(50, 50)
	var buf bytes.Buffer
	if err := Encode(&buf, src, GIF); err != nil {
		t.Fatal(err)
	}
	img, err := Decode(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != 50 || img.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 50x50", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestEncodeDecodeBMP(t *testing.T) {
	src := testImage(50, 50)
	var buf bytes.Buffer
	if err := Encode(&buf, src, BMP); err != nil {
		t.Fatal(err)
	}
	img, err := Decode(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != 50 || img.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 50x50", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestEncodeDecodeTIFF(t *testing.T) {
	src := testImage(50, 50)
	var buf bytes.Buffer
	if err := Encode(&buf, src, TIFF); err != nil {
		t.Fatal(err)
	}
	img, err := Decode(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != 50 || img.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 50x50", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestSaveOpen(t *testing.T) {
	src := testImage(50, 50)
	dir := t.TempDir()
	path := filepath.Join(dir, "test.png")

	if err := Save(src, path); err != nil {
		t.Fatal(err)
	}
	img, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != 50 || img.Bounds().Dy() != 50 {
		t.Errorf("got %dx%d, want 50x50", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestSaveJPEG(t *testing.T) {
	src := testImage(50, 50)
	dir := t.TempDir()
	path := filepath.Join(dir, "test.jpg")
	if err := Save(src, path, JPEGQuality(85)); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}

func TestFormatString(t *testing.T) {
	tests := []struct {
		f    Format
		want string
	}{
		{JPEG, "jpeg"},
		{PNG, "png"},
		{GIF, "gif"},
		{BMP, "bmp"},
		{TIFF, "tiff"},
		{WebP, "webp"},
		{Format(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.f.String(); got != tt.want {
			t.Errorf("Format(%d).String() = %q, want %q", tt.f, got, tt.want)
		}
	}
}

func TestEncodeUnsupportedFormat(t *testing.T) {
	src := testImage(10, 10)
	var buf bytes.Buffer
	err := Encode(&buf, src, WebP)
	if err == nil {
		t.Error("expected error for WebP encoding")
	}
}

func TestSaveUnsupportedExt(t *testing.T) {
	src := testImage(10, 10)
	dir := t.TempDir()
	err := Save(src, filepath.Join(dir, "test.xyz"))
	if err == nil {
		t.Error("expected error for unsupported extension")
	}
}

func TestOpenNonexistent(t *testing.T) {
	_, err := Open("/nonexistent/path/img.png")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestReadJPEGOrientation(t *testing.T) {
	// Test with a non-JPEG reader - should return 1
	var buf bytes.Buffer
	_ = png.Encode(&buf, image.NewNRGBA(image.Rect(0, 0, 1, 1)))
	orientation := readJPEGOrientation(&buf)
	if orientation != 1 {
		t.Errorf("expected orientation 1 for non-JPEG, got %d", orientation)
	}
}

func TestApplyOrientation(t *testing.T) {
	src := testImage(100, 50)
	for _, o := range []int{1, 2, 3, 4, 5, 6, 7, 8, 0, 99} {
		result := applyOrientation(src, o)
		if result == nil {
			t.Errorf("orientation %d: nil result", o)
		}
	}
}

func TestDecodeReaderNotSeekable(t *testing.T) {
	// Create a PNG in a buffer
	img := image.NewNRGBA(image.Rect(0, 0, 10, 10))
	for y := range 10 {
		for x := range 10 {
			img.SetNRGBA(x, y, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)

	// Wrap in a non-seekable reader
	decoded, err := Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Bounds().Dx() != 10 {
		t.Errorf("got width %d, want 10", decoded.Bounds().Dx())
	}
}
