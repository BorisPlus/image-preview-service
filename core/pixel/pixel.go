package pixel

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"strconv"
)

func monotonous(color color.RGBA, height, width int) ([]byte, error) {
	m := image.NewRGBA(image.Rect(0, 0, height, width))
	draw.Draw(m, m.Bounds(), &image.Uniform{color}, image.Point{}, draw.Src)
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, m, nil); err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

func pixel(color color.RGBA) []byte {
	pixel, _ := monotonous(color, 1, 1)
	return pixel
}

var (
	WhitePixel     = pixel(color.RGBA{255, 255, 255, 255})
	WhitePixelLen  = strconv.Itoa(len(WhitePixel))
	GrayPixel      = pixel(color.RGBA{220, 220, 220, 255})
	BlackPixel     = pixel(color.RGBA{0, 0, 0, 255})
	OrangePixel    = pixel(color.RGBA{204, 85, 0, 255})
	RedPixel       = pixel(color.RGBA{255, 0, 0, 255})
	GreenPixel     = pixel(color.RGBA{0, 255, 0, 255})
)
