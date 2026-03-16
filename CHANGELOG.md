# Changelog

## 0.1.0 — 2026-03-16

Initial release.

- Pure Go image processing library (single dependency: `golang.org/x/image`)
- **Resize**: Nearest-neighbor, bilinear, bicubic, Lanczos with configurable output size
- **Transform**: Rotate (90/180/270), flip horizontal/vertical, transpose, transverse, arbitrary rotation
- **Crop**: Rectangle crop, center crop, auto-crop by color tolerance
- **Adjust**: Brightness, contrast, gamma, saturation, hue, color inversion
- **Effects**: Gaussian blur, sharpen, edge detection, emboss, grayscale, sepia
- **Compose**: Alpha compositing with over, multiply, screen, overlay blend modes
- **Histogram**: Per-channel histogram computation
- **I/O**: JPEG, PNG, GIF, BMP, TIFF encode/decode with WebP decode; format auto-detection; EXIF auto-orientation
- **Compatibility layers**:
  - `compat/resize`: drop-in replacement for `nfnt/resize`
  - `compat/imaging`: drop-in replacement for `disintegration/imaging`
- Comprehensive test suite (87.6% coverage) with golden image tests
- 15 benchmarks covering resize, effects, transforms, and I/O
