package main

import (
	"fmt"

	"shit/models"
	"shit/pather"
	"shit/solver"
	"shit/utils"
)

var harita solver.Map
var train *solver.Tile
var workstations solver.Tiles
var options solver.Options
var bestOption *solver.Option

var trainModel models.Train
var workstationsModels models.Workstations

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
	models.LoadMaterials()
	workstationsModels = models.LoadWorkstations()
	trainModel = models.LoadTrain()
	trainModel.LoadFromStorage(workstationsModels)
	getResult()

}

func getResult() {
	harita = solver.ParseMap(utils.GetMaze())
	train = harita.GetKind(solver.Train)[0]
	workstations = harita.GetKind(solver.Workstation)
	options = buildRoute(workstations, train)

	bestOption = options.GetBestResult()

	printRoute(*bestOption)

}

func printRoute(o solver.Option) {
	totalCost := 0

	fmt.Printf("Train Specs:\n-----------------\nLocation: %d,%d\nCapacity: %d\n\n", trainModel.X, trainModel.Y, trainModel.MaxCapacity)

	fmt.Printf("Workstations\n-----------------\n")

	for i := 0; i < len(workstationsModels); i++ {
		var mat_sep = workstationsModels[i].PrintRequirements()

		fmt.Printf("\nName: %s\nLocation: %d,%d\nProcess Time: %d\nLoad Time: %d\nUnload Time: %d\nMaterials Demand:\n%s\n-----------------\n",
			workstationsModels[i].Name, workstationsModels[i].X, workstationsModels[i].Y,
			workstationsModels[i].Speed, workstationsModels[i].LoadTime, workstationsModels[i].UnloadTime, mat_sep)
	}

	for i := 0; i < len(o.Path); i++ {
		cost := 0
		if i == 0 {
			to := o.Path[i][len(o.Path[i])-1].(*solver.Tile).Get().(*models.Workstation)

			fmt.Printf("\nFrom starting point to %s\n", to.Name)

			fmt.Printf("Train Stock: %s\n", trainModel.Stock.Details())

			fmt.Printf("\nWarehouse demands:\n%s", to.PrintRequirements())
			fmt.Printf("Load Time: %d\n", to.LoadTime)
			cost += trainModel.Unload(to)

			harita.PrintMap(o.Path[i])

		} else if i == len(o.Path)-1 {
			from := o.Path[i][0].(*solver.Tile).Get().(*models.Workstation)

			fmt.Printf("From %s back to storage\n", from.Name)

			harita.PrintMap(o.Path[i])

		} else {
			from := o.Path[i][0].(*solver.Tile).Get().(*models.Workstation)
			to := o.Path[i][len(o.Path[i])-1].(*solver.Tile).Get().(*models.Workstation)
			fmt.Printf("\nFrom %s to %s\n", from.Name, to.Name)
			fmt.Printf("Train Stock: %s\n", trainModel.Stock.Details())
			fmt.Printf("\nWarehouse demands:\n%s", to.PrintRequirements())
			fmt.Printf("Load Time: %d\n", to.LoadTime)
			trainModel.Unload(to)

			harita.PrintMap(o.Path[i])

		}

	}
}

// TODO sure carpanlari kullanip ne kadar vakit harcandigi artik raporlansin.
// TODO geri donsun, hatta toplayarak donsun.
