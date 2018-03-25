package models

import (
	"io/ioutil"
	"fmt"
	"github.com/go-yaml/yaml"
)

type Workstations []*Workstation

type Workstation struct {
	X            int
	Y            int
	Name         string
	Speed        int
	LoadTime     int `yaml:"load_time"`
	UnloadTime   int `yaml:"unload_time"`
	Requirements []string
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
