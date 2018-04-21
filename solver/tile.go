package solver

import (
	"fmt"

	"shit/models"
	"shit/pather"
)

type Tile struct {
	Kind int
	X    int
	Y    int
	W    Map
}

func (t *Tile) PathNeighbors() []pather.Pather {
	var neighbors []pather.Pather
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

func (t *Tile) PathNeighborCost(to pather.Pather) float64 {
	toT := to.(*Tile)
	return Costs[toT.Kind]
}

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

func (t *Tile) Get() interface{} {
	if t.Kind == Train {
		train := models.LoadTrain()
		r := models.GetTrain(t.X, t.Y, train)
		if r == nil {
			fmt.Println("Train does not match")
			return nil
		} else {
			return r
		}
	} else if t.Kind == Workstation {
		stations := models.LoadWorkstations()
		r := models.GetWorkstation(t.X, t.Y, stations)
		if r == nil {
			fmt.Println(t.X, t.Y, " does not match")
			return nil
		} else {
			return r
		}

	}
	return nil

}
