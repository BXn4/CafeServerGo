package cafe

import (
	"cafego/internal/models/simple"

	"github.com/charmbracelet/log"
)

func (c *Cafe) AddFurnitures(id, amount int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.ID < 0 {
		log.Printf("You tried to add something to the marketplaces furniture inventory")
		return
	}

	c.FurnitureInventory[id] += amount
}

func (c *Cafe) RemoveFurnitures(id, amount int) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

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

func (c *Cafe) GetFurnitureInventory() simple.IntMap {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.FurnitureInventory
}
