package netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           int
}

// ReadPGM reads a PGM image from a file and returns a struct that represents the image.
func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read magic number
	scanner.Scan()
	magicNumber := scanner.Text()

	// Read comment lines
	for strings.HasPrefix(scanner.Text(), "#") {
		scanner.Scan()
	}

	// Read width, height, and max value
	scanner.Scan()
	dimensions := strings.Fields(scanner.Text())
	width, _ := strconv.Atoi(dimensions[0])
	height, _ := strconv.Atoi(dimensions[1])

	scanner.Scan()
	max, _ := strconv.Atoi(scanner.Text())

	// Read pixel data
	data := make([][]uint8, height)
	for i := range data {
		data[i] = make([]uint8, width)
		for j := 0; j < width; j++ {
			scanner.Scan()
			value, _ := strconv.Atoi(scanner.Text())
			data[i][j] = uint8(value)
		}
	}

	return &PGM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         max,
	}, nil
}

// Size returns the width and height of the image.
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At returns the value of the pixel at (x, y).
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Save saves the PGM image to a file and returns an error if there was a problem.
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write magic number, width, height, and max value
	file.WriteString(fmt.Sprintf("%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max))

	// Write pixel data
	for _, row := range pgm.data {
		for _, value := range row {
			file.WriteString(fmt.Sprintf("%d\n", value))
		}
	}

	return nil
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {
	for i := range pgm.data {
		for j := range pgm.data[i] {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {
	for i := range pgm.data {
		for j, k := 0, len(pgm.data[i])-1; j < k; j, k = j+1, k-1 {
			pgm.data[i][j], pgm.data[i][k] = pgm.data[i][k], pgm.data[i][j]
		}
	}
}

// Flop flops the PGM image vertically.
func (pgm *PGM) Flop() {
	for i, j := 0, len(pgm.data)-1; i < j; i, j = i+1, j-1 {
		pgm.data[i], pgm.data[j] = pgm.data[j], pgm.data[i]
	}
}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	pgm.max = int(maxValue)
}

// Rotate90CW rotates the PGM image 90° clockwise.
func (pgm *PGM) Rotate90CW() {
	newData := make([][]uint8, pgm.width)
	for i := 0; i < pgm.width; i++ {
		newData[i] = make([]uint8, pgm.height)
		for j := 0; j < pgm.height; j++ {
			newData[i][j] = pgm.data[pgm.height-j-1][i]
		}
	}
	pgm.data = newData
	pgm.width, pgm.height = pgm.height, pgm.width
}

// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM {

	return nil
}

/*
func main() {
	pgm, err := ReadPGM("example.pgm")
	if err != nil {
		fmt.Println("Error reading PGM:", err)
		return
	}

	width, height := pgm.Size()
	fmt.Printf("Image size: %d x %d\n", width, height)

	// Perform other operations as needed.
}*/
