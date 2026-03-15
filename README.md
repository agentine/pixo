# pixo

[![Go Reference](https://pkg.go.dev/badge/github.com/agentine/pixo.svg)](https://pkg.go.dev/github.com/agentine/pixo)

Pure Go image processing library. Drop-in replacement for [nfnt/resize](https://github.com/nfnt/resize) and [disintegration/imaging](https://github.com/disintegration/imaging).

## Install

```
go get github.com/agentine/pixo
```

## Features

- **Resize** with 7 interpolation filters (NearestNeighbor, Box, Linear, CatmullRom, MitchellNetravali, Lanczos2, Lanczos3)
- **Fit** and **Fill** with aspect ratio preservation
- **Rotate** (arbitrary angle + 90/180/270), **Flip**, **Transpose**, **Transverse**
- **Crop** with anchor-based positioning (9 anchor points)
- **Adjust** brightness, contrast, gamma, saturation, sigmoid, hue
- **Effects**: Gaussian blur, sharpen, grayscale, invert, 3x3/5x5 convolution
- **Compose**: paste, overlay with alpha blending
- **I/O**: Open/Save/Encode/Decode with format detection, EXIF auto-orientation
- **Formats**: JPEG, PNG, GIF, BMP, TIFF (read/write), WebP (read)
- **Compatibility packages** for painless migration from nfnt/resize and disintegration/imaging

## Usage

### Resize an image

```go
package main

import (
    "log"
    "github.com/agentine/pixo"
)

func main() {
    img, err := pixo.Open("input.jpg")
    if err != nil {
        log.Fatal(err)
    }

    // Resize to 800x600
    resized := pixo.Resize(img, 800, 600, pixo.Lanczos3)
    pixo.Save(resized, "output.jpg", pixo.JPEGQuality(90))

    // Fit within 800x800, preserving aspect ratio
    fitted := pixo.Fit(img, 800, 800, pixo.Lanczos3)
    pixo.Save(fitted, "fitted.jpg")

    // Create a 200x200 thumbnail
    thumb := pixo.Thumbnail(img, 200, 200, pixo.Lanczos3)
    pixo.Save(thumb, "thumb.jpg")
}
```

### Transform

```go
rotated := pixo.Rotate90(img)
flipped := pixo.FlipH(img)
cropped := pixo.CropCenter(img, 500, 500)
```

### Adjust colors

```go
bright := pixo.AdjustBrightness(img, 20)   // +20%
contrast := pixo.AdjustContrast(img, 30)   // +30%
gray := pixo.Grayscale(img)
sharp := pixo.Sharpen(img, 1.0)
blurred := pixo.Blur(img, 2.0)
```

### Compose images

```go
result := pixo.Overlay(background, foreground, image.Pt(100, 100), 0.7)
```

## Migration from nfnt/resize

Change your import path:

```go
// Before
import "github.com/nfnt/resize"

// After
import "github.com/agentine/pixo/compat/resize"
```

The API is identical:

```go
// Before
m := resize.Resize(800, 0, img, resize.Lanczos3)

// After — same code, just different import
m := resize.Resize(800, 0, img, resize.Lanczos3)
```

Available interpolation functions: `NearestNeighbor`, `Bilinear`, `Bicubic`, `MitchellNetravali`, `Lanczos2`, `Lanczos3`.

## Migration from disintegration/imaging

Change your import path:

```go
// Before
import "github.com/disintegration/imaging"

// After
import "github.com/agentine/pixo/compat/imaging"
```

All functions have the same signatures:

```go
// Before
dst := imaging.Resize(src, 800, 600, imaging.Lanczos)
dst = imaging.Blur(src, 3.5)
dst = imaging.AdjustBrightness(src, 20)
err := imaging.Save(dst, "output.jpg")

// After — same code, just different import
dst := imaging.Resize(src, 800, 600, imaging.Lanczos)
dst = imaging.Blur(src, 3.5)
dst = imaging.AdjustBrightness(src, 20)
err := imaging.Save(dst, "output.jpg")
```

## Filters

| Filter | Description |
|--------|-------------|
| `NearestNeighbor` | Fastest, pixelated |
| `Box` | Fast, averaging |
| `Linear` | Bilinear, good balance |
| `CatmullRom` | Bicubic, sharp |
| `MitchellNetravali` | Bicubic, balanced |
| `Lanczos2` | High quality, a=2 |
| `Lanczos3` | Highest quality, a=3 |

## Benchmarks

Run benchmarks:

```
go test -bench=. -benchmem
```

## License

MIT
