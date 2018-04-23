package models

import (
	"fmt"
	"io/ioutil"

	"shit/utils"

	"github.com/go-yaml/yaml"
)

type Train struct {
	X               int
	Y               int
	Name            string
	MaxCapacity     int `yaml:"capacity"`
	Speed           int
	Stock           Inventories
	CurrentCapacity int
}

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

func (t *Train) Unload(w *Workstation) {
	reqs, req_am := w.GetRequirements()

	for i, material := range reqs {
		for j := 0; j < req_am[i]; j++ {
			if t.unloadMaterial(material.Name) {
				w.LoadMaterial(material)
			} else {
				fmt.Println("Burasin sonra halledicez.")
			}
		}

	}
}

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

func GetTrain(x, y int, t Train) *Train {
	if t.X == x && t.Y == y {
		return &t
	}
	return nil
}
