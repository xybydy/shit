package models

type Workstations []Workstation

type Workstation struct {
	X int
	Y int
	Name         string
	Speed        int
	LoadTime     int `yaml:"load_time"`
	UnloadTime   int `yaml:"unload_time"`
	Requirements []string
}

func (w *Workstation) parseRequirements() {
	//ret := make(Materials, len(w.Requirements))

}
