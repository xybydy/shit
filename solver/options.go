package solver

import (
	"shit/pather"
)

type Options []*Option

type Option struct {
	Route      Tiles
	Cost       float64
	Path       [][]pather.Pather
	InnerCosts []float64
}

func (o Options) Append(op ...*Option) Options {
	for _, k := range op {
		if k.Cost != 0 {
			o = append(o, k)
		}
	}
	return o
}

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
