package solver

import (
	"fmt"
	"io"
	"strings"

	"semesterproject/models"
	"semesterproject/pather"
)

// `Map` object is basically a 2-dimensional matrix of `Tile` objects
type Map map[int]map[int]*Tile

// Inserts `Tile` object into map of given coordinations
func (m Map) SetTile(t *Tile, x, y int) {
	if m[x] == nil {
		m[x] = map[int]*Tile{}
	}
	m[x][y] = t
	t.X = x
	t.Y = y
	t.W = m
}

// Returns `Tile` object of given coordinations.
func (m Map) GetTile(x, y int) *Tile {
	if m[x] == nil {
		return nil
	}
	return m[x][y]
}

// Returns the first occurrence of the given `Tile` type.
func (m Map) FirstOfKind(kind int) *Tile {
	for _, row := range m {
		for _, t := range row {
			if t.Kind == kind {
				return t
			}
		}
	}
	return nil
}

// Acts like `FirstOfKind` method however this method returns all objects of given type of `Tile`
func (m Map) GetKind(kind int) Tiles {
	var kinds Tiles
	for _, row := range m {
		for _, col := range row {
			if col.Kind == kind {
				kinds = append(kinds, col)
			}
		}
	}
	return kinds
}

// Crosschecks all tiles and model coordinations and if one of the type does not match return `false` and
// simulation stops running.
func (m Map) CrossCheck() bool {
	for _, row := range m {
		for _, col := range row {
			if col.Kind == Workstation {
				stations := models.LoadWorkstations()
				r := models.GetWorkstation(col.X, col.Y, stations)
				if r == nil {
					fmt.Println("Workstation at", col.X, col.Y, " does not match with any model")
					return false
				}
			} else if col.Kind == Train {
				train := models.LoadTrain()
				r := models.GetTrain(col.X, col.Y, train)
				if r == nil {
					fmt.Println("Train does not match")
					return false
				}
			}

		}
	}
	return true
}

// Renders the map and the path.
func (m Map) renderMap(path []pather.Pather) string {
	width := len(m)
	if width == 0 {
		return ""
	}
	height := len(m[0])

	pathLocs := map[string]bool{}
	for _, p := range path {
		pT := p.(*Tile)
		pathLocs[fmt.Sprintf("%d,%d", pT.X, pT.Y)] = true
	}

	rows := make([]string, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			t := m.GetTile(x, y)
			r := ' '
			if pathLocs[fmt.Sprintf("%d,%d", x, y)] {
				r = TypeRunes[FinalPath]
			} else if t != nil {
				r = TypeRunes[t.Kind]
			}
			rows[y] += string(r)
		}
	}
	return strings.Join(rows, "\n")
}

// Renders the map.
func (m Map) PrintMap(i io.Writer, path []pather.Pather) {

	// fmt.Printf("\n%s\r", m.renderMap(path))
	fmt.Fprintf(i, "\n%s\r", m.renderMap(path))

}

// Returns the x,y dimension of the `Map`.
func (m Map) GetSize() (row, col int) {
	return len(m), len(m[0])
}

// Goes through the given `string` input and builds Map with full of `Tiles`
func ParseMap(input string) Map {
	m := Map{}
	for x, row := range strings.Split(strings.TrimSpace(input), "\n") {
		for y, raw := range row {
			if raw == 13 {continue} //boslugu okursa bu turu atla
			
			kind, ok := RuneType[raw]
			if !ok {
				kind = Wall
			}

			m.SetTile(&Tile{Kind: kind}, x, y)
		}
	}
	return m
}
