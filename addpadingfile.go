package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var (
	regex = regexp.MustCompile("[[:xdigit:]]+")
)

// Input represents the data structure received by user.
type Input struct {
	Filename string
	Offset   string
	Data     string
}

//Main function
func main() {
	var input Input
	obtainInput(&input)
	addData(input)
}

//obtainInput takes user input data.
func obtainInput(input *Input) {
	filename := flag.String("file", "path/filename", "a string")
	offset := flag.String("offset", "0", "an integer")
	data := flag.String("data", "\x00\x00", "a string in hex format")

	flag.Parse()

	fmt.Println("- Parametro file: ", *filename)
	fmt.Println("- Parametro offset: ", *offset)
	fmt.Println("- Parametro data: ", *data)

	input.Filename = *filename
	input.Offset = *offset
	input.Data = *data
}

//addData read and create new file with the data that user input.
func addData(input Input) {
	reader, err := ioutil.ReadFile(input.Filename)
	check(err)
	file, err := os.Create(input.Filename)
	check(err)

	defer file.Close()

	//TODO: add input data from a offset.

	regMatch := regex.FindAllString(input.Data, -1)

	string2byte, err := hex.DecodeString(strings.Join(regMatch, ""))

	cleanData := [][]byte{[]byte(string2byte), reader}
	dataFinal := bytes.Join(cleanData, []byte(""))
	writer, err := file.Write(dataFinal)
	check(err)

	file.Sync()

	fmt.Printf("- Bytes writes: %d\n", writer)
	//fmt.Println(hex.Dump(file))
}

//Check error.
func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
