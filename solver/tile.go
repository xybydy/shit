package solver

import (
	"github.com/gitchander/permutation"
)

type Tiles []*Tile

func (t Tiles) Len() int {
	return len(t)
}

func (t Tiles) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

type Tile struct {
	Kind int
	X    int
	Y    int
	W    Map
	Base interface{}
}

func (t *Tile) PathNeighbors() []Pather {
	var neighbors []Pather
	for _, offset := range [][]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	} {
		if n := t.W.Tile(t.X+offset[0], t.Y+offset[1]); n != nil && n.Kind != Wall {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

func (t *Tile) PathNeighborCost(to Pather) float64 {
	toT := to.(*Tile)
	return Costs[toT.Kind]
}

func (t *Tile) PathEstimatedCost(to Pather) float64 {
	toT := to.(*Tile)
	absX := toT.X - t.X
	if absX < 0 {
		absX = -absX
	}
	absY := toT.Y - t.Y
	if absY < 0 {
		absY = -absY
	}
	return float64(absX + absY)
}

func (t *Tile) SetBase(){
	

}





func GetPermutation(tiles Tiles) []Tiles {
	var routes []Tiles

	p := permutation.New(tiles)
	for p.Scan() {
		q := make(Tiles, len(tiles))
		copy(q, tiles)
		routes = append(routes, q)
	}

	return routes

}
