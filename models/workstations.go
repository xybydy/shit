package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-yaml/yaml"
)

// Workstations that are fetched from `yml` inputs.
type Workstations []*Workstation

type Workstation struct {
	// X coordination of ´Workstation´
	X int
	// Y coordination of ´Workstation´
	Y int
	// Unique name of ´Workstation´
	Name string
	// Processing speed of workstation. Higher speed is faster processing time.
	Speed float64
	// Loading time of all demanded materials. Higher LoadTime causes
	// longer time to load materials.
	LoadTime float64 `yaml:"load_time"`
	// Unloading time of all demanded materials. Higher UnloadTime causes
	// longer time to unload materials.
	UnloadTime float64 `yaml:"unload_time"`
	// Slice of requirements in raw string type.
	Requirements []string
	// LoadedItems is the list of materials which are loaded into Workstation from the Train.
	LoadedItems Materials `yaml:"-"`
	// The time of when the Train has loaded the materials into workstation
	filledTime float64
}

// Returns the list of of materials and required amounts of those as a slice.
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

// Returns the string of required items with amounts as a string
func (w *Workstation) PrintRequirements() string {
	var finalString string
	reqs, amts := w.GetRequirements()

	for i := 0; i < len(reqs); i++ {
		finalString += fmt.Sprintf("  %s: %d\n", reqs[i].Name, amts[i])

	}
	return finalString
}

// Loads given material into workstation and updates filledTime accordingly.
func (w *Workstation) LoadMaterial(m Material, in float64) float64 {
	w.LoadedItems = append(w.LoadedItems, m)
	w.filledTime = in
	return w.LoadTime
}

// Returns the time of when Workstation can unload the processed materials.
func (w *Workstation) GetReadyTime() float64 {
	var totalProcessTime float64
	for _, l := range LoadedMaterials {
		totalProcessTime += l.ProcessTime

	}
	return w.filledTime + totalProcessTime*(100/w.Speed)
}

// Initilization function of ´Workstations´.
// Reads input file and creats ´Workstations´ object with full of ´Workstation´ objects.
func LoadWorkstations(input ...string) Workstations {
	var in string

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	if len(input) == 0 {
		in = filepath.Join(dir, "inputs/worker.yml")
	} else if len(input) == 1 {
		in = input[0]
	}

	var workstations Workstations

	f, err := ioutil.ReadFile(in)
	if err != nil {
		fmt.Print(err)
	}

	err = yaml.Unmarshal(f, &workstations)
	if err != nil {
		fmt.Println(err)
	}

	return workstations
}

// Returns the workstation object of given coordinations.
func GetWorkstation(x, y int, w Workstations) *Workstation {
	for _, station := range w {
		if station.X == x && station.Y == y {
			return station
		}
	}
	return nil
}
