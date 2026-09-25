/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package cafe

import (
	"cafego/internal/models/simple"

	"github.com/charmbracelet/log"
)

func (c *Cafe) AddToFridge(id, amount int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.ID < 0 {
		log.Printf("You tried to add something to the marketplaces fridge")
		return
	}

	c.FridgeInventory[id] += amount
}

func (c *Cafe) GetFridgeMaxCapacity() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.ID < 0 {
		log.Printf("You tried to get the max capacity of the markeplaces fridge ")
		return 0
	}

	fridgeCount := 0
	for _, obj := range c.Objects {
		if obj.IsFridge() {
			fridgeCount++
		}
	}
	return fridgeCount * 50
}

func (c *Cafe) GetFridgeCapacity() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.ID < 0 {
		log.Printf("You tried to get the capacity of the markeplaces fridge")
		return 0
	}

	capacity := 0
	for ingredientID := range c.FridgeInventory {
		// Fancy does not take space in the fridge inventory. Fancys starts at ID 1401
		if ingredientID >= 1401 {
			capacity++
		}
	}
	return capacity
}

func (c *Cafe) GetFridgeFreeSpace() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.ID < 0 {
		log.Printf("You tried to get the free space in the markeplaces fridge")
		return 0
	}

	freeSpace := c.fridgeCapacity

	for ingredientID := range c.FridgeInventory {
		// Fancy does not take space in the fridge inventory. Fancys starts at ID 1401
		if ingredientID < 1401 {
			freeSpace -= 1
		}
	}

	return freeSpace
}

func (c *Cafe) RemoveFromFridge(id, amount int) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.ID < 0 {
		log.Printf("You tried to remove something from the marketplaces fridge")
		return false
	}

	_, ok := c.FridgeInventory[id]
	if ok {
		if c.FridgeInventory[id]-amount < 0 {
			return false
		}
		c.FridgeInventory[id] -= amount
	}
	if c.FridgeInventory[id] == 0 {
		delete(c.FridgeInventory, id)
	}
	return false
}

func (c *Cafe) GetFridgeInventory() simple.IntMap {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.FridgeInventory
}
