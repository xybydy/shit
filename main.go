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

var usedStations models.Workstations
var availableWorkstations models.Workstations
var collectedWorkstations models.Workstations

var globalTime float64

type cheaper struct {
	ws   *models.Workstation
	cost float64
}

func buildRoute(workstations solver.Tiles, startPoint *solver.Tile, returning bool) solver.Options {
	var options solver.Options
	routes := solver.GetPermutation(workstations)

	for _, route := range routes {
		var totalCost float64
		var innerCosts []float64

		paths := make([][]pather.Pather, 0)
		route = append([]*solver.Tile{startPoint}, route...)
		if !returning {
			route = append(route, startPoint)
		}

		for i := 1; i < len(route); i++ {
			path, cost, found := pather.Path(route[i], route[i-1])

			if !found {
				fmt.Println("Cant't find the route")
			} else {
				cost = cost * 1 / trainModel.Speed
				totalCost += cost
				paths = append(paths, path)
				innerCosts = append(innerCosts, cost)
			}
		}
		options = options.Append(&solver.Option{route, totalCost, paths, innerCosts})

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
	options = buildRoute(workstations, train, false)

	bestOption = options.GetBestResult()

	printRoute(*bestOption)

}

func printRoute(o solver.Option) {
	fmt.Printf("Dene: %v", o.Cost)

	// fmt.Printf("Train Specs:\n-----------------\nLocation: %d,%d\nCapacity: %d\n\n", trainModel.X, trainModel.Y, trainModel.MaxCapacity)

	// fmt.Printf("Workstations\n-----------------\n")

	// for i := 0; i < len(workstationsModels); i++ {
	// 	var mat_sep = workstationsModels[i].PrintRequirements()
	//
	// 	fmt.Printf("\nName: %s\nLocation: %d,%d\nProcess Time: %d\nLoad Time: %d\nUnload Time: %d\nMaterials Demand:\n%s\n-----------------\n",
	// 		workstationsModels[i].Name, workstationsModels[i].X, workstationsModels[i].Y,
	// 		workstationsModels[i].Speed, workstationsModels[i].LoadTime, workstationsModels[i].UnloadTime, mat_sep)
	// }

	// fmt.Printf("Best: %.2f", o.Cost)

	DeliverStuff(o)
}

func DeliverStuff(o solver.Option) {
	var (
		totalPathCost   float64
		totalUnloadCost float64
		totalLoadCost   float64
	)

	for i := 0; i < len(o.Path); i++ {
		var (
			pathCost   float64
			unloadCost float64
			loadCost   float64
		)

		if i == 0 {
			to := o.Path[i][len(o.Path[i])-1].(*solver.Tile).Get().(*models.Workstation)

			usedStations = append(usedStations, to)

			pathCost = o.InnerCosts[i]
			loadCost = to.LoadTime
			totalCost := pathCost + loadCost

			fmt.Printf("\nFrom starting point to %s\n", to.Name)
			fmt.Printf("Train Stock: %s\n", trainModel.Stock.Details())

			trainModel.Unload(to, globalTime+totalCost)

			fmt.Printf("\nWarehouse demands:\n%s", to.PrintRequirements())
			fmt.Printf("Load Time: %.2f\n", loadCost)
			fmt.Printf("Time to reach: %.2f\n", pathCost)
			fmt.Printf("Total cost to deliver the product: %.2f", totalCost)

			harita.PrintMap(o.Path[i])

		} else if i == len(o.Path)-1 {
			from := o.Path[i][0].(*solver.Tile).Get().(*models.Workstation)

			fmt.Printf("\nTotal delivery cost: %.2f\n", totalPathCost)
			fmt.Printf("Total loading cost: %.2f\n", totalLoadCost)
			fmt.Printf("Total delivery cost: %.2f\n", totalPathCost+totalLoadCost)

			fmt.Printf("of %v", from.Name)
			collectAll(from)

			pathCost += o.InnerCosts[i]

			fmt.Printf("From %s back to storage\n", from.Name)
			fmt.Printf("Time to reach: %.2f\n", pathCost)

			harita.PrintMap(o.Path[i])

		} else {
			from := o.Path[i][0].(*solver.Tile).Get().(*models.Workstation)
			to := o.Path[i][len(o.Path[i])-1].(*solver.Tile).Get().(*models.Workstation)

			usedStations = append(usedStations, to)

			pathCost = o.InnerCosts[i]
			loadCost = to.LoadTime
			totalCost := pathCost + loadCost

			fmt.Printf("\nFrom %s to %s\n", from.Name, to.Name)
			fmt.Printf("Train Stock: %s\n", trainModel.Stock.Details())

			trainModel.Unload(to, globalTime+totalCost)

			fmt.Printf("\nWarehouse demands:\n%s", to.PrintRequirements())
			fmt.Printf("Load Time: %.2f\n", loadCost)
			fmt.Printf("Time to reach: %.2f\n", pathCost)
			fmt.Printf("Total cost to deliver the product: %.2f", totalCost)

			harita.PrintMap(o.Path[i])

		}

		totalPathCost += pathCost
		totalLoadCost += loadCost
		totalUnloadCost += unloadCost

		globalTime += pathCost + loadCost

	}

}

func collectAll(startPoint *models.Workstation) *models.Workstation {
	var totalPath, totalUnload float64
	var pathCost, unLoadCost float64

	fmt.Printf("\nCollection from warehouses\n")

	station := startPoint

	for i := 0; i < len(usedStations); i++ {
		checkAvailableWorkstation(station)
		nextStation := getWorkstationAvailable(station, availableWorkstations, true)

		pathCost, unLoadCost = collectOne(station, nextStation)

		totalPath += pathCost
		totalUnload += unLoadCost
		globalTime += pathCost + unLoadCost

		if !models.IsIn(nextStation, collectedWorkstations, false) {
			collectedWorkstations = append(collectedWorkstations, nextStation)
		}

		station = nextStation
	}
	return station

}

func collectOne(from, to *models.Workstation) (float64, float64) {
	unloadCost := to.UnloadTime
	path, pathCost, _ := pather.Path(harita.GetTile(from.X, from.Y), harita.GetTile(to.X, to.Y))

	fmt.Printf("\nCollecting from %s to %s\n", from.Name, to.Name)
	fmt.Printf("Time to reach: %.2f\n", pathCost)
	fmt.Printf("Workstation Unload Time: %.2f\n", unloadCost)
	fmt.Printf("Total time spent: %.2f\n", unloadCost+pathCost)
	harita.PrintMap(path)

	return pathCost, unloadCost

}

func checkAvailableWorkstation(startPoint *models.Workstation) {
	for _, workstation := range usedStations {
		if workstation.GetReadyTime() <= globalTime && !models.IsIn(workstation, availableWorkstations, false) && !models.IsIn(workstation, collectedWorkstations, false) {
			availableWorkstations = append(availableWorkstations, workstation)
		}
	}

	if len(availableWorkstations) <= len(collectedWorkstations) {
		availableWorkstations = append(availableWorkstations, getWorkstationAvailable(startPoint, usedStations, false))

	}

}

func getWorkstationAvailable(startPoint *models.Workstation, workstations models.Workstations, available bool) *models.Workstation {
	var cheap cheaper

	for _, workstation := range workstations {
		if !models.IsIn(workstation, availableWorkstations, available) {
			_, pathCost, _ := pather.Path(harita.GetTile(startPoint.X, startPoint.Y), harita.GetTile(workstation.X, workstation.Y))
			cost := pathCost + workstation.GetReadyTime() + globalTime

			fmt.Printf("qqqq %v - %v", workstation.Name, available)

			if cheap.ws == nil {
				cheap = cheaper{ws: workstation, cost: cost}
				fmt.Printf("getwork: %v - %v - %v - %v\n", workstation.Name, workstation.GetReadyTime(), workstation.UnloadTime, cost)

			} else {
				fmt.Printf("getwork2: %v - %v - %v - %v\n", workstation.Name, workstation.GetReadyTime(), workstation.UnloadTime, cost)
				if cost < cheap.cost {
					cheap = cheaper{ws: workstation, cost: cost}
				}

			}

		} else {
			fmt.Printf("var %v\n", workstation.Name)
		}
	}
	return cheap.ws
}

// TODO filledtime workstation calismiyor.
