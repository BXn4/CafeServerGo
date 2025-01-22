package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/objects"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"strconv"
	"strings"
)

func UseGift(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	slot, err := strconv.Atoi(req.Args[2])
	if err != nil {
		c.SendExtensionResponse("gus", "-1", "-1")
		return err
	}

	gift := c.Player.Gifts[0]
	c.Player.RemoveGift(slot)

	item, err := utils.GetItem(gift.ID)
	if err != nil {
		c.SendExtensionResponse("gus", "-1", "-1")
		return err
	}

	strId := strconv.Itoa(item.ID)
	strAmount := strconv.Itoa(gift.Amount)

	if strings.ToLower(item.Group) == "ingredient" {
		if c.Location.Cafe().GetFridgeCapacity() == c.Location.Cafe().GetFridgeMaxCapacity() {
			c.SendExtensionResponse("gus", "-1", "40")
			return nil
		}

		// Add to fridge
		_, ok := c.Location.Cafe().FridgeInventory[item.ID]
		if ok {
			c.Location.Cafe().FridgeInventory[item.ID] += gift.Amount
		} else {
			c.Location.Cafe().FridgeInventory[item.ID] = gift.Amount
		}

		// Send response
		c.SendExtensionResponse("gus", "-1", "0", req.Args[2], strId, strAmount)
	} else if strings.ToLower(item.Group) == "dish" {

		// Get empty counter
		var counter *objects.CafeObject
		for _, obj := range c.Location.Cafe().Objects {
			if !obj.IsCounter() {
				continue
			}
			if obj.DishID == -1 {
				counter = obj
				break
			}
		}

		// If cant find empty counter
		if counter == nil {
			c.SendExtensionResponse("gus", "-1", "20")
			return nil
		}

		// Set counter
		counter.DishID = item.ID
		counter.DishAmount = gift.Amount

		// Convert counter pos to string
		posX := strconv.Itoa(counter.Pos[0])
		posY := strconv.Itoa(counter.Pos[1])

		c.SendExtensionResponse("gus", "-1", "0", req.Args[2], strId, strAmount, posX, posY)
	}

	SendPlayerGifts(req, c, gm)
	return nil
}
