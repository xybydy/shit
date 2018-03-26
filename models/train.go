package models

import (
	"io/ioutil"
	"fmt"
	"github.com/go-yaml/yaml"
)

type Train struct {
	X               int
	Y               int
	Name            string
	MaxCapacity     int `yaml:"capacity"`
	Speed           int
	ingredients     Materials
	CurrentCapacity int
}

//Ilk stok doldurmayi neye gore yapacaklar.
//Simdilik tren bir kerelik tum workstationlara ait olan mallari topluyor
func (t *Train) LoadFromStorage(w Workstations) {
	for _, station := range w {
		reqs, reqAmount := station.GetRequirements()
		for i, amount := range reqAmount {
			for j := 0; j < amount; j++ {
				t.ingredients = append(t.ingredients, reqs[i])
				t.CurrentCapacity += reqs[i].Size
			}
		}
	}
}

func (t *Train) UnloadMaterial(materialName string) bool {
	if t.checkStock(materialName) {
		q, err := t.ingredients.Pop(materialName)
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
	if t.ingredients.Get(material) == (Material{}) {
		return false
	} else {
		return true
	}
}

func (t *Train) RemainingMaterials() []Material {

	return t.ingredients

}

func (t *Train) Unload(w *Workstation) {
	reqs, req_am := w.GetRequirements()

	for i, material := range reqs {
		for j := 0; j < req_am[i]; j++ {
			if t.UnloadMaterial(material.Name) {
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
