# pixo — Pure Go Image Processing Library

## Overview

**Replaces:** nfnt/resize (3,319 importers, archived Sept 2022) and disintegration/imaging (2,899 importers, 5.7k stars, dormant since March 2023, 31 open issues)

**Package:** `github.com/agentine/pixo`

**Language:** Go

## Why

### nfnt/resize
- Archived September 2022; maintainer explicitly stated "I won't really look at PRs and issues anymore"
- No go.mod file, no tagged version, no stable v1 release
- Last meaningful commit: February 2018
- 3,319 importers with no maintained drop-in replacement
- Maintainer recommended golang.org/x/image/draw (too low-level) and imagick (requires CGo)
- Community suggested gofrs adoption, never happened

### disintegration/imaging
- Last commit: March 2023 (3+ years ago)
- 31 open issues, 3 unmerged PRs
- No EXIF support (requested since 2016, issue #15)
- No WebP support
- 5.7k stars, 2,899 importers
- Effectively unmaintained — no response to issues or PRs

### Market Gap
- golang.org/x/image/draw: low-level scaling only, no crop/rotate/blur/adjust operations
- bimg/govips: require libvips C dependency — not pure Go
- bild: semi-active (v0.14.0, July 2024) but different API design, not a drop-in for either target
- No maintained pure Go library provides the simple, batteries-included API that imaging users expect

## Architecture

### Core Package (`pixo`)
Single-package design matching imaging's simplicity. All functions accept `image.Image` and return `*image.NRGBA`.

### Modules

1. **Resize** — `Resize()`, `Fit()`, `Fill()`, `Thumbnail()`
   - Interpolation filters: NearestNeighbor, Box, Linear, CatmullRom, MitchellNetravali, Lanczos2, Lanczos3
   - Aspect ratio preservation
   - Concurrent processing for large images

2. **Transform** — `Rotate()`, `Rotate90/180/270()`, `FlipH()`, `FlipV()`, `Transpose()`, `Transverse()`
   - Arbitrary angle rotation with configurable background color

3. **Crop** — `Crop()`, `CropAnchor()`, `CropCenter()`
   - Anchor-based cropping (TopLeft, Top, TopRight, Left, Center, Right, BottomLeft, Bottom, BottomRight)

4. **Adjust** — `AdjustBrightness()`, `AdjustContrast()`, `AdjustGamma()`, `AdjustSaturation()`, `AdjustSigmoid()`, `AdjustHue()`
   - Color space-aware adjustments

5. **Effects** — `Blur()`, `Sharpen()`, `Grayscale()`, `Invert()`, `Convolve3x3()`, `Convolve5x5()`
   - Gaussian blur with configurable sigma

6. **Compose** — `Paste()`, `PasteCenter()`, `Overlay()`, `OverlayCenter()`
   - Alpha-aware compositing

7. **I/O** — `Open()`, `Save()`, `Decode()`, `Encode()`
   - Format detection from extension/content
   - JPEG, PNG, GIF, BMP, TIFF support
   - WebP decode support (via golang.org/x/image/webp)
   - EXIF auto-orientation on decode
   - Configurable encode options (quality, compression)

8. **Histogram** — `Histogram()`
   - Per-channel histogram computation

### Compatibility Packages

- `pixo/compat/resize` — Drop-in API compatibility with nfnt/resize
  - Same function signatures: `Resize()`, `Thumbnail()`
  - Same interpolation function constants
  - Users change import path only

- `pixo/compat/imaging` — Drop-in API compatibility with disintegration/imaging
  - Same function signatures for all operations
  - Same option types and constants
  - Users change import path only

## Deliverables

1. Core image processing library with all modules above
2. Compatibility packages for nfnt/resize and disintegration/imaging
3. Comprehensive test suite with golden image tests
4. Benchmarks comparing against nfnt/resize, imaging, and x/image/draw
5. Migration guide from both target libraries
6. Examples and documentation

## Technical Decisions

- **Go 1.21+** minimum (slices, maps packages)
- **Zero external dependencies** for core (only stdlib)
- **golang.org/x/image** as optional dependency for WebP decode
- **Parallel processing** using worker pools for large image operations
- **image.NRGBA** as canonical output format (matching imaging convention)
- **Functional options** for I/O operations (e.g., `Save(img, "out.jpg", pixo.JPEGQuality(90))`)
