package main

import (
	"fmt"

	"shit/pather"
	"shit/solver"
	"shit/utils"
)

var harita solver.Map
var train *solver.Tile
var workstations solver.Tiles

func buildRoute(workstations solver.Tiles, train *solver.Tile) solver.Options {
	var options solver.Options
	routes := solver.GetPermutation(workstations)

	for _, route := range routes {
		var totalCost float64
		paths := make([][]pather.Pather, 0)
		route = append(route, train)
		route = append([]*solver.Tile{train}, route...)

		for i := 1; i < len(route); i++ {
			path, cost, found := pather.Path(route[i], route[i-1])

			if !found {
				fmt.Println("Cant't find the route")
			} else {
				totalCost += cost
				paths = append(paths, path)
			}
		}
		options = options.Append(&solver.Option{route, totalCost, paths})

	}
	return options
}

func main() {
	//models.LoadMaterials()
	//qwe := models.LoadWorkstations()
	//t := models.LoadTrain()
	//
	//t.LoadFromStorage(qwe)
	//w1 := qwe[0]
	//
	//fmt.Println(t.CurrentCapacity)
	//t.Unload(w1)
	//fmt.Println(t.CurrentCapacity)

	//alt kismi map
	getResult()

}

func getResult() {
	harita = solver.ParseMap(utils.GetMaze())
	train = harita.GetKind(solver.Train)[0]
	workstations = harita.GetKind(solver.Workstation)
	routes := buildRoute(workstations, train)

	bestRoute := routes.GetBestResult()

	printRoute(*bestRoute)

}

func printRoute(o solver.Option) {

	fmt.Printf("The best route's cost is %v\n", o.Cost)

	harita.PrintMap(o.Path)

}
