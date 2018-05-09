/*
Objects and methods for pathfinding algorithm
*/

package solver

const (
	Road = iota
	Wall
	Train
	Workstation
	FinalPath
)

// Tile representation as a char in the map.
//
// Available road represented with `.`
// Each wall represented with `#`
// `W` character represents Workstation objects
// Train / Storage represented with `S`
// Selected path in optimum path represented with `●`
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
