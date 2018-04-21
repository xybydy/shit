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
	Train:       'T',
	FinalPath:   '●',
}

var RuneType = map[rune]int{
	'.': Road,
	'#': Wall,
	'W': Workstation,
	'T': Train,
}

var Costs = map[int]float64{
	Road:  1.0,
	Train: 1.0,
}
