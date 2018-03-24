package models

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
