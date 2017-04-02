package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"time"

	"vallon.me/place"
)

const timeout = 1 * time.Minute

func main() {
	for {
		name := fmt.Sprintf("place-%d.png", time.Now().UnixNano())

		log.Printf("[%s] taking snapshot", name)
		if err := takeSnapshot(name); err != nil {
			log.Printf("[%s] %s", name, err)
		} else {
			log.Printf("[%s] saved", name)
		}

		time.Sleep(timeout)
	}
}

func takeSnapshot(name string) error {
	img, err := place.GetBitmap()
	if err != nil {
		return err
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}

	return f.Close()
}
