package main

import (
	"flag"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"math/rand"
	"mime"
	"os"
	"path/filepath"
	"time"
)

func main() {
	in := flag.String("i", "", "image to fracture")
	out := flag.String("o", "", "output file")

	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	imgData := load(*in)

	mimeType := mime.TypeByExtension(filepath.Ext(*in))

	save(*out, imgData, mimeType)
}

func load(filePath string) (colorGrid [][]color.Color) {
	file, err := os.Open(filePath)

	if err != nil {
		log.Println("Error reading file:", err)
	}

	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Println("Error decoding file:", err)
	}

	size := img.Bounds().Size()

	next := size.Y

	minAngle := 10
	maxAngle := 360

	angle := rand.Intn(maxAngle-minAngle) + minAngle

	for x := 0; x < size.X; x++ {
		var colors []color.Color

		xPos := x

		// Don't make too random. It looks bad
		if x%int((size.X+x)/2) == 0 {
			angle = rand.Intn(maxAngle-minAngle) + minAngle
		}

		for y := 0; y < size.Y; y++ {
			yPos := y

			xCenter := float64(size.X / 2)
			yCenter := float64(size.Y / 2)

			dx := float64(x) - xCenter
			dy := float64(y) - yCenter

			c := math.Cos(float64(angle))
			s := math.Sin(float64(angle))

			rotX := dx*c - dy*s
			rotY := dx*s + dy*c

			xPos = int(rotX + xCenter)

			// Produces a more interesting result
			xPos += int(math.Abs(float64(x - y)))

			// Handle position out of range
			if xPos < 0 || xPos > size.X {
				xPos = x
			}

			if y <= next {
				// Append pixels that haven't been rotated
				colors = append(colors, img.At(x, y))

				next -= yPos
			} else {
				yPos = int(rotY + yCenter)

				// Handle position out of range
				if xPos < 0 || xPos > size.X {
					xPos = x
				}

				if yPos < 0 || yPos > size.Y {
					yPos = y
				}

				// Append rotated pixels
				colors = append(colors, img.At(xPos, yPos))

				next = size.Y - yPos - y
			}
		}

		colorGrid = append(colorGrid, colors)
	}

	return
}

func save(filePath string, colorGrid [][]color.Color, mimeType string) {
	xLen, yLen := len(colorGrid), len(colorGrid[0])

	rect := image.Rect(0, 0, xLen, yLen)

	img := image.NewNRGBA(rect)

	for x := 0; x < xLen; x++ {
		for y := 0; y < yLen; y++ {
			img.Set(x, y, colorGrid[x][y])
		}
	}

	file, err := os.Create(filePath)

	if err != nil {
		log.Println("Error creating image:", err)
	}

	defer file.Close()

	switch mimeType {
	case "image/jpeg":
		err = jpeg.Encode(file, img.SubImage(img.Rect), nil)
	case "image/png":
		err = png.Encode(file, img.SubImage(img.Rect))
	case "image/gif":
		err = gif.Encode(file, img.SubImage(img.Rect), nil)
	}

	if err != nil {
		log.Println("Error encoding image:", err)
	}
}
