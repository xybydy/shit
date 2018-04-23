package solver

import (
	"fmt"
	"strings"
	"time"

	"shit/models"
	"shit/pather"
)

type Map map[int]map[int]*Tile

func (m Map) Tile(x, y int) *Tile {
	if m[x] == nil {
		return nil
	}
	return m[x][y]
}

func (m Map) SetTile(t *Tile, x, y int) {
	if m[x] == nil {
		m[x] = map[int]*Tile{}
	}
	m[x][y] = t
	t.X = x
	t.Y = y
	t.W = m
}

func (m Map) GetTile(x, y int) *Tile {
	return m[x][y]
}

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

func ParseMap(input string) Map {
	m := Map{}
	for x, row := range strings.Split(strings.TrimSpace(input), "\n") {
		for y, raw := range row {
			kind, ok := RuneType[raw]
			if !ok {
				kind = Wall
			}

			m.SetTile(&Tile{Kind: kind}, x, y)
		}
	}
	return m
}

func (m Map) CrossCheck() bool {
	for _, row := range m {
		for _, col := range row {
			if col.Kind == Workstation {
				stations := models.LoadWorkstations()
				r := models.GetWorkstation(col.X, col.Y, stations)
				if r == nil {
					fmt.Println(col.X, col.Y, " does not match")
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
			t := m.Tile(x, y)
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

func (m Map) PrintMap(path []pather.Pather) {

	fmt.Printf("\n%s\r", m.renderMap(path))
	time.Sleep(100 * time.Millisecond)

}
