//Generate a colored version of Mandelbrot set
package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"math/rand"
	"os"
)

var palette = make(map[int]color.Color)

func main() {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 1024, 1024
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < height; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			img.Set(px, py, mandelbrot(z))
		}
	}
	png.Encode(os.Stdout, img)
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		value := int(cmplx.Abs(v))
		if value > 2 {
			if palette[value] == nil {
				r := uint8(rand.Intn(255))
				g := uint8(rand.Intn(255))
				b := uint8(rand.Intn(255))
				palette[value] = color.RGBA{r, g, b, 255 - contrast*n}
			}
			return palette[value]
		}
	}
	return color.Black
}
