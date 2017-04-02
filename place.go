package place

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"net/http"
)

const (
	infoURL = "https://www.reddit.com/api/place/pixel.json"
	bmpURL  = "https://www.reddit.com/api/place/board-bitmap"
)

func GetPixel(x, y int) (*Info, error) {
	query := fmt.Sprintf("?x=%d&y=%d", x, y)

	res, err := http.Get(infoURL + query)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	defer res.Body.Close()

	var inf Info
	if err := json.NewDecoder(res.Body).Decode(&inf); err != nil {
		return nil, err
	}

	return &inf, nil
}

type Info struct {
	Username  string  `json:"user_name"`
	Timestamp float64 `json:"timestamp"`

	Tile
}

type Tile struct {
	Y     int          `json:"y"`
	X     int          `json:"x"`
	Color paletteIndex `json:"color"`
}

type paletteIndex int

func (i paletteIndex) Color() color.Color {
	return DefaultColorPalette[i]
}

func GetBitmap() (image.Image, error) {
	res, err := http.Get(bmpURL)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	defer res.Body.Close()

	img := image.NewPaletted(image.Rect(0, 0, 1000, 1000), DefaultColorPalette)

	r := bufio.NewReader(res.Body)

	var whatIsThis [4]byte
	r.Read(whatIsThis[:])

	for i := 0; i < len(img.Pix); i += 2 {
		pix, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		img.Pix[i] = pix >> 4
		img.Pix[i+1] = pix & 15
	}

	return img, nil
}

var DefaultColorPalette = color.Palette{
	color.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
	color.RGBA{0xE4, 0xE4, 0xE4, 0xFF},
	color.RGBA{0x88, 0x88, 0x88, 0xFF},
	color.RGBA{0x22, 0x22, 0x22, 0xFF},
	color.RGBA{0xFF, 0xA7, 0xD1, 0xFF},
	color.RGBA{0xE5, 0x00, 0x00, 0xFF},
	color.RGBA{0xE5, 0x95, 0x00, 0xFF},
	color.RGBA{0xA0, 0x6A, 0x42, 0xFF},
	color.RGBA{0xE5, 0xD9, 0x00, 0xFF},
	color.RGBA{0x94, 0xE0, 0x44, 0xFF},
	color.RGBA{0x02, 0xBE, 0x01, 0xFF},
	color.RGBA{0x00, 0xD3, 0xDD, 0xFF},
	color.RGBA{0x00, 0x83, 0xC7, 0xFF},
	color.RGBA{0x00, 0x00, 0xEA, 0xFF},
	color.RGBA{0xCF, 0x6E, 0xE4, 0xFF},
	color.RGBA{0x82, 0x00, 0x80, 0xFF},
}
