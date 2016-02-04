package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"strings"
)

const (
	black = "0"
	red   = "1"
)

func init() {
	var err string
	fmt.Sprintf(err, "Correct usage:\n%s <image>\n", os.Args[0])
	flag.ErrHelp = errors.New(err)

	flag.Usage = func() {
		fmt.Fprintf(os.Stdin, "Correct usage:\n%s <image>\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if len(os.Args) == 2 {
		decodeImage(os.Args[1])
	} else {
		fmt.Printf("More information: -h\n")
	}
}

func decodeImage(file string) {
	info := ""

	reader, err := os.Open(file)
	if err != nil {
		fmt.Printf("--- Open error ---\n")
	}
	defer reader.Close()

	//image contain RGBA colors
	image, _, err := image.Decode(reader)
	if err != nil {
		fmt.Printf("--- Decode error ---\n")
	}

	bound := image.Bounds()

	var msg []string
	for i := 0; i < bound.Max.Y; i++ {
		for j := 0; j < bound.Max.X; j++ {
			pixel := image.At(j, i)
			r, _, _, _ := pixel.RGBA()

			if r == 0 {
				info = info + black
			} else {
				info = info + red
			}
		}

		num, err := strconv.ParseInt(info, 2, 0)
		if err != nil {
			panic(err)
		}

		msg = append(msg, string(int(num)))
		info = ""
	}

	fmt.Printf("\nMessage:\n%s\n\n", strings.Join(msg, ""))
}
