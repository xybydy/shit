package models

import (
	"fmt"
	"io/ioutil"

	"shit/utils"

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

func (t *Train) unloadMaterial(materialName string) bool {
	if t.checkStock(materialName) {
		q, err := t.Stock.Pop(materialName)
		if err != nil {
			fmt.Printf("There is no item %s\n", materialName)
		}
		t.CurrentCapacity -= q.Size
		return true

	} else {
		fmt.Printf("There's no %s\n", materialName)
		return false
	}

}

func (t *Train) checkStock(material string) bool {
	if t.Stock.Get(material) == (Material{}) {
		return false
	} else {
		return true
	}
}

// Unloads the material from train to workstation.
func (t *Train) Unload(w *Workstation, start float64) float64 {
	time := 0.0
	reqs, req_am := w.GetRequirements()

	for i, material := range reqs {
		for j := 0; j < req_am[i]; j++ {
			if t.unloadMaterial(material.Name) {
				time += w.LoadMaterial(material, start)
			} else {
				fmt.Println("Burasin sonra halledicez.")
				return 0.0
			}
		}
	}
	return time
}

// Initilization function of ´Train´.
// Reads input file and creates ´Train´ object.
func LoadTrain() Train {
	var train Train

	f, err := ioutil.ReadFile("inputs/train.yml")
	if err != nil {
		fmt.Print(err)
	}

	err = yaml.Unmarshal(f, &train)
	if err != nil {
		fmt.Println(err)
	}
	return train
}

// TODO Bunu kaldıralım, LoadTrain pointer dönerse bununla aynı işi yapar bence.
func GetTrain(x, y int, t Train) *Train {
	if t.X == x && t.Y == y {
		return &t
	}
	return nil
}
