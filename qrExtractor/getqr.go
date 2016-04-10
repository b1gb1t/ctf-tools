package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
)

var (
	background color.Color = color.RGBA{200, 200, 200, 255}
	green color.Color = color.RGBA{100, 150, 90, 255}
	lightgreen color.Color = color.RGBA{100, 200, 90, 255}
	yellow color.Color = color.RGBA{200, 200, 90, 255}
	blue color.Color = color.RGBA{100, 150, 200, 255}
	lightblue color.Color = color.RGBA{100, 200, 200, 255}
	pink color.Color = color.RGBA{200, 150, 200, 255}
	orange color.Color = color.RGBA{200, 150, 90, 255}
	white color.Color = color.RGBA{255, 255, 255, 255}
	//black := color.RGBA{0, 0, 0, 255}
	color_selected color.Color = lightgreen
)
func init() {
	var err string
	fmt.Sprintf(err, "Correct usage:\n%s -in <image> -e | -in <image> -o <output_file>\n", os.Args[0])
	flag.ErrHelp = errors.New(err)

	flag.Usage = func() {
		fmt.Fprintf(os.Stdin, "Correct usage:\n%s -in <image> -e | -in <image> -o <output_file>\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	var imagen string
	flag.StringVar(&imagen, "in", "", "a string value")
	var extract bool
	flag.BoolVar(&extract, "e", false, "true or false")
	var fout string
	flag.StringVar(&fout, "o", "", "output file")
	flag.Parse()

	if flag.Parsed() && flag.NFlag() == 2 && imagen != "" {
		decodeImage(imagen, extract, fout)
	} else {
		fmt.Printf("More information: -h\n")
	}
}

func decodeImage(file string, extract bool, fout string) {
	reader, err := os.Open(file)
	if err != nil {
		fmt.Printf("--- Open error ---\n")
	}
	defer reader.Close()

	//image contain RGBA colors
	img, _, err := image.Decode(reader)
	if err != nil {
		fmt.Printf("--- Decode error ---\n")
	}

	bound := img.Bounds()

	if extract {
		getColors(img, bound)
	} else if fout != "" {
		getQRs(img, bound, fout)
	} else {
		fmt.Printf("More information: -h\n")
	}
}

// Get unique image colors
func getColors(img image.Image, bound image.Rectangle) {
	var colores []color.Color
	aux_colores := map[color.Color]bool{}
	for i := 0; i < bound.Max.Y; i++ {
		for j := 0; j < bound.Max.X; j++ {
			pixel := img.At(j, i)
			if !aux_colores[pixel] {
				aux_colores[pixel] = true
				colores = append(colores, pixel)
			}
		}
	}
	fmt.Printf("Colores:\n%v\n", colores)
	fmt.Printf("Bound: %v", bound) //Medidas en cuadricula (xi,yi)(xf,yf)
}

func getQRs(img image.Image, bound image.Rectangle, fout string) {
	file, err := os.Create(fout)
	if err != nil {
		fmt.Println("Error - Not create new image")
	}
	defer file.Close()

	newImg := image.NewRGBA(image.Rect(bound.Min.X, bound.Min.Y, bound.Max.X, bound.Max.Y))
	// 1 module = 17.08 px
	md2px := 16
	corner1 := image.Rect(14*md2px, 0, 23*md2px, 8*md2px)
	corner2 := image.Rect(0, 14*md2px, 8*md2px, 23*md2px)
	corner3 := image.Rect(14*md2px, 14*md2px, 23*md2px, 23*md2px)

	config := map[color.Color]bool{}
	config[background] = false // Always false
	config[green] = false
	config[lightgreen] = true
	config[yellow] = true
	config[blue] = true
	config[lightblue] = true
	config[pink] = true
	config[orange] = true
	
	if config[green] == true {
		for i := bound.Min.Y; i < bound.Max.Y; i++ {
			for j := bound.Min.X; j < bound.Max.X; j++ {
				pixel := img.At(j, i)
				if config[pixel] {
					newImg.Set(j, i, pixel)
				} else {
					newImg.Set(j, i, white)
				}
			}
		}
	} else {
		for i := bound.Min.Y; i < bound.Max.Y; i++ {
			for j := bound.Min.X; j < bound.Max.X; j++ {
				pixel := img.At(j, i)

				//Print the 3 corners
				if i >= corner1.Min.Y && i <= corner1.Max.Y && j >= corner1.Min.X && j <= corner1.Max.X {
					if pixel == green {

						newImg.Set(j, i, color_selected)
					} else {
						newImg.Set(j, i, white)
					}
				} else if i >= corner2.Min.Y && i <= corner2.Max.Y && j >= corner2.Min.X && j <= corner2.Max.X {
					if pixel == green {
						newImg.Set(j, i, color_selected)
					} else {
						newImg.Set(j, i, white)
					}
				} else if i >= corner3.Min.Y && i <= corner3.Max.Y && j >= corner3.Min.X && j <= corner3.Max.X {
					if pixel == green {
						newImg.Set(j, i, color_selected)
					} else {
						newImg.Set(j, i, white)
					}
				} else if config[pixel] { 		// Print the rest of image
					newImg.Set(j, i, pixel)
				} else {
					newImg.Set(j, i, white)
				}
			}
		}
	}

	png.Encode(file, newImg)

}
