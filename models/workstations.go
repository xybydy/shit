package models

import (
	"io/ioutil"
	"fmt"
	"github.com/go-yaml/yaml"
	"strings"
	"strconv"
)

type Workstations []*Workstation

type Workstation struct {
	X            int
	Y            int
	Name         string
	Speed        int
	LoadTime     int       `yaml:"load_time"`
	UnloadTime   int       `yaml:"unload_time"`
	Requirements []string
	LoadedItems  Materials `yaml:"-"`
	ReadyItems   Materials `yaml:"-"`
}

func (w *Workstation) GetRequirements() ([]Material, []int) {
	materials := make([]Material, 0)
	amounts := make([]int, 0)

	for _, r := range w.Requirements {
		s := strings.Split(r, ",")
		material := LoadedMaterials.Get(s[0])
		materials = append(materials, material)
		amount, _ := strconv.Atoi(s[1])
		amounts = append(amounts, amount)
	}

	return materials, amounts
}

func (w *Workstation) LoadMaterial(m Material) {
	w.LoadedItems = append(w.LoadedItems, m)
}

func LoadWorkstations() Workstations {
	var workstations Workstations

	f, err := ioutil.ReadFile("inputs/worker.yml")
	if err != nil {
		fmt.Print(err)
	}

	err = yaml.Unmarshal(f, &workstations)
	if err != nil {
		fmt.Println(err)
	}

	return workstations

}

func GetWorkstation(x, y int, w Workstations) *Workstation {
	for _, station := range w {
		if station.X == x && station.Y == y {
			return station
		}
	}
	return nil
}
