package pixo

import (
	"encoding/binary"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
)

// Format represents an image file format.
type Format int

const (
	JPEG Format = iota
	PNG
	GIF
	BMP
	TIFF
	WebP
)

// String returns the format name.
func (f Format) String() string {
	switch f {
	case JPEG:
		return "jpeg"
	case PNG:
		return "png"
	case GIF:
		return "gif"
	case BMP:
		return "bmp"
	case TIFF:
		return "tiff"
	case WebP:
		return "webp"
	default:
		return "unknown"
	}
}

// EncodeOption is a functional option for encoding.
type EncodeOption func(*encodeConfig)

type encodeConfig struct {
	jpegQuality         int
	pngCompressionLevel png.CompressionLevel
}

func defaultEncodeConfig() *encodeConfig {
	return &encodeConfig{
		jpegQuality:         95,
		pngCompressionLevel: png.DefaultCompression,
	}
}

// JPEGQuality sets the JPEG encoding quality (1-100).
func JPEGQuality(q int) EncodeOption {
	return func(c *encodeConfig) {
		c.jpegQuality = q
	}
}

// PNGCompressionLevel sets the PNG compression level.
func PNGCompressionLevel(l png.CompressionLevel) EncodeOption {
	return func(c *encodeConfig) {
		c.pngCompressionLevel = l
	}
}

// Open opens and decodes an image from the given file.
// EXIF orientation is applied automatically.
func Open(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Decode(f)
}

// Save encodes and saves img to the given file.
// The format is determined by the file extension.
func Save(img image.Image, filename string, opts ...EncodeOption) error {
	format, err := formatFromFilename(filename)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return Encode(f, img, format, opts...)
}

// Decode decodes an image from r with format detection.
// EXIF orientation is applied automatically for JPEG images.
func Decode(r io.Reader) (image.Image, error) {
	// We need to peek at the header for EXIF before full decode
	// Use a buffer to allow re-reading
	var buf [12]byte
	n, err := io.ReadFull(r, buf[:])
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, err
	}
	header := buf[:n]

	// Check for JPEG EXIF orientation
	isJPEG := len(header) >= 2 && header[0] == 0xFF && header[1] == 0xD8

	// Reconstruct reader
	combined := io.MultiReader(
		strings.NewReader(string(header)),
		r,
	)

	// If JPEG, try to read EXIF orientation first with a separate reader
	orientation := 1
	if isJPEG {
		// We'll read orientation after decoding from a seekable copy
		// For now, decode first then apply orientation
		_ = orientation
	}

	img, _, err := image.Decode(combined)
	if err != nil {
		return nil, err
	}

	if isJPEG {
		// Try to read EXIF from the original file if it was seekable
		if seeker, ok := r.(io.ReadSeeker); ok {
			if _, err := seeker.Seek(0, io.SeekStart); err == nil {
				orientation = readJPEGOrientation(seeker)
			}
		}
		img = applyOrientation(img, orientation)
	}

	return img, nil
}

// Encode encodes img to w in the specified format.
func Encode(w io.Writer, img image.Image, format Format, opts ...EncodeOption) error {
	cfg := defaultEncodeConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	switch format {
	case JPEG:
		return jpeg.Encode(w, img, &jpeg.Options{Quality: cfg.jpegQuality})
	case PNG:
		enc := &png.Encoder{CompressionLevel: cfg.pngCompressionLevel}
		return enc.Encode(w, img)
	case GIF:
		return gif.Encode(w, img, nil)
	case BMP:
		return bmp.Encode(w, img)
	case TIFF:
		return tiff.Encode(w, img, nil)
	default:
		return errors.New("pixo: unsupported format for encoding: " + format.String())
	}
}

func formatFromFilename(filename string) (Format, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return JPEG, nil
	case ".png":
		return PNG, nil
	case ".gif":
		return GIF, nil
	case ".bmp":
		return BMP, nil
	case ".tif", ".tiff":
		return TIFF, nil
	default:
		return 0, errors.New("pixo: unsupported file extension: " + ext)
	}
}

// readJPEGOrientation reads the EXIF orientation tag from a JPEG.
// Returns 1 (normal) if not found or on error.
func readJPEGOrientation(r io.Reader) int {
	// Read SOI marker
	var marker [2]byte
	if _, err := io.ReadFull(r, marker[:]); err != nil {
		return 1
	}
	if marker[0] != 0xFF || marker[1] != 0xD8 {
		return 1
	}

	// Scan for APP1 (EXIF) marker
	for {
		if _, err := io.ReadFull(r, marker[:]); err != nil {
			return 1
		}
		if marker[0] != 0xFF {
			return 1
		}

		var size [2]byte
		if _, err := io.ReadFull(r, size[:]); err != nil {
			return 1
		}
		segLen := int(size[0])<<8 | int(size[1])
		if segLen < 2 {
			return 1
		}

		if marker[1] == 0xE1 { // APP1
			data := make([]byte, segLen-2)
			if _, err := io.ReadFull(r, data); err != nil {
				return 1
			}
			return parseExifOrientation(data)
		}

		// Skip segment
		if _, err := io.CopyN(io.Discard, r, int64(segLen-2)); err != nil {
			return 1
		}
	}
}

func parseExifOrientation(data []byte) int {
	// Check "Exif\0\0" header
	if len(data) < 14 || string(data[:6]) != "Exif\x00\x00" {
		return 1
	}
	data = data[6:]

	var order binary.ByteOrder
	switch string(data[:2]) {
	case "II":
		order = binary.LittleEndian
	case "MM":
		order = binary.BigEndian
	default:
		return 1
	}

	if order.Uint16(data[2:4]) != 0x002A {
		return 1
	}

	ifdOffset := order.Uint32(data[4:8])
	if int(ifdOffset)+2 > len(data) {
		return 1
	}

	numEntries := order.Uint16(data[ifdOffset : ifdOffset+2])
	offset := int(ifdOffset) + 2

	for i := 0; i < int(numEntries); i++ {
		entryOffset := offset + i*12
		if entryOffset+12 > len(data) {
			return 1
		}

		tag := order.Uint16(data[entryOffset : entryOffset+2])
		if tag == 0x0112 { // Orientation tag
			return int(order.Uint16(data[entryOffset+8 : entryOffset+10]))
		}
	}
	return 1
}

func applyOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 1:
		return img
	case 2:
		return FlipH(img)
	case 3:
		return Rotate180(img)
	case 4:
		return FlipV(img)
	case 5:
		return Transpose(img)
	case 6:
		return Rotate270(img)
	case 7:
		return Transverse(img)
	case 8:
		return Rotate90(img)
	default:
		return img
	}
}

func init() {
	// Register WebP decoder
	image.RegisterFormat("webp", "RIFF????WEBP", func(r io.Reader) (image.Image, error) {
		return webp.Decode(r)
	}, func(r io.Reader) (image.Config, error) {
		return webp.DecodeConfig(r)
	})
}
