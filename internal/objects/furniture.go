package objects

import "log"

func (c *Cafe) AddFurnitures(id, amount int) {

	if c.ID < 0 {
		log.Printf("You tried to add something to the marketplaces furniture inventory")
		return
	}

	_, ok := c.FurnitureInventory[id]
	if ok {
		c.FurnitureInventory[id] += amount
	} else {
		c.FurnitureInventory[id] = amount
	}
}

func (c *Cafe) RemoveFurnitures(id, amount int) bool {

	if c.ID < 0 {
		log.Printf("You tried to remove something from the marketplaces furniture inventory")
		return false
	}

	_, ok := c.FurnitureInventory[id]
	if ok {
		if c.FurnitureInventory[id]-amount < 0 {
			return false
		}
		c.FurnitureInventory[id] -= amount
	} else {
		return false
	}
	return true
}
