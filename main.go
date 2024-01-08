package netpbm

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}
