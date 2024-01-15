package netpbm

import (
	"bufio"
	"os"
	"strconv"
)

// PPM represents a Portable PixMap image.
type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           int
}

// Pixel represents a single pixel with red, green, and blue values.
type Pixel struct {
	R, G, B uint8
}

// ReadPPM reads a PPM image from a file and returns a struct that represents the image.
func ReadPPM(filename string) (*PPM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	ppm := &PPM{}

	// Read magic number
	scanner.Scan()
	ppm.magicNumber = scanner.Text()

	// Read width, height, and max value
	scanner.Scan()
	ppm.width, _ = strconv.Atoi(scanner.Text())
	scanner.Scan()
	ppm.height, _ = strconv.Atoi(scanner.Text())
	scanner.Scan()
	ppm.max, _ = strconv.Atoi(scanner.Text())

	// Read pixel data
	ppm.data = make([][]Pixel, ppm.height)
	for i := 0; i < ppm.height; i++ {
		ppm.data[i] = make([]Pixel, ppm.width)
		for j := 0; j < ppm.width; j++ {
			scanner.Scan()
			r, _ := strconv.Atoi(scanner.Text())
			scanner.Scan()
			g, _ := strconv.Atoi(scanner.Text())
			scanner.Scan()
			b, _ := strconv.Atoi(scanner.Text())
			ppm.data[i][j] = Pixel{uint8(r), uint8(g), uint8(b)}
		}
	}

	return ppm, nil
}
