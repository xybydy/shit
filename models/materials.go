// Models to be used in the Project.
package models

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

// Materials that are fetched from `yml` inputs.
var LoadedMaterials Materials

// Material containers
type Materials []Material

type Material struct {
	// Name of the material
	Name string
	// The size of a material. Quantity of material will be multiplied with `Size`
	Size int
	// Process time muliplier to be handled by workstation.
	ProcessTime float64 `yaml:"process_time"`
}

func (m Materials) Get(name string) (Material, error) {
	for i := 0; i < len(m); i++ {
		if m[i].Name == name {
			return m[i], nil
		}
	}
	return Material{}, errors.New(fmt.Sprintf("There is no such material specs: %s", name))
}

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
