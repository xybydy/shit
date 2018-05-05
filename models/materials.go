package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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
	return Material{}, fmt.Errorf("There is no such material specs: %s", name)
}

// Initilization function of ´Materials´.
// Reads input file and creats ´Materials´ object with full of ´Material´ objects.
func LoadMaterials(input ...string) {
	var in string

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	if len(input) == 0 {
		in = filepath.Join(dir, "inputs/materials.yml")
	} else if len(input) == 1 {
		in = input[0]
	}

	f, err := ioutil.ReadFile(in)
	if err != nil {
		fmt.Print(err)
	}

	err = yaml.Unmarshal(f, &LoadedMaterials)
	if err != nil {
		fmt.Println(err)
	}
}
