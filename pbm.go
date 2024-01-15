package Netpbm

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read magic number
	magicNumber := readNonCommentLine(scanner)

	// Read width and height
	width, height := 0, 0
	fmt.Sscanf(readNonCommentLine(scanner), "%d %d", &width, &height)

	// Read image data
	data := make([][]bool, height)
	if magicNumber == "P1" {
		// P1 format (ASCII)
		for i := 0; i < height && scanner.Scan(); i++ {
			line := strings.Fields(scanner.Text()) // Split the line into fields
			data[i] = make([]bool, width)
			for j, field := range line {
				if j >= width {
					return nil, fmt.Errorf("Index out of range")
				}
				if field == "1" {
					data[i][j] = true
				} else if field == "0" {
					data[i][j] = false
				}
			}
		}
	} else if magicNumber == "P4" {
		// P4 format (binary)
		for i := 0; i < height && scanner.Scan(); i++ {
			line := scanner.Text()
			data[i] = make([]bool, width)

			for j := 0; j < width/8; j++ { // Modifier l'indice ici
				// Read rune by rune
				if j < len(line) {
					char := line[j]

					// Convert to ASCII
					asciiValue := int(char)

					// Convert ASCII to hex
					hexValue := fmt.Sprintf("%02X", asciiValue)

					// Convert hex to binary
					binaryValue, err := strconv.ParseUint(hexValue, 16, 8)
					fmt.Println(binaryValue)
					if err != nil {
						// Handle error
						fmt.Println("Error:", err)
						return nil, err
					}

					// Set the corresponding bits in data
					for k := 0; k < 8; k++ {
						data[i][j*8+k] = (binaryValue>>uint(7-k))&1 == 1
					}
				} else {
					// Padding case
					data[i][j*8] = false
				}
			}
		}

	} else {
		return nil, fmt.Errorf("Invalid magic number")
	}

	return &PBM{data, width, height, magicNumber}, nil
}

// readNonCommentLine reads the next non-comment line from the scanner.
func readNonCommentLine(scanner *bufio.Scanner) string {
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[0] != '#' {
			return line
		}
	}
	return ""
}

// Size returns the width and height of the image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// At returns the value of the pixel at (x, y).
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

// Save saves the PBM image to a file and returns an error if there was a problem.
func (pbm *PBM) Save(filename string) error {
	if pbm == nil {
		return errors.New("cannot save a nil PBM")
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write magic number, width, and height
	fmt.Fprintf(file, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	// Write image data
	for _, row := range pbm.data {
		for _, pixel := range row {
			if pixel {
				fmt.Fprint(file, "1")
			} else {
				fmt.Fprint(file, "0")
			}
		}
		fmt.Fprintln(file) // Add a newline after each row
	}

	return nil
}

// Invert inverts the colors of the PBM image.
func (pbm *PBM) Invert() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			pbm.data[i][j] = !pbm.data[i][j]
		}
	}
}

// Flip flips the PBM image horizontally.
func (pbm *PBM) Flip() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width/2; j++ {
			pbm.data[i][j], pbm.data[i][pbm.width-j-1] = pbm.data[i][pbm.width-j-1], pbm.data[i][j]
		}
	}
}

// Flop flops the PBM image vertically.
func (pbm *PBM) Flop() {
	for i := 0; i < pbm.height/2; i++ {
		pbm.data[i], pbm.data[pbm.height-i-1] = pbm.data[pbm.height-i-1], pbm.data[i]
	}
}

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}

func main() {
	// Example usage
	pbm, err := ReadPBM("p1.pbm")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	pbm.Save("original_output.pbm") // Save original image

	pbm.Invert()                    // Invert colors
	pbm.Save("inverted_output.pbm") // Save inverted image

	pbm.Flip()                     // Flip horizontally
	pbm.Save("flipped_output.pbm") // Save flipped image

	pbm.Flop()                     // Flop vertically
	pbm.Save("flopped_output.pbm") // Save flopped image
}
