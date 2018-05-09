package solver

import (
	"semesterproject/pather"
)

// The list of `Option` objects
type Options []*Option

// Each `Option` object holds whole the route options throughout the all workstations.
type Option struct {
	// The `Tiles` that the `Train` object went through the route
	Route Tiles
	// Number of `Tiles` that `Train` object went through, each `Tile` costs 1 unit time.
	Cost float64
	// Sub-path in the `Route` from each point to next point
	Path [][]pather.Pather
	// The costs of inner paths
	InnerCosts []float64
}

// Add `Option` object to `Options` list
func (o Options) Append(op ...*Option) Options {
	for _, k := range op {
		if k.Cost != 0 {
			o = append(o, k)
		}
	}
	return o
}

// Returns the Option with lowest cost - faster - route.
func (o Options) GetBestResult() *Option {
	var bestOption = &Option{}

	for i, option := range o {
		if i == 0 {
			bestOption = option
		}
		if option.Cost < bestOption.Cost {
			bestOption = option
		} else if option.Cost == bestOption.Cost {
			if option.InnerCosts[0] < bestOption.InnerCosts[0] {
				bestOption = option
			}
		}
	}
	return bestOption
}

func (o Option) Len() int {
	return len(o.Route)
}
