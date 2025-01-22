package objects

func (c *Cafe) AddToFridge(id, amount int) {

	println("added to fridge: ", id, amount)
	_, ok := c.FridgeInventory[id]
	if ok {
		c.FridgeInventory[id] += amount
	} else {
		c.FridgeInventory[id] = amount
	}
}

func (c *Cafe) GetFridgeMaxCapacity() int {
	fridgeCount := 0
	for _, obj := range c.Objects {
		if obj.IsFridge() {
			fridgeCount++
		}
	}
	return fridgeCount * 50
}

func (c *Cafe) GetFridgeCapacity() int {
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
	freeSpace := c.GetFridgeMaxCapacity()

	for ingredientID := range c.FridgeInventory {
		// Fancy does not take space in the fridge inventory. Fancys starts at ID 1401
		if ingredientID < 1401 {
			freeSpace -= 1
		}
	}

	return freeSpace
}
