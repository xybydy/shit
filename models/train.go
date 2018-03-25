package models

import (
	"io/ioutil"
	"fmt"
	"github.com/go-yaml/yaml"
)

type Train struct {
	X        int
	Y        int
	Name     string
	Capacity int
	Speed    int
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
