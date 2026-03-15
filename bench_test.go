package pixo

import (
	"bytes"
	"image"
	"image/color"
	"testing"
)

func benchImage(w, h int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.SetNRGBA(x, y, color.NRGBA{
				R: uint8(x % 256),
				G: uint8(y % 256),
				B: uint8((x + y) % 256),
				A: 255,
			})
		}
	}
	return img
}

func BenchmarkResizeNearestNeighbor(b *testing.B) {
	src := benchImage(1920, 1080)
	for range b.N {
		Resize(src, 640, 360, NearestNeighbor)
	}
}

func BenchmarkResizeLinear(b *testing.B) {
	src := benchImage(1920, 1080)
	for range b.N {
		Resize(src, 640, 360, Linear)
	}
}

func BenchmarkResizeCatmullRom(b *testing.B) {
	src := benchImage(1920, 1080)
	for range b.N {
		Resize(src, 640, 360, CatmullRom)
	}
}

func BenchmarkResizeLanczos3(b *testing.B) {
	src := benchImage(1920, 1080)
	for range b.N {
		Resize(src, 640, 360, Lanczos3)
	}
}

func BenchmarkResizeSmall(b *testing.B) {
	src := benchImage(256, 256)
	for range b.N {
		Resize(src, 128, 128, Lanczos3)
	}
}

func BenchmarkResizeLarge(b *testing.B) {
	src := benchImage(4096, 2160)
	for range b.N {
		Resize(src, 1920, 1080, Lanczos3)
	}
}

func BenchmarkBlur(b *testing.B) {
	src := benchImage(640, 480)
	for range b.N {
		Blur(src, 2.0)
	}
}

func BenchmarkSharpen(b *testing.B) {
	src := benchImage(640, 480)
	for range b.N {
		Sharpen(src, 1.0)
	}
}

func BenchmarkGrayscale(b *testing.B) {
	src := benchImage(1920, 1080)
	for range b.N {
		Grayscale(src)
	}
}

func BenchmarkRotate90(b *testing.B) {
	src := benchImage(1920, 1080)
	for range b.N {
		Rotate90(src)
	}
}

func BenchmarkFlipH(b *testing.B) {
	src := benchImage(1920, 1080)
	for range b.N {
		FlipH(src)
	}
}

func BenchmarkEncodePNG(b *testing.B) {
	src := benchImage(640, 480)
	for range b.N {
		var buf bytes.Buffer
		_ = Encode(&buf, src, PNG)
	}
}

func BenchmarkEncodeJPEG(b *testing.B) {
	src := benchImage(640, 480)
	for range b.N {
		var buf bytes.Buffer
		_ = Encode(&buf, src, JPEG, JPEGQuality(85))
	}
}

func BenchmarkDecodePNG(b *testing.B) {
	src := benchImage(640, 480)
	var buf bytes.Buffer
	_ = Encode(&buf, src, PNG)
	data := buf.Bytes()

	for range b.N {
		r := bytes.NewReader(data)
		_, _ = Decode(r)
	}
}

func BenchmarkDecodeJPEG(b *testing.B) {
	src := benchImage(640, 480)
	var buf bytes.Buffer
	_ = Encode(&buf, src, JPEG, JPEGQuality(85))
	data := buf.Bytes()

	for range b.N {
		r := bytes.NewReader(data)
		_, _ = Decode(r)
	}
}
