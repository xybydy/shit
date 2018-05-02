package models

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

// Materials that are fetched from `yml` inputs.
var LoadedMaterials Materials

// Slice of ´Material´ type.
type Materials []Material

// Type of material with exact the same type as yml file.
type Material struct {
	// Name of the material
	Name string
	// The size of a material. Quantity of material is being multiplied with `Size`
	Size int
	// Process time muliplier to be handled by workstation.
	ProcessTime float64 `yaml:"process_time"`
}

// ´Get´ method brings ´Material´ of name parameter given
// and raises error if there is no such material requested.
func (m Materials) Get(name string) (Material, error) {
	for i := 0; i < len(m); i++ {
		if m[i].Name == name {
			return m[i], nil
		}
	}
	return Material{}, errors.New(fmt.Sprintf("There is no such material specs: %s", name))
}

// TODO TO BE DELETED or make it pointer method.
func (m Materials) Pop(name string) (Material, error) {
	for i := 0; i < len(m); i++ {
		if m[i].Name == name {
			item := m[i]
			copy(m[i:], m[i+1:])
			m[len(m)-1] = Material{}
			m = m[:len(m)-1]
			return item, nil
		}
	}
	return Material{}, errors.New(fmt.Sprintf("There is no available item: %s", name))
}

// Initilization function of ´Materials´.
// Reads input file and creats ´Materials´ object with full of ´Material´ objects.
func LoadMaterials() {
	f, err := ioutil.ReadFile("inputs/materials.yml")
	if err != nil {
		fmt.Print(err)
	}

	err = yaml.Unmarshal(f, &LoadedMaterials)
	if err != nil {
		fmt.Println(err)
	}
}
