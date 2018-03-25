package solver

import "fmt"

type Options []*Option

type Option struct {
	Route Tiles
	Cost  float64
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
		}
	}
	return bestOption
}

func (o Options) ShowAllResults() {
	for _, option := range o {
		fmt.Printf("Cost: %v, Route: %v\n", option.Cost, option.Route)
	}
}
