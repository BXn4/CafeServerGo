package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/object"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"strconv"
)

// min level: 7

func init() {
	RegisterCommand(requests.C2S_GIFT_USE,
		CommandConfig{
			Name:       "GiftUse",
			Identifier: responses.S2C_GIFT_USE,
			MinArgs:    3,
			MaxArgs:    3,
		},
		UseGiftValidator,
		UseGift,
	)
}

func UseGift(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	slot, _ := strconv.Atoi(req.Args[2])
	gift := c.Player.Gifts[slot]
	c.Player.Gifts.RemoveGift(slot)

	item, _ := utils.GetItem(gift.ID)
	strId := strconv.Itoa(item.ID)
	strAmount := strconv.Itoa(gift.Amount)

	if item.Group == "Ingredient" {
		// Add to fridge
		c.Location.Cafe().AddToFridge(item.ID, gift.Amount)
		// Send response
		c.SendExtensionResponse("gus", "-1", "0", req.Args[2], strId, strAmount)

	} else if item.Group == "Dish" {
		// Get empty counter
		var counter *object.Object
		for _, obj := range c.Location.Cafe().GetObjects() {
			if !obj.IsCounter() {
				continue
			}
			if obj.GetDishID() == -1 || obj.GetDishID() == item.ID {
				counter = obj
				break
			}
		}

		// Set counter
		if counter.GetDishID() != item.ID {
			counter.SetDishID(item.ID)
		} else {
			counter.AddDishAmount(gift.Amount)
		}

		// Convert counter pos to string
		posX := strconv.Itoa(counter.GetPos().X)
		posY := strconv.Itoa(counter.GetPos().Y)

		c.SendExtensionResponse("gus", "-1", "0", req.Args[2], strId, strAmount, posX, posY)
	}

	c.DB.UpdateObjects(c.Location.Cafe().ID, c.Location.Cafe().Objects.StringForDB())
	c.DB.UpdateFridgeInventory(c.Location.Cafe().ID, c.Location.Cafe().GetFridgeInventory().String())
	c.DB.UpdateFurnitureInventory(c.Location.Cafe().ID, c.Location.Cafe().FurnitureInventory.String())

	SendPlayerGifts(req, c, gm)
	return nil
}

func UseGiftValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	slot, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	gift := c.Player.Gifts[slot]
	if gift == nil {
		return "Gift not found!", NOT_DECLARED
	}

	item, err := utils.GetItem(gift.ID)
	if err != nil {
		return "Item info not found!", NOT_DECLARED
	}

	if item.Group == "Ingredient" {
		if c.Location.Cafe().GetFridgeFreeSpace() < gift.Amount {
			return "Cant add gift items to the fridge!", ERROR_FRIDGE_FULL
		}

	} else if item.Group == "Dish" {
		// Get empty counter
		var counter *object.Object
		for _, obj := range c.Location.Cafe().GetObjects() {
			if !obj.IsCounter() {
				continue
			}
			if obj.GetDishID() == -1 || obj.GetDishID() == item.ID {
				counter = obj
				break
			}
		}

		// If cant find empty counter
		if counter == nil {
			return "Cant use gift, no spare counters found!", ERROR_GIFT_USE_NO_FREE_COUNTER
		}
	}

	return "Command ran without any errors.", NO_ERROR
}
