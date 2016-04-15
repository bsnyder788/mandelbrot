//Generate a color version of Mandelbrot set using super sampling.
package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"math/rand"
	"os"
)

const (
	xmin, ymin, xmax, ymax = -2, -2, +2, +2
	width, height          = 8194, 8194
)

var palette = make(map[int]color.Color)

// number of sub samples to randomly take per output pixel
var super int

func init() {
	//default to a sub sampling 4 pixels
	flag.IntVar(&super, "s", 4, "number of sub samples per pixel")
}

// Generate n color samples for the given output pixel.
func sample(x, y float64, n int) []color.Color {
	var samples []color.Color
	for i := 0; i < n; i++ {
		samples = append(samples, mandelbrot(complex(x+(rand.Float64()/width*4), y+(rand.Float64()/height*4))))
	}
	return samples
}

// Get the average RGB color from the slice of colors. (super sampling)
func getAvgSample(samples []color.Color) color.Color {
	rTot := uint32(0)
	bTot := uint32(0)
	gTot := uint32(0)
	aTot := uint32(0)
	for _, c := range samples {
		r, b, g, a := c.RGBA()
		rTot += r
		bTot += b
		gTot += g
		aTot += a
	}
	length := uint32(len(samples))
	rAvg := uint8(rTot / length)
	gAvg := uint8(gTot / length)
	bAvg := uint8(bTot / length)
	aAvg := uint8(aTot / length)
	return color.RGBA{rAvg, gAvg, bAvg, aAvg}
}

func main() {
	flag.Parse()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	//map each output pixel to a color for the final PNG render
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < height; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			//super sample using random sub pixels
			samples := sample(x, y, super)
			img.Set(px, py, getAvgSample(samples))
		}
	}
	png.Encode(os.Stdout, img)
}

// Get the color for a given point in complex plan of Mandelbrot set
func mandelbrot(z complex128) color.Color {
	const iterations = 500
	var v complex128
	for n := 0; n < iterations; n++ {
		v = v*v + z
		value := cmplx.Abs(v)
		if value > 2 {
			// calculate re-normalized escape count
			norm := n + 1 - int(math.Log(math.Log(value))/math.Log(2))
			if palette[norm] == nil {
				r, g, b := getAlteredBernsteinColors(n, iterations)
				palette[norm] = color.RGBA{r, g, b, 255}
			}
			return palette[norm]
		}
	}
	return color.Black
}

// Get continous band RGB color using altered bernstein polynomials
func getAlteredBernsteinColors(n, maxIter int) (uint8, uint8, uint8) {
	t := float64(n) / float64(maxIter)
	r := 9 * (1 - t) * t * t * t * 255
	g := 15 * (1 - t) * (1 - t) * t * t * 255
	b := 8.5 * (1 - t) * (1 - t) * (1 - t) * t * 255
	return uint8(r), uint8(g), uint8(b)
}
