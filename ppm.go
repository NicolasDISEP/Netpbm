package Netpbm

import (
	"bufio"
	"fmt"
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

func (ppm *PPM) PrintPPM() {
	fmt.Printf("Magic Number: %s\n", ppm.magicNumber)
	fmt.Printf("Width: %d\n", ppm.width)
	fmt.Printf("Height: %d\n", ppm.height)
	fmt.Printf("Max Value: %d\n", ppm.max)

	fmt.Println("Pixel Data:")
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			pixel := ppm.data[y][x]
			fmt.Printf("(%d, %d, %d) ", pixel.R, pixel.G, pixel.B)
		}
		fmt.Println()
	}
}

func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

func (ppm *PPM) At(x, y int) Pixel {
	// Vérification des limites pour éviter les erreurs d'index
	if x < 0 || x >= ppm.width || y < 0 || y >= ppm.height {
		// Vous pouvez également gérer cela différemment, comme renvoyer une valeur par défaut ou une erreur.
		panic("Index out of bounds")
	}

	return ppm.data[y][x]
}

func (ppm *PPM) Set(x, y int, value Pixel) {
	// Vérification des limites pour éviter les erreurs d'index
	if x < 0 || x >= ppm.width || y < 0 || y >= ppm.height {
		// Vous pouvez également gérer cela différemment, comme renvoyer une valeur par défaut ou une erreur.
		panic("Index out of bounds")
	}

	ppm.data[y][x] = value
}

func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if ppm.magicNumber == "P6" || ppm.magicNumber == "P3" {
		fmt.Fprintf(file, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)
	} else {
		err = fmt.Errorf("magic number error")
		return err
	}
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			pixel := ppm.data[y][x]
			if ppm.magicNumber == "P6" {
				fmt.Fprintf(file, "%c%c%c", pixel.R, pixel.G, pixel.B)
			} else if ppm.magicNumber == "P3" {
				fmt.Fprintf(file, "%d %d %d ", pixel.R, pixel.G, pixel.B)
			}
		}
		if ppm.magicNumber == "P3" {
			fmt.Fprint(file, "\n")
		}
	}

	return nil
}

func (ppm *PPM) Invert() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			pixel := &ppm.data[y][x]
			pixel.R = 255 - pixel.R
			pixel.G = 255 - pixel.G
			pixel.B = 255 - pixel.B
		}
	}
}

func (ppm *PPM) Flip() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width/2; x++ {
			ppm.data[y][x], ppm.data[y][ppm.width-x-1] = ppm.data[y][ppm.width-x-1], ppm.data[y][x]
		}
	}
}

func (ppm *PPM) Flop() {
	for y := 0; y < ppm.height/2; y++ {
		ppm.data[y], ppm.data[ppm.height-y-1] = ppm.data[ppm.height-y-1], ppm.data[y]
	}
}

func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

func (ppm *PPM) SetMaxValue(maxValue uint8) {
	ppm.max = int(maxValue)
}
func (ppm *PPM) Rotate90CW() {
	newPPM := PPM{
		data:        make([][]Pixel, ppm.width),
		width:       ppm.height,
		height:      ppm.width,
		magicNumber: ppm.magicNumber,
		max:         ppm.max,
	}

	for i := range newPPM.data {
		newPPM.data[i] = make([]Pixel, newPPM.width)
	}

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			newPPM.data[x][ppm.height-y-1] = ppm.data[y][x]
		}
	}

	*ppm = newPPM
}

func (ppm *PPM) ToPGM() *PGM {
	pgm := &PGM{
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P2",
		max:         ppm.max,
	}

	pgm.data = make([][]uint8, ppm.height)
	for i := range pgm.data {
		pgm.data[i] = make([]uint8, ppm.width)
	}

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			gray := rgbToGray(ppm.data[y][x])
			pgm.data[y][x] = gray
		}
	}

	return pgm
}

type Point struct {
	X, Y int
}

// rgbToGray converts an RGB color to a grayscale value.
func rgbToGray(color Pixel) uint8 {
	// Use luminosity method for converting RGB to grayscale
	// Gray = 0.299*R + 0.587*G + 0.114*B
	return uint8(0.299*float64(color.R) + 0.587*float64(color.G) + 0.114*float64(color.B))
}

func (ppm *PPM) ToPBM() *PBM {
	pbm := &PBM{
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P1",
	}

	pbm.data = make([][]bool, ppm.height)
	for i := range pbm.data {
		pbm.data[i] = make([]bool, ppm.width)
	}

	threshold := uint8(ppm.max / 2)
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			average := (uint8(ppm.data[y][x].R) + uint8(ppm.data[y][x].G) + uint8(ppm.data[y][x].B)) / 3
			pbm.data[y][x] = average > uint8(threshold)
		}
	}

	return pbm
}
