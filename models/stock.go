package models

import (
	"errors"
	"fmt"
)

type Inventories []*Inventory

type Inventory struct {
	Material Material
	Amount   int
}

func (m Inventories) Add(material Material, amount int) Inventories {
	item, err := m.getStock(material)
	if !err {
		item.Amount += amount
		return m
	} else {
		return append(m, &Inventory{material, amount})

	}

}

func (m Inventories) getStock(material Material) (*Inventory, bool) {
	for i := 0; i < len(m); i++ {
		if m[i].Material.Name == material.Name {
			return m[i], false
		}
	}
	return &Inventory{}, true

}

func (m Inventories) Pop(name string) (Material, error) {
	for i := 0; i < len(m); i++ {
		if m[i].Material.Name == name {
			if m[i].Amount <= 0 {
				return Material{}, errors.New(fmt.Sprintf("There is no available item: %s", name))
			} else {
				m[i].Amount -= m[i].Material.Size
				return m[i].Material, nil
			}

		}
	}
	return Material{}, errors.New(fmt.Sprintf("There is no available item: %s", name))
}

func (m Inventories) Get(name string) Material {
	for i := 0; i < len(m); i++ {
		if m[i].Material.Name == name {
			return m[i].Material
		}
	}
	return Material{}
}

func (m Inventories) Details() string {
	var final string

	for _, i := range m {
		final += fmt.Sprintf("\n  %v: %d", i.Material.Name, i.Amount/i.Material.Size)
	}
	return final
}

func (m Inventories) GetAmount(name string) int {

	for i := 0; i < len(m); i++ {
		if m[i].Material.Name == name {
			return m[i].Amount
		}
	}
	return 0

}
