package models

import (
	"io/ioutil"
	"fmt"
	"github.com/go-yaml/yaml"
)

type Materials []*Material

type Material struct {
	Name        string
	Type        string
	Size        int
	ProcessTime int `yaml:"process_time"`
}

func (m Materials) Get(name string) *Material {
	for i := 0; i < len(m); i++ {
		if m[i].Name == name {
			return m[i]
		}
	}
	return &Material{}
}

func LoadMaterials() Materials {
	var materials Materials

	f, err := ioutil.ReadFile("inputs/materials.yml")
	if err != nil {
		fmt.Print(err)
	}

	err = yaml.Unmarshal(f, &materials)
	if err != nil {
		fmt.Println(err)
	}

	return materials

}
