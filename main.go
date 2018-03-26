package main

import (
	"shit/solver"
	"fmt"
	"shit/models"
)

func buildRoute(workstations solver.Tiles, storage *solver.Tile) solver.Options {
	var options solver.Options
	routes := solver.GetPermutation(workstations)

	for _, route := range routes {
		var totalCost float64
		route = append(route, storage)
		route = append([]*solver.Tile{storage}, route...)

		for i := 1; i < len(route); i++ {
			_, cost, found := solver.Path(route[i], route[i-1])

			if !found {
				fmt.Println("Cant't find the route")
			} else {
				totalCost += cost
			}
		}
		options = options.Append(&solver.Option{route, totalCost})

	}
	return options
}

func main() {
	models.LoadMaterials()
	qwe := models.LoadWorkstations()
	t := models.LoadTrain()

	t.LoadFromStorage(qwe)
	w1 := qwe[0]

	fmt.Println(t.CurrentCapacity)
	t.Unload(w1)
	fmt.Println(t.CurrentCapacity)


	//harita := solver.ParseMap(utils.GetMaze())

	//storage := harita.GetKind(solver.Start)[0]
	//routes := buildRoute(workstations, storage)
	//routes.ShowAllResults()

}
