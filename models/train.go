/*
Models package contains all objects and methods of objects used in the simulation.
The objects are `Material`, `Stock`, `Workstation`, `Train`.
*/
package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"semesterproject/utils"

	"github.com/go-yaml/yaml"
)

// Type of train with exact the same type as yml file.
// Although it is named as "Train", this object also acts as warehouse at the beginning of the simulation.
// So the location of the train will also be startpoint and the endpoint of simulation.
type Train struct {
	// X coordination of ´Train´
	X int
	// X coordination of ´Train´
	Y int
	// Unique name of ´Train´
	Name string
	// The maximum capacity that Train can contain at once.
	// Please note that the capacity contains materials with size and quantities multiplied.
	// But for now to prevent unexpected behaviour it is recommended for users to give the Train higher capacity than Workstations demands
	MaxCapacity int `yaml:"capacity"`
	// The transportation speed of ´Train´
	Speed float64
	// All materials are brought from warehouse added in ´Stock´
	Stock Inventories
	// CurrentCapacity used for to check if any available space is left in Train Stock.
	CurrentCapacity int
}

// This method goes through each workstations, collects all demands of workstations
// and add the materials to Train Stock.
func (t *Train) LoadFromStorage(w Workstations) {
	for _, station := range w {
		reqs, reqAmount := station.GetRequirements()
		for i, req := range reqs {

			if t.MaxCapacity-t.CurrentCapacity > reqAmount[i]*req.Size {
				t.Stock = t.Stock.Add(req, reqAmount[i]*req.Size)
				t.CurrentCapacity += req.Size * reqAmount[i]
			} else {
				div, _ := utils.DivMod(t.MaxCapacity-t.CurrentCapacity,
					req.Size)
				t.CurrentCapacity += req.Size * div
				t.Stock = t.Stock.Add(req, div*req.Size)
			}
		}
	}
}

// Checks whether the given material exists. If the material exists
// the material unload from `Train` and opens space.
func (t *Train) unloadMaterial(materialName string) bool {
	if t.checkStock(materialName) {
		q, err := t.Stock.Pop(materialName)
		if err != nil {
			fmt.Printf("There is no item %s\n", materialName)
		}
		t.CurrentCapacity -= q.Size
		return true

	}
	fmt.Printf("There's no %s\n", materialName)
	return false

}

// Returns `true` if `Train` has the given material name, returns `false` otherwise.
func (t *Train) checkStock(material string) bool {
	if t.Stock.Get(material) == (Material{}) {
		return false
	}
	return true

}

// Unloads the material from train to workstation.
func (t *Train) Unload(w *Workstation, start float64) float64 {
	time := 0.0
	reqs, requestAmount := w.GetRequirements()

	for i, material := range reqs {
		for j := 0; j < requestAmount[i]; j++ {
			if t.unloadMaterial(material.Name) {
				time += w.LoadMaterial(material, start)
			}
			fmt.Println("Burasin sonra halledicez.")
			return 0.0

		}
	}
	return time
}

// Initilization function of ´Train´.
// Reads input file and creates ´Train´ object.
func LoadTrain(input ...string) Train {
	var in string
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	if len(input) == 0 {
		in = filepath.Join(dir, "inputs/train.yml")
	} else if len(input) == 1 {
		in = input[0]
	}

	var train Train

	f, err := ioutil.ReadFile(in)
	if err != nil {
		fmt.Print(err)
	}

	err = yaml.Unmarshal(f, &train)
	if err != nil {
		fmt.Println(err)
	}
	return train
}

// Returns `Train` object if given coordinates matches with `Train` object given. Return `nil` if coordinates does not match with the object.
// This function used to cross check `Tile` object and `Train` object.
func GetTrain(x, y int, t Train) *Train {
	if t.X == x && t.Y == y {
		return &t
	}
	return nil
}
