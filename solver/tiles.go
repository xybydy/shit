package solver

import "github.com/gitchander/permutation"

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

type Tiles []*Tile

func (t Tiles) Len() int {
	return len(t)
}

func (t Tiles) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
