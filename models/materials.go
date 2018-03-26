package models

import (
	"io/ioutil"
	"fmt"
	"github.com/go-yaml/yaml"
	"errors"
)

var LoadedMaterials Materials

type Materials []Material

type Material struct {
	Name        string
	Type        string
	Size        int
	ProcessTime int `yaml:"process_time"`
}

func (m Materials) Get(name string) Material {
	for i := 0; i < len(m); i++ {
		if m[i].Name == name {
			return m[i]
		}
	}
	return Material{}
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
