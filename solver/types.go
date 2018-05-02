package solver

const (
	Road = iota
	Wall
	Train
	Workstation
	FinalPath
)

var TypeRunes = map[int]rune{
	Road:        '.',
	Wall:        '#',
	Workstation: 'W',
	Train:       'S',
	FinalPath:   '●',
}

var RuneType = map[rune]int{
	'.': Road,
	'#': Wall,
	'W': Workstation,
	'S': Train,
}

var Costs = map[int]float64{
	Road:  1.0,
	Train: 1.0,
}
