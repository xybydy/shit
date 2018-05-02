package main

import (
	"fmt"
	"strings"

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
var collectedWorkstations models.Workstations

var globalTime float64

type cheaper struct {
	ws   *models.Workstation
	cost float64
}

var border = strings.Repeat("-", 45)

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
	fmt.Printf("\n%s", border)

	for i := 0; i < len(o.Path); i++ {
		var (
			pathCost float64
			loadCost float64
		)

		fmt.Printf("\n%s\n", border)
		if i == 0 {
			to := o.Path[i][len(o.Path[i])-1].(*solver.Tile).Get().(*models.Workstation)

			usedStations = append(usedStations, to)

			pathCost = o.InnerCosts[i]
			loadCost = to.LoadTime
			totalCost := pathCost + loadCost

			fmt.Printf("Delivering from warehouse point to %s\n\n", to.Name)
			fmt.Printf("Available Stock: %s\n", trainModel.Stock.Details())

			trainModel.Unload(to, globalTime+totalCost)

			fmt.Printf("\nWarehouse demands:\n%s", to.PrintRequirements())
			fmt.Printf("\nLoad Time: %.2f\n", loadCost)
			fmt.Printf("Time to reach: %.2f\n", pathCost)
			fmt.Printf("Total time to deliver the product: %.2f\n\n", totalCost)

			fmt.Printf("The route on the map")
			harita.PrintMap(o.Path[i])

		} else if i == len(o.Path)-1 {
			from := o.Path[i][0].(*solver.Tile).Get().(*models.Workstation)

			fmt.Printf("\nALL MATERIALS ARE DELIVERED!\n", totalPathCost)
			fmt.Printf("\nTotal delivery cost: %.2f\n", totalPathCost)
			fmt.Printf("Total loading cost: %.2f\n", totalLoadCost)
			fmt.Printf("Total delivery cost: %.2f\n", totalPathCost+totalLoadCost)

			from, unloadPathCost, UnloadCost := collectAll(from)

			path, cost, _ := pather.Path(harita.GetTile(from.X, from.Y), train)

			pathCost += cost

			fmt.Printf("\n\n%s\n", border)
			fmt.Printf("\nALL MATERIALS ARE COLLECTED!\n\n")
			fmt.Printf("Total unload time: %v\n", unloadPathCost)
			fmt.Printf("Total time to return time: %v\n", unloadPathCost+pathCost)
			fmt.Printf("\n%s\n", border)

			fmt.Printf("\nFrom %s back to storage\n\n", from.Name)
			fmt.Printf("Time to reach: %.2f\n", pathCost)
			harita.PrintMap(path)

			pathCost += unloadPathCost
			totalUnloadCost += UnloadCost

		} else {
			from := o.Path[i][0].(*solver.Tile).Get().(*models.Workstation)
			to := o.Path[i][len(o.Path[i])-1].(*solver.Tile).Get().(*models.Workstation)

			usedStations = append(usedStations, to)

			pathCost = o.InnerCosts[i]
			loadCost = to.LoadTime
			totalCost := pathCost + loadCost

			fmt.Printf("\nFrom %s to %s\n\n", from.Name, to.Name)
			fmt.Printf("Train Stock: %s\n", trainModel.Stock.Details())

			trainModel.Unload(to, globalTime+totalCost)

			fmt.Printf("\nWarehouse demands:\n%s", to.PrintRequirements())
			fmt.Printf("\nLoad Time: %.2f\n", loadCost)
			fmt.Printf("Time to reach: %.2f\n", pathCost)
			fmt.Printf("Total time to deliver the product: %.2f\n\n", totalCost)

			fmt.Printf("The route on the map")
			harita.PrintMap(o.Path[i])

		}

		totalPathCost += pathCost
		totalLoadCost += loadCost
		globalTime += pathCost + loadCost + totalUnloadCost

	}
	fmt.Printf("---------------------------------------\n")
	fmt.Printf("Total simulation time: %v\n", globalTime)

}

func collectAll(startPoint *models.Workstation) (*models.Workstation, float64, float64) {

	var pathCost, unLoadCost, idleTime float64

	fmt.Printf("\n%s\n", border)
	fmt.Printf("%s\n", border)
	fmt.Printf("\nCOLLECTING FROM WORKSTATIONS\n")

	station := startPoint

	for i := 0; i < len(usedStations); i++ {
		nextStation := getWorkstationAvailable(station, usedStations)

		fmt.Printf("Test: %v", nextStation.GetReadyTime())

		pathCost, unLoadCost, idleTime = collectOne(station, nextStation)

		globalTime += pathCost + unLoadCost + idleTime

		if !models.IsIn(nextStation, collectedWorkstations, false) {
			collectedWorkstations = append(collectedWorkstations, nextStation)
		}

		station = nextStation
	}

	return station, pathCost, unLoadCost

}

func collectOne(from, to *models.Workstation) (float64, float64, float64) {
	var idleTime float64
	unloadCost := to.UnloadTime
	path, pathCost, _ := pather.Path(harita.GetTile(from.X, from.Y), harita.GetTile(to.X, to.Y))

	if to.GetReadyTime() < globalTime+pathCost {
		idleTime = 0.0
	} else {
		idleTime = to.GetReadyTime() - globalTime - pathCost
	}

	fmt.Printf("\nCollecting from %s to %s\n", from.Name, to.Name)
	fmt.Printf("\nWorkstation %s will be ready at: %.2f\n", to.Name, to.GetReadyTime())
	fmt.Printf("\nTrain idle time %.2f\n", idleTime)
	fmt.Printf("Time to reach: %.2f\n", pathCost)
	fmt.Printf("Workstation Unload Time: %.2f\n", unloadCost)
	fmt.Printf("Total time spent: %.2f\n", unloadCost+pathCost+idleTime)
	harita.PrintMap(path)

	return pathCost, unloadCost, idleTime

}

func getWorkstationAvailable(startPoint *models.Workstation, workstations models.Workstations) *models.Workstation {
	var cheap cheaper
	var topPriorityWS []cheaper
	var lowPriorityWS []cheaper

	for _, workstation := range workstations {
		if !models.IsIn(workstation, collectedWorkstations, false) {
			_, pathCost, _ := pather.Path(harita.GetTile(startPoint.X, startPoint.Y), harita.GetTile(workstation.X, workstation.Y))
			cost := pathCost + workstation.GetReadyTime() + globalTime
			topPriorityWS = append(topPriorityWS, cheaper{workstation, cost})
		}
	}

	for _, workstation := range workstations {
		if !models.IsIn(workstation, collectedWorkstations, false) {
			_, pathCost, _ := pather.Path(harita.GetTile(startPoint.X, startPoint.Y), harita.GetTile(workstation.X, workstation.Y))
			cost := pathCost + workstation.GetReadyTime() + globalTime

			if workstation.GetReadyTime() < globalTime {
				topPriorityWS = append(topPriorityWS, cheaper{workstation, cost})
			} else {
				lowPriorityWS = append(lowPriorityWS, cheaper{workstation, cost})
			}

		}
	}

	if len(topPriorityWS) != 0 {
		for _, workstation := range topPriorityWS {

			if cheap.ws == nil {
				cheap = workstation
			} else {
				if workstation.cost < cheap.cost {
					cheap = workstation
				}
			}
		}

	} else {
		for _, workstation := range lowPriorityWS {
			if cheap.ws == nil {
				cheap = workstation
			} else {
				if workstation.cost < cheap.cost {
					cheap = workstation
				}
			}
		}
	}

	return cheap.ws
}

// TODO filledtime workstation calismiyor.
