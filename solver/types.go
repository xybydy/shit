package solver

const (
	Road        = iota
	Wall
	Start
	Workstation
	FinalPath
)

var TypeRunes = map[int]rune{
	Road:        '.',
	Wall:        '#',
	Workstation: 'W',
	Start:       'S',
	FinalPath:   '●',
}

var RuneType = map[rune]int{
	'.': Road,
	'#': Wall,
	'W': Workstation,
	'S': Start,
}

var Costs = map[int]float64{
	Road:  1.0,
	Start: 1.0,
}
