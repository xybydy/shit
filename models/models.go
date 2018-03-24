package models

import "errors"

type Nodes []*Node

type Node struct {
	X          int
	Y          int
	Type       string
	Neighbours []*Node
}

func (n Nodes) Get(x, y int) (*Node, error) {
	for i := 0; i < len(n); i++ {
		if n[i].X == x && n[i].Y == y {
			return n[i], nil
		}
	}
	return nil, errors.New("Out of bounds")
}
