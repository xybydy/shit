package utils

import (
	"fmt"
	"io/ioutil"
)

func GetMaze(filename ...string) string {
	var file string

	if len(filename) == 0 {
		file = "plant.map"
	} else if len(filename) == 1 {
		file = filename[0]
	}

	f, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}

	return fmt.Sprintf("%s", f)

}

func DivMod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}
