package main

import (
	"flag"
	"image"
	"log"
	"os"
	"time"

	"vallon.me/place"

	"image/color"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (
	username, password, imagePath string
	startX, startY                int
)

func init() {
	flag.StringVar(&username, "u", "", "reddit username")
	flag.StringVar(&password, "p", "", "reddit password")
	flag.StringVar(&imagePath, "i", "", "input image to plot")

	flag.IntVar(&startX, "x", 0, "starting x coordinate")
	flag.IntVar(&startY, "y", 0, "starting y coordinate")

	flag.Parse()

	switch "" {
	case username, password, imagePath:
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func Same(c1, c2 color.Color) bool {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2
}

func main() {
	img, err := OpenImage(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	client, err := place.Login(username, password)
	if err != nil {
		log.Fatal(err)
	}

	snap, err := place.GetBitmap()
	if err != nil {
		log.Fatal(err)
	}

	dem := img.Bounds().Max
	for y := 0; y < dem.Y; y++ {
		for x := 0; x < dem.X; x++ {
			color := img.At(x, y)

			if Same(color, snap.At(startX+x, startY+y)) {
				continue
			}

			log.Printf("placing tile (%d, %d) at (%d, %d)\n", x, y, startX+x, startY+y)

			d, err := client.Draw(startX+x, startY+y, color)
			if err != nil {
				log.Fatal(err)
			}

			log.Println("sleeping for:", d)
			time.Sleep(d)
		}
	}
}

func OpenImage(name string) (image.Image, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}
