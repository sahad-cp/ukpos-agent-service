package escpos

import (
	"bytes"
	"encoding/base64"
	"image"
	// "image/color"
	_ "image/jpeg" // Register JPG decoder
	_ "image/png"  // Register PNG decoder
	"strings"
)

// ProcessBase64Logo converts a Base64 string to ESC/POS Raster Bytes
func ProcessBase64Logo(b64Str string) []byte {
	if b64Str == "" {
		return nil
	}

	// 1. Clean up string (Remove data URI prefix if present)
	if idx := strings.Index(b64Str, ","); idx != -1 {
		b64Str = b64Str[idx+1:]
	}

	// 2. Decode Base64 to Bytes
	rawDecoded, err := base64.StdEncoding.DecodeString(b64Str)
	if err != nil {
		return nil // Invalid Base64
	}

	// 3. Decode Image (PNG/JPG)
	img, _, err := image.Decode(bytes.NewReader(rawDecoded))
	if err != nil {
		return nil // Not a valid image
	}

	// 4. Resize & Rasterize (Width 300 is safe for 80mm)
	return imageToRaster(img, 576)
}

func imageToRaster(img image.Image, paperWidth int) []byte {
	bounds := img.Bounds()
	imgW := bounds.Dx()
	imgH := bounds.Dy()

	// Target logo width (safe for 80mm printers)
	targetW := 300
	if targetW > paperWidth {
		targetW = paperWidth
	}
	targetH := (imgH * targetW) / imgW

	// Center horizontally
	leftPad := (paperWidth - targetW) / 2
	if leftPad < 0 {
		leftPad = 0
	}

	byteWidth := (paperWidth + 7) / 8

	var buf bytes.Buffer

	// GS v 0  (Raster bit image)
	buf.Write([]byte{0x1D, 0x76, 0x30, 0x00})
	buf.WriteByte(byte(byteWidth % 256))
	buf.WriteByte(byte(byteWidth / 256))
	buf.WriteByte(byte(targetH % 256))
	buf.WriteByte(byte(targetH / 256))

	for y := 0; y < targetH; y++ {
		srcY := y * imgH / targetH

		for xByte := 0; xByte < byteWidth; xByte++ {
			var outByte byte

			for bit := 0; bit < 8; bit++ {
				x := xByte*8 + bit

				// Outside printable area → white
				if x < leftPad || x >= leftPad+targetW {
					continue
				}

				srcX := (x - leftPad) * imgW / targetW

				r, g, b, a := img.At(srcX, srcY).RGBA()

				// Transparent pixels → white
				if a < 0x8000 {
					continue
				}

				// Convert to grayscale (ITU-R BT.601)
				gray := (299*r + 587*g + 114*b) / 1000
				gray8 := gray >> 8

				// Thermal-friendly threshold
				if gray8 < 180 {
					outByte |= 1 << (7 - bit)
				}
			}

			buf.WriteByte(outByte)
		}
	}

	return buf.Bytes()
}


// func imageToRaster(img image.Image, paperWidth int) []byte {
// 	bounds := img.Bounds()
// 	imgW := bounds.Dx()
// 	imgH := bounds.Dy()

// 	// Resize image to max width (300)
// 	targetW := 300
// 	targetH := (imgH * targetW) / imgW

// 	// Calculate padding
// 	leftPad := (paperWidth - targetW) / 2


// 	byteWidth := (paperWidth + 7) / 8

// 	var buf bytes.Buffer
// 	buf.WriteString("\x1dv0\x00")
// 	buf.WriteByte(byte(byteWidth % 256))
// 	buf.WriteByte(byte(byteWidth / 256))
// 	buf.WriteByte(byte(targetH % 256))
// 	buf.WriteByte(byte(targetH / 256))

// 	for y := 0; y < targetH; y++ {
// 		for xByte := 0; xByte < byteWidth; xByte++ {
// 			var b byte
// 			for bit := 0; bit < 8; bit++ {
// 				x := xByte*8 + bit

// 				// Inside image region?
// 				if x >= leftPad && x < leftPad+targetW {
// 					srcX := (x - leftPad) * imgW / targetW
// 					srcY := y * imgH / targetH
// 					c := color.GrayModel.Convert(img.At(srcX, srcY)).(color.Gray)
// 					if c.Y < 128 {
// 						b |= 1 << (7 - bit)
// 					}
// 				}
// 			}
// 			buf.WriteByte(b)
// 		}
// 	}
// 	return buf.Bytes()
// }