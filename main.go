package main

import (
	"shit/solver"
	"fmt"
)

var finalRoutes []solver.Tiles

func buildRoute(workstations solver.Tiles, storage *solver.Tile) []solver.Tiles {
	routes := solver.GetPermutation(workstations)

	for _, route := range routes {
		route = append(route, storage)
		route = append([]*solver.Tile{storage}, route...)
		finalRoutes = append(finalRoutes, route)
	}

	return finalRoutes

}

func main() {

	harita := solver.ParseMap(`
################
#S.............#
#.###.########.#
#.##W.########.#
#.###.#####.##.#
#.###.......W#.#
#.###.###.#....#
#.###.###.#.##.#
#.###.W##.#.##.#
#.###.###.#.##.#
#.###.###.#.##.#
#.###.###.#..#.#
#.........W.##.#
#.#W#####......#
#.##############
################
	`)

	storage := harita.GetKind(solver.Start)[0]
	workstations := harita.GetKind(solver.Workstation)

	routes := buildRoute(workstations, storage)

	for _, route := range routes {
		//var totalCost float64
		//sols := make(map[float64][]solver.Pather)

		for i := 1; i < len(route); i++ {
			p, _, found := solver.Path(route[i], route[i-1])

			if !found {
				fmt.Println("Cant't find the route")
			} else {

				for _, j := range p {
					k := j.(*solver.Tile)
					fmt.Printf("%d: %d.%d\n", k.Kind, k.X, k.Y)
				}

				//for _, k := range p {
				//	pT := k.(*solver.Tile)
				//	fmt.Printf("%d.%d, ", pT.X, pT.Y)
				//}

			}

		}
		fmt.Println("qeweqe")

	}

}
