package models

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/go-yaml/yaml"
)

type Workstations []*Workstation

type Workstation struct {
	X            int
	Y            int
	Name         string
	Status       bool `yaml:"-"`
	Speed        float64
	LoadTime     float64 `yaml:"load_time"`
	UnloadTime   float64 `yaml:"unload_time"`
	Requirements []string
	LoadedItems  Materials `yaml:"-"`
	filledTime   float64
}

func (w *Workstation) GetRequirements() ([]Material, []int) {
	materials := make([]Material, 0)
	amounts := make([]int, 0)

	for _, r := range w.Requirements {
		s := strings.Split(r, ",")
		material, err := LoadedMaterials.Get(s[0])
		if err != nil {
			fmt.Println(err)
		}
		materials = append(materials, material)
		amount, _ := strconv.Atoi(s[1])
		amounts = append(amounts, amount)
	}

	return materials, amounts
}

func (w *Workstation) PrintRequirements() string {
	var finalString string
	reqs, amts := w.GetRequirements()

	for i := 0; i < len(reqs); i++ {
		finalString += fmt.Sprintf("  %s: %d\n", reqs[i].Name, amts[i])

	}
	return finalString
}

func (w *Workstation) LoadMaterial(m Material, in float64) float64 {
	w.LoadedItems = append(w.LoadedItems, m)
	w.filledTime = in
	return w.LoadTime
}

func (w *Workstation) GetReadyTime() float64 {
	var totalProcessTime float64
	for _, l := range LoadedMaterials {
		totalProcessTime += l.ProcessTime

	}
	return w.filledTime + totalProcessTime*(100/w.Speed)
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
