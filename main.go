package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"

	"github.com/disintegration/imaging"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Action = func(c *cli.Context) error {
		imagePath := "fall.jpg"

		if c.NArg() > 0 {
			imagePath = c.Args().Get(0)
		}

		src, err := imaging.Open(imagePath)
		if err != nil {
			log.Fatalf("failed to open image: %v", err)
		}

		blurred := imaging.Blur(src, 20)

		// Open a test image.
		beer, err := imaging.Open("./beer_hand.png")
		if err != nil {
			log.Fatalf("failed to open image: %v", err)
		}

		blrect := blurred.Bounds()
		brect := beer.Bounds()

		dstRectY := float64(blrect.Dy()) * 0.7
		scaleY := (dstRectY / float64(brect.Dy()))

		scaledBeer := imaging.Resize(beer, 0, int(float64(brect.Dy())*scaleY), imaging.Lanczos)
		err = imaging.Save(scaledBeer, "scaled_beer.jpg")
		if err != nil {
			log.Fatalf("failed to save image: %v", err)
		}

		offset := image.Pt(0, blrect.Dy()-scaledBeer.Bounds().Dy())
		rgba := image.NewRGBA(blrect)
		draw.Draw(rgba, blrect, blurred, image.Point{0, 0}, draw.Src)
		draw.Draw(rgba, scaledBeer.Bounds().Add(offset), scaledBeer, image.Point{0, 0}, draw.Over)

		out, err := os.Create("out.jpg")
		if err != nil {
			fmt.Println(err)
		}

		var opt jpeg.Options
		opt.Quality = 100

		jpeg.Encode(out, rgba, &opt)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
