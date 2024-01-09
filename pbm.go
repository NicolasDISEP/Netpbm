package netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// ReadPBM reads a PBM image from a file and returns a struct that represents the image.
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	cleanedData := removeComments(scanner)

	// Now, create a new scanner for the cleaned content
	scanner = bufio.NewScanner(strings.NewReader(cleanedData))

	// Read magic number
	scanner.Scan()
	magicNumber := scanner.Text()

	// Read width and height
	scanner.Scan()
	width, height := 0, 0
	fmt.Sscanf(scanner.Text(), "%d %d", &width, &height)

	// Read image data
	data := make([][]bool, height)
	if magicNumber == "P1" {
		// P1 format (ASCII)
		for i := 0; i < height && scanner.Scan(); i++ {
			line := strings.Fields(scanner.Text()) // Split the line into fields
			data[i] = make([]bool, width)
			for j, field := range line {
				if field == "1" {
					data[i][j] = true
				} else if field == "0" {
					data[i][j] = false
				} else {
					return nil, fmt.Errorf("Invalid character in image data")
				}
			}
		}

	} else if magicNumber == "P4" {
		// For P4 (binary), read bytes
		for i := 0; i < height; i++ {
			scanner.Scan()
			line := scanner.Text()
			data[i] = make([]bool, width)
			for j := 0; j < width; j++ {
				byteIndex := j / 8
				bitIndex := 7 - (j % 8)
				bit := (line[byteIndex] >> uint(bitIndex)) & 1
				data[i][j] = bit == 1
			}
		}
	} else {
		return nil, fmt.Errorf("Invalid magic number")
	}

	return &PBM{data, width, height, magicNumber}, nil
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

// removeComments removes comments (lines starting with '#') from the scanner.
func removeComments(scanner *bufio.Scanner) string {
	var cleanedLines []string

	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line starts with '#'
		if strings.HasPrefix(line, "#") {
			// Ignore the comment line
			continue
		}
		cleanedLines = append(cleanedLines, line)
	}

	// Return cleaned lines as a single string
	return strings.Join(cleanedLines, "\n")
}

/*
func main() {
	// Example usage
	pbm, err := ReadPBM("p4.pbm")
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
}*/
