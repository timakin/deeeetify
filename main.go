package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/urfave/cli"
)

const (
	defaultBlurStrength = 20
	defaultBackGround   = "https://github.com/timakin/deeeetify/blob/master/fall.jpg"
	beerHandFile        = "https://github.com/timakin/deeeetify/blob/master/beer_hand.png"
)

func main() {
	app := cli.NewApp()

	app.Action = func(c *cli.Context) error {
		imagePath := defaultBackGround
		blurStrength := defaultBlurStrength

		if c.NArg() > 0 {
			imagePath = c.Args().Get(0)
			if c.Args().Get(1) != "" {
				bs, err := strconv.Atoi(c.Args().Get(1))
				if err != nil {
					log.Fatalf("blur strength must be specified with number")
				}
				blurStrength = bs
			}
		}

		src, err := imaging.Open(imagePath)
		if err != nil {
			log.Fatalf("failed to open image: %v", err)
		}

		blurred := imaging.Blur(src, float64(blurStrength))

		beer, err := imaging.Open(beerHandFile)
		if err != nil {
			log.Fatalf("failed to open image: %v", err)
		}

		blrect := blurred.Bounds()
		brect := beer.Bounds()

		dstRectY := float64(blrect.Dy()) * 0.7
		scaleY := (dstRectY / float64(brect.Dy()))

		scaledBeer := imaging.Resize(beer, 0, int(float64(brect.Dy())*scaleY), imaging.Lanczos)

		offset := image.Pt(0, blrect.Dy()-scaledBeer.Bounds().Dy())
		rgba := image.NewRGBA(blrect)
		draw.Draw(rgba, blrect, blurred, image.Point{0, 0}, draw.Src)
		draw.Draw(rgba, scaledBeer.Bounds().Add(offset), scaledBeer, image.Point{0, 0}, draw.Over)

		out, err := os.Create("deeeeted.jpg")
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
