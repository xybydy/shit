package models

import (
	"os"
	"log"
	"io/ioutil"
	"strings"
)

type Map struct {
	Start        *Node
	End          *Node
	rightBorder  int
	bottomBorder int
	Nodes        Nodes
}

func (m *Map) parseMap(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}

	e, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}

	str := strings.Split(string(e), "\n")
	if m.bottomBorder == 0 {
		m.bottomBorder = len(str) - 1

		if m.rightBorder == 0 {
			m.rightBorder = len(str[0]) - 1
		}

		for i, x := range str {

			for j, y := range x {
				var nodeType string

				switch y {
				case '#':
					nodeType = "wall"
				case '.':
					nodeType = "path"
				case 'S':
					nodeType = "start"
				case 'W':
					nodeType = "workstation"
				}

				m.Nodes = append(m.Nodes, &Node{X: i, Y: j, Type: nodeType})
			}
		}
	}
}

func (m *Map) buildNeighbours() {
	for i := 0; i < len(m.Nodes); i++ {
		node := m.Nodes[i]

		right := node.Y + 1
		left := node.Y - 1
		down := node.X + 1
		up := node.X - 1

		if node.Type == "start" {
			m.Start = node
		}

		if node.Type == "finish" {
			m.End = node
		}

		upNode, err := m.Nodes.Get(up, node.Y)
		if err != nil {
			log.Println("Out of bounds, to up!")
		}

		downNode, err := m.Nodes.Get(down, node.Y)
		if err != nil {
			log.Println("Out of bounds, to deep down!")
		}

		rightNode, err := m.Nodes.Get(right, node.Y)
		if err != nil {
			log.Println("Out of bounds, to the right")
		}

		leftNode, err := m.Nodes.Get(left, node.Y)
		if err != nil {
			log.Println("Out of bounds, to the left")
		}

		if 0 <= up && up <= m.bottomBorder && upNode != nil && upNode.Type != "wall" { //UP
			node.Neighbours = append(node.Neighbours, upNode)
		}

		if 0 <= right && right <= m.rightBorder && rightNode != nil && rightNode.Type != "wall" { //RIGHT
			node.Neighbours = append(node.Neighbours, rightNode)
		}

		if 0 <= down && down <= m.bottomBorder && downNode != nil && downNode.Type != "wall" { //DOWN
			node.Neighbours = append(node.Neighbours, downNode)
		}

		if 0 <= left && left <= m.rightBorder && leftNode != nil && leftNode.Type != "wall" { //LEFT
			node.Neighbours = append(node.Neighbours, leftNode)
		}

	}
}

func (m *Map) DoIt(path string) {
	m.parseMap(path)
	m.buildNeighbours()
}



