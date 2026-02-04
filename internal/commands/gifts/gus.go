package gifts

import (
	"cafego/internal/client"
	"cafego/internal/commands"
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
	commands.RegisterCommand(requests.C2S_GIFT_USE,
		commands.CommandConfig{
			Name:         "GiftUse",
			Identifier:   responses.S2C_GIFT_USE,
			Description:  "Using an gift from the gift inventory",
			Args:         "{slotID} {giftID} {giftAmount} {posX} {posY}",
			MinArgs:      3,
			MaxArgs:      7,
			FeatureLevel: 7,
		},
		UseGiftValidator,
		UseGift,
		UseGiftDBSaver,
	)
}

func UseGift(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {

	slot, _ := strconv.Atoi(req.Args[2])
	gifts := c.Player.GetGifts()
	gift := gifts[slot]
	c.Player.RemoveGift(slot)

	item, _ := utils.GetItem(gift.ID)
	strId := strconv.Itoa(item.ID)
	strAmount := strconv.Itoa(gift.Amount)

	if item.Group == "Ingredient" {
		// Add to fridge
		c.Location.Cafe().AddToFridge(item.ID, gift.Amount)
		// Send response
		c.SendExtensionResponse(cm.Identifier, "-1", "0", req.Args[2], strId, strAmount)

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

		c.SendExtensionResponse(cm.Identifier, "-1", "0", req.Args[2], strId, strAmount, posX, posY)
	}

	SendPlayerGifts(req, c, gm, nil)
	return nil
}

func UseGiftValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	if c.Player.GetLevel() < cm.FeatureLevel {
		return "Player not yet unlocked the feature!", commands.NOT_DECLARED
	}

	slot, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	gifts := c.Player.GetGifts()
	gift := gifts[slot]
	if gift == nil {
		return "Gift not found!", commands.NOT_DECLARED
	}

	item, err := utils.GetItem(gift.ID)
	if err != nil {
		return "Item info not found!", commands.NOT_DECLARED
	}

	if item.Group == "Ingredient" {
		if c.Location.Cafe().GetFridgeFreeSpace() < gift.Amount {
			return "Cant add gift items to the fridge!", commands.ERROR_FRIDGE_FULL
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
			return "Cant use gift, no spare counters found!", commands.ERROR_GIFT_USE_NO_FREE_COUNTER
		}
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func UseGiftDBSaver(c *client.Client) error {
	c.DB.UpdateObjects(c.Location.Cafe().GetID(), c.Location.Cafe().GetObjects().StringForDB())
	c.DB.UpdateFridgeInventory(c.Location.Cafe().GetID(), c.Location.Cafe().GetFridgeInventory().String())
	c.DB.UpdateFurnitureInventory(c.Location.Cafe().GetID(), c.Location.Cafe().GetFurnitureInventory().String())
	c.DB.UpdateGifts(c.Player.GetID(), c.Player.GetGifts().String())

	return nil
}
