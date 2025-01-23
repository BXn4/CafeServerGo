package objects

import "log"

func (c *Cafe) AddToFridge(id, amount int) {

	if c.ID < 0 {
		log.Printf("You tried to add something to the marketplaces fridge")
		return
	}

	_, ok := c.FridgeInventory[id]
	if ok {
		c.FridgeInventory[id] += amount
	} else {
		c.FridgeInventory[id] = amount
	}
}

func (c *Cafe) GetFridgeMaxCapacity() int {

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
	if c.ID < 0 {
		log.Printf("You tried to get the free space in the markeplaces fridge")
		return 0
	}

	freeSpace := c.GetFridgeMaxCapacity()

	for ingredientID := range c.FridgeInventory {
		// Fancy does not take space in the fridge inventory. Fancys starts at ID 1401
		if ingredientID < 1401 {
			freeSpace -= 1
		}
	}

	return freeSpace
}
