/*
We prepare our map in `plant.map` file, then go to inputs folder to define our models such as workstation, train and materials.
Then program loads `plant.map` and the model files first.
Afterwards the simulation does is checking coordinations in given model files if they're in-line with the map.
If they are not in-line, simulation stops running.
After checking the models, simulation finds all possible ways from warehouse stops at through all workstations
and then back to warehouse again with calculating the cost of the route. Amongst the all routes the program picks the shortest route,
Then bring all the output into "output.txt" file.
*/
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"semesterproject/models"
	"semesterproject/pather"
	"semesterproject/solver"
	"semesterproject/utils"
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

var fileWriter io.Writer

type cheaper struct {
	ws   *models.Workstation
	cost float64
}

var border = strings.Repeat("-", 45)

// Gets all the possible routes within the map and returns the options as `Options` type.
func buildRoute(workstations solver.Tiles, startPoint *solver.Tile) solver.Options {
	var options solver.Options
	routes := solver.GetPermutation(workstations)

	for _, route := range routes {
		var totalCost float64
		var innerCosts []float64

		paths := make([][]pather.Pather, 0)
		route = append([]*solver.Tile{startPoint}, route...)

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
		options = options.Append(&solver.Option{Route: route, Cost: totalCost, Path: paths, InnerCosts: innerCosts})

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

// Returns the results of the simulation and prints to the user
func getResult() {
	harita = solver.ParseMap(utils.GetMaze())
	train = harita.GetKind(solver.Train)[0]
	workstations = harita.GetKind(solver.Workstation)
	options = buildRoute(workstations, train)

	if !harita.CrossCheck() {
		fmt.Println("Please check all locations in models and the map. Locations are not in-line")
		fmt.Println("Simulation stops.")
		os.Exit(1)
	}

	bestOption = options.GetBestResult()

	f, err := os.Create("output.txt")
	if err != nil {
		fmt.Errorf("couldn't create the file")
	}
	defer f.Close()

	fileWriter = io.MultiWriter(os.Stdout, f)

	printRoute(*bestOption)

}

func printRoute(o solver.Option) {

	row, col := harita.GetSize()

	fmt.Fprintf(fileWriter, "\n%s\nMAP SPECS:\n%s\n", border, border)
	fmt.Fprintf(fileWriter, "Map Size: %dx%d\n", row, col)
	fmt.Fprintf(fileWriter, "Number of Workstations: %d\n%s\n", len(workstationsModels), border)
	fmt.Fprintf(fileWriter, "\n%s\nTRAIN SPECS:\n%s\nLocation: %d,%d\nCapacity: %d\n%s\n\n", border, border, trainModel.X, trainModel.Y, trainModel.MaxCapacity, border)

	fmt.Fprintf(fileWriter, "%s\nWORKSTATIONS\n%s\n", border, border)

	for i := 0; i < len(workstationsModels); i++ {
		var matSep = workstationsModels[i].PrintRequirements()

		fmt.Fprintf(fileWriter, "Name: %s\nLocation: %d,%d\nProcess Time: %.2f\nLoad Time: %.2f\nUnload Time: %.2f\nMaterials Demand:\n%s%s\n",
			workstationsModels[i].Name, workstationsModels[i].X, workstationsModels[i].Y,
			workstationsModels[i].Speed, workstationsModels[i].LoadTime, workstationsModels[i].UnloadTime, matSep, border)
	}

	deliverStuff(o)
}

func deliverStuff(o solver.Option) {
	var (
		totalPathCost   float64
		totalUnloadCost float64
		totalLoadCost   float64
		totalIdleCost   float64
	)
	fmt.Fprintf(fileWriter, "\n\n%s", border)

	fmt.Fprintf(fileWriter, "\nSIMULATION REPORT\n%s\n", border)

	for i := 0; i < len(o.Path); i++ {
		var (
			pathCost float64
			loadCost float64
		)

		fmt.Fprintf(fileWriter, "\n%s\n", border)
		if i == 0 {
			to := o.Path[i][len(o.Path[i])-1].(*solver.Tile).Get().(*models.Workstation)

			usedStations = append(usedStations, to)

			pathCost = o.InnerCosts[i]
			loadCost = to.LoadTime
			totalCost := pathCost + loadCost

			fmt.Fprintf(fileWriter, "Delivering from warehouse point to %s\n\n", to.Name)
			fmt.Fprintf(fileWriter, "Available Stock: %s\n", trainModel.Stock.Details())

			trainModel.Unload(to, globalTime+totalCost)

			fmt.Fprintf(fileWriter, "\nWarehouse Demand:\n%s", to.PrintRequirements())
			fmt.Fprintf(fileWriter, "\nTime to load the material: %.2f\n", loadCost)
			fmt.Fprintf(fileWriter, "Time to reach to workstation: %.2f\n", pathCost)
			fmt.Fprintf(fileWriter, "Total time to deliver the product: %.2f\n\n", totalCost)

			fmt.Fprintf(fileWriter, "MAP PREVIEW\n")
			harita.PrintMap(fileWriter, o.Path[i])

		} else if i == len(o.Path)-1 {
			from := o.Path[i][0].(*solver.Tile).Get().(*models.Workstation)

			fmt.Fprintf(fileWriter, "\nALL MATERIALS ARE DELIVERED!\n")
			fmt.Fprintf(fileWriter, "\nTotal delivery cost: %.2f\n", totalPathCost)
			fmt.Fprintf(fileWriter, "Total loading cost: %.2f\n", totalLoadCost)
			fmt.Fprintf(fileWriter, "TOTAL COST: %.2f\n", globalTime)

			from, unloadPathCost, unloadCost, idleTime := collectAll(from)

			path, cost, _ := pather.Path(harita.GetTile(from.X, from.Y), train)

			pathCost += cost

			fmt.Fprintf(fileWriter, "\nFrom %s back to storage\n\n", from.Name)
			fmt.Fprintf(fileWriter, "Time to reach: %.2f\n", pathCost)
			fmt.Fprintf(fileWriter, "MAP PREVIEW\n")
			harita.PrintMap(fileWriter, path)

			pathCost += unloadPathCost
			totalUnloadCost += unloadCost
			totalIdleCost = idleTime

			fmt.Fprintf(fileWriter, "\n\n%s\n", border)
			fmt.Fprintf(fileWriter, "\nALL MATERIALS ARE COLLECTED AND RETURNED TO WAREHOUSE!\n\n")
			fmt.Fprintf(fileWriter, "Total return path cost: %v\n", unloadPathCost+cost)
			fmt.Fprintf(fileWriter, "Total unload cost: %v\n", unloadCost)
			fmt.Fprintf(fileWriter, "Total idle cost: %v\n", idleTime)
			fmt.Fprintf(fileWriter, "TOTAL COLLECTION COST: %v\n", idleTime+unloadCost+unloadPathCost)
			fmt.Fprintf(fileWriter, "\n%s\n", border)

		} else {
			from := o.Path[i][0].(*solver.Tile).Get().(*models.Workstation)
			to := o.Path[i][len(o.Path[i])-1].(*solver.Tile).Get().(*models.Workstation)

			usedStations = append(usedStations, to)

			pathCost = o.InnerCosts[i]
			loadCost = to.LoadTime
			totalCost := pathCost + loadCost

			fmt.Fprintf(fileWriter, "From %s to %s\n\n", from.Name, to.Name)
			fmt.Fprintf(fileWriter, "Train Stock: %s\n", trainModel.Stock.Details())

			trainModel.Unload(to, globalTime+totalCost)

			fmt.Fprintf(fileWriter, "\nWarehouse Demand:\n%s", to.PrintRequirements())
			fmt.Fprintf(fileWriter, "\nLoad Time: %.2f\n", loadCost)
			fmt.Fprintf(fileWriter, "Time to reach: %.2f\n", pathCost)
			fmt.Fprintf(fileWriter, "Total time to deliver the product: %.2f\n\n", totalCost)

			fmt.Fprintf(fileWriter, "MAP PREVIEW\n")
			harita.PrintMap(fileWriter, o.Path[i])

		}

		totalPathCost += pathCost
		totalLoadCost += loadCost
		globalTime += pathCost + loadCost + totalUnloadCost + totalIdleCost

	}
	fmt.Fprintf(fileWriter, "%s\n", border)
	fmt.Fprintf(fileWriter, "Total simulation cost: %v\n", globalTime)
	fmt.Fprintf(fileWriter, "Total path cost: %v\n", totalPathCost)
	fmt.Fprintf(fileWriter, "Total loading cost: %v\n", totalLoadCost)
	fmt.Fprintf(fileWriter, "Total unloading cost: %v\n", totalUnloadCost)
	fmt.Fprintf(fileWriter, "Total idle time cost: %v\n", totalIdleCost)
	fmt.Fprintf(fileWriter, "%s\n", border)

}

func collectAll(startPoint *models.Workstation) (*models.Workstation, float64, float64, float64) {

	var totalPathCost, totalUnLoadCost, totalIdleTime float64
	var pathCost, unLoadCost, idleTime float64

	fmt.Fprintf(fileWriter, "\n%s\n", border)
	fmt.Fprintf(fileWriter, "COLLECTION FROM WORKSTATIONS\n")
	fmt.Fprintf(fileWriter, "%s\n", border)

	station := startPoint

	for i := 0; i < len(usedStations); i++ {
		nextStation := getWorkstationAvailable(station, usedStations)

		pathCost, unLoadCost, idleTime = collectOne(station, nextStation)

		// globalTime += pathCost + unLoadCost + idleTime

		totalPathCost += pathCost
		totalUnLoadCost += unLoadCost
		totalIdleTime += idleTime

		if !models.IsIn(nextStation, collectedWorkstations, false) {
			collectedWorkstations = append(collectedWorkstations, nextStation)
		}

		station = nextStation
	}

	return station, totalPathCost, totalUnLoadCost, totalIdleTime

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

	fmt.Fprintf(fileWriter, "\nCollecting from %s to %s\n", from.Name, to.Name)
	fmt.Fprintf(fileWriter, "\nWorkstation %s will be ready at: %.2f\n", to.Name, to.GetReadyTime())
	fmt.Fprintf(fileWriter, "\nTrain idle time %.2f\n", idleTime)
	fmt.Fprintf(fileWriter, "Time to reach: %.2f\n", pathCost)
	fmt.Fprintf(fileWriter, "Workstation Unload Time: %.2f\n", unloadCost)
	fmt.Fprintf(fileWriter, "Total time spent: %.2f\n", unloadCost+pathCost+idleTime)
	fmt.Fprintf(fileWriter, "\nMAP PREVIEW\n")
	harita.PrintMap(fileWriter, path)
	fmt.Fprintf(fileWriter, "%s\n", border)

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
