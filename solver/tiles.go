package solver

import "github.com/gitchander/permutation"

// Tiles object contain list of Tile objects
type Tiles []*Tile

// Returns the number of how many Tile object that Tiles object contations
func (t Tiles) Len() int {
	return len(t)
}

func (t Tiles) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// Returns the permuation results of Tiles object to check the all possible route and
// calculate the cost
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
