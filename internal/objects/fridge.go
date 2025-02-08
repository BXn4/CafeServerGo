package objects

import "github.com/charmbracelet/log"

func (c *Cafe) AddToFridge(id, amount int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.id < 0 {
		log.Printf("You tried to add something to the marketplaces fridge")
		return
	}

	_, ok := c.fridgeInventory[id]
	if ok {
		c.fridgeInventory[id] += amount
	} else {
		c.fridgeInventory[id] = amount
	}
}

func (c *Cafe) GetFridgeMaxCapacity() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.id < 0 {
		log.Printf("You tried to get the max capacity of the markeplaces fridge ")
		return 0
	}

	fridgeCount := 0
	for _, obj := range c.objects {
		if obj.IsFridge() {
			fridgeCount++
		}
	}
	return fridgeCount * 50
}

func (c *Cafe) GetFridgeCapacity() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.id < 0 {
		log.Printf("You tried to get the capacity of the markeplaces fridge")
		return 0
	}

	capacity := 0
	for ingredientID := range c.fridgeInventory {
		// Fancy does not take space in the fridge inventory. Fancys starts at ID 1401
		if ingredientID >= 1401 {
			capacity++
		}
	}
	return capacity
}

func (c *Cafe) GetFridgeFreeSpace() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.id < 0 {
		log.Printf("You tried to get the free space in the markeplaces fridge")
		return 0
	}

	freeSpace := c.GetFridgeMaxCapacity()

	for ingredientID := range c.fridgeInventory {
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
	if c.id < 0 {
		log.Printf("You tried to remove something from the marketplaces fridge")
		return false
	}

	_, ok := c.fridgeInventory[id]
	if ok {
		if c.fridgeInventory[id]-amount < 0 {
			return false
		}
		c.fridgeInventory[id] -= amount
	}
	if c.fridgeInventory[id] == 0 {
		delete(c.fridgeInventory, id)
	}
	return false
}
