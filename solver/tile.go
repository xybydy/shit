package solver

import (
	"fmt"

	"semesterproject/models"
	"semesterproject/pather"
)

// Tile object in the Map.
// This point in the map is an `Tile` object.
type Tile struct {
	// The type of the `Tile` object
	Kind int
	// X coordination of `Tile`
	X int
	// Y coordination of `Tile`
	Y int
	// The `Map` object which contains the `Tile` object
	W Map
}

// Returns the available - non-wall - neighbors of the `Tile`
func (t *Tile) PathNeighbors() []pather.Pather {
	var neighbors []pather.Pather
	for _, offset := range [][]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	} {
		if n := t.W.GetTile(t.X+offset[0], t.Y+offset[1]); n != nil && n.Kind != Wall {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

// Returns the cost of the road to Neighbor.
func (t *Tile) PathNeighborCost(to pather.Pather) float64 {
	toT := to.(*Tile)
	return Costs[toT.Kind]
}

// Basic implementation to calculate the estimated path cost.
func (t *Tile) PathEstimatedCost(to pather.Pather) float64 {
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

// Returns the related Object of Tile.
func (t *Tile) Get() interface{} {
	if t.Kind == Train {
		train := models.LoadTrain()
		r := models.GetTrain(t.X, t.Y, train)
		if r == nil {
			fmt.Println("Train does not match")
			return nil
		}
		return r

	} else if t.Kind == Workstation {
		stations := models.LoadWorkstations()
		r := models.GetWorkstation(t.X, t.Y, stations)
		if r == nil {
			fmt.Println(t.X, t.Y, " does not match")
			return nil
		}
		return r

	}
	return nil
}
