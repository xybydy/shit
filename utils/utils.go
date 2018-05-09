/*
This package has miscellaneous functions used within the project
*/
package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Read `plant.map` file and return it as a string.
func GetMaze(filename ...string) string {
	var file string

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	if len(filename) == 0 {
		file = filepath.Join(dir, "inputs/plant.map")
	} else if len(filename) == 1 {
		file = filename[0]
	}

	f, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}

	return fmt.Sprintf("%s", f)

}

// Returns quotient and remainder of divison.
func DivMod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}
