package models

import (
	"errors"
	"fmt"
)

// ´Inventories´ is the sub-object of ´Train´.
// It holds the material and the amount of material that the train carries.
type Inventories []*Inventory

// Inventory object is a single element of inventories
// Every material is mapped with ´Material´ itself and the amount of the ´Material´.
type Inventory struct {
	// Material that the ´Train´ has.
	Material Material
	// The amount of the material. The size of the material also considered.
	// If the size of the ´Material´ is 5 and the quantity is 2, 10 will be added as ´Amount´
	Amount int
}

// Add method of Inventories is used to add the materials from warehouse and back from workstations.
// The method returns ´Inventories´ object with material added from parameters.
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

// ´Pop´ method returns the ´Material´ of given name and removes it from ´Inventories´
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

// ´Get´ method returns the Material of given name but unlike ´Pop´ method,
// ´Get´ does not removes the ´Material´ from ´Inventories´
func (m Inventories) Get(name string) Material {
	for i := 0; i < len(m); i++ {
		if m[i].Material.Name == name {
			return m[i].Material
		}
	}
	return Material{}
}

// ´Details´ method returns a string of inventory details.
func (m Inventories) Details() string {
	var final string

	for _, i := range m {
		final += fmt.Sprintf("\n  %v: %d", i.Material.Name, i.Amount/i.Material.Size)
	}
	return final
}

// ´GetAmount´ method returns the available amount of ´Material´ of given name.
func (m Inventories) GetAmount(name string) int {
	for i := 0; i < len(m); i++ {
		if m[i].Material.Name == name {
			return m[i].Amount
		}
	}
	return 0
}
