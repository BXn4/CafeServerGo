package objects

import "github.com/charmbracelet/log"

func (c *Cafe) AddFurnitures(id, amount int) {

	if c.id < 0 {
		log.Printf("You tried to add something to the marketplaces furniture inventory")
		return
	}

	_, ok := c.furnitureInventory[id]
	if ok {
		c.furnitureInventory[id] += amount
	} else {
		c.furnitureInventory[id] = amount
	}
}

func (c *Cafe) RemoveFurnitures(id, amount int) bool {

	if c.id < 0 {
		log.Printf("You tried to remove something from the marketplaces furniture inventory")
		return false
	}

	_, ok := c.furnitureInventory[id]
	if ok {
		if c.furnitureInventory[id]-amount < 0 {
			return false
		}
		c.furnitureInventory[id] -= amount
	} else {
		return false
	}
	return true
}
