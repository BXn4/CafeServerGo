package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/object"
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

	gift := c.Player.Gifts[slot]
	c.Player.Gifts.RemoveGift(slot)

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
		c.Location.Cafe().AddToFridge(item.ID, gift.Amount)

		// Send response
		c.SendExtensionResponse("gus", "-1", "0", req.Args[2], strId, strAmount)
	} else if strings.ToLower(item.Group) == "dish" {

		// Get empty counter
		var counter *object.Object
		for _, obj := range c.Location.Cafe().GetObjects() {
			if !obj.IsCounter() {
				continue
			}
			if obj.GetDishID() == -1 {
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
		counter.SetDishID(item.ID)
		counter.SetDishAmount(gift.Amount)

		// Convert counter pos to string
		posX := strconv.Itoa(counter.GetPos().X)
		posY := strconv.Itoa(counter.GetPos().Y)

		c.SendExtensionResponse("gus", "-1", "0", req.Args[2], strId, strAmount, posX, posY)
	}

	SendPlayerGifts(req, c, gm)
	return nil
}
