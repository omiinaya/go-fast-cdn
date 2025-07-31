package testutils

import (
	"bytes"
	"crypto/md5"
	"image"
	"image/color"
	"image/png"
	"io"
)

// CreateDummyImage creates a dummy image for testing
func CreateDummyImage(width, height int) (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Fill with a simple pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{uint8(x % 256), uint8(y % 256), uint8((x + y) % 256), 255})
		}
	}
	return img, nil
}

// EncodeImage encodes an image to a writer
func EncodeImage(w io.Writer, img image.Image) error {
	return png.Encode(w, img)
}

// CreateDummyDocument creates a dummy document content for testing
func CreateDummyDocument() []byte {
	return []byte(`Lorem Ipsum is simply dummy text of the printing and typesetting industry. 
Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown 
printer took a galley of type and scrambled it to make a type specimen book. It has survived not 
only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged.`)
}

// CalculateImageChecksum calculates the MD5 checksum of an image
func CalculateImageChecksum(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := EncodeImage(&buf, img)
	if err != nil {
		return nil, err
	}
	checksum := md5.Sum(buf.Bytes())
	return checksum[:], nil
}

// CalculateDocumentChecksum calculates the MD5 checksum of a document
func CalculateDocumentChecksum(docContent []byte) []byte {
	checksum := md5.Sum(docContent)
	return checksum[:]
}
