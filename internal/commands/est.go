package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"strconv"
)

// est - C2S_EDITOR_STORE_OBJECT
func StoreObject(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	objX, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}
	objY, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return err
	}
	// Dont allow players to modify the packet and sending us EST while not in editor.
	if !c.Location.Cafe().InEditorMode {
		c.SendExtensionResponse("est", "-1", "38", strconv.Itoa(objX), strconv.Itoa(objY))
		return nil
	}
	obj := c.Location.Cafe().GetObjectByPos(objX, objY)
	if obj == nil || obj.IsDoor() {
		c.SendExtensionResponse("est", "-1", "51", strconv.Itoa(objX), strconv.Itoa(objY))
		return nil
	}
	// You have to have at least one stove / counter / fridge
	// NEED A BETTER LOGIC!
	/*var stovesCount uint = 0
	var countersCount uint = 0
	var fridgesCount int = 0
	if obj.IsStove() || obj.IsCounter() || obj.IsFridge() {
		for _, object := range c.Location.Cafe().Objects {
			if object.IsStove() {
				stovesCount++
			}
			if object.IsCounter() {
				countersCount++
			}
			if object.IsFridge() {
				fridgesCount++
			}
		}
		if stovesCount == 2 || countersCount == 2 || fridgesCount == 2 || c.Location.Cafe().GetFridgeFreeSpace() < 50*fridgesCount {
			c.SendExtensionResponse("est", "-1", "39", strconv.Itoa(objX), strconv.Itoa(objY))
			return nil
		}
	} */
	if obj.DishID > 0 || obj.DishAmount > 0 {
		c.SendExtensionResponse("est", "-1", "37", strconv.Itoa(objX), strconv.Itoa(objY))
		return nil
	}
	if c.Location.Cafe().FurnitureInventory[int(obj.Kind)] != 0 {
		c.Location.Cafe().FurnitureInventory[int(obj.Kind)] += 1
	} else {
		c.Location.Cafe().FurnitureInventory[int(obj.Kind)] = 1
	}
	objectInfo, err := utils.GetItem(int(obj.Kind))
	if err != nil {
		return nil
	}
	c.Location.Cafe().RemoveObject(obj.Pos[0], obj.Pos[1])
	c.Location.Cafe().Luxury -= (objectInfo.Cash / 4000) + (objectInfo.Gold * 2)
	c.SendExtensionResponse("est", "-1", "0", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(int(obj.Kind)))
	return nil
}
