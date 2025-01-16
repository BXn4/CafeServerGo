package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"fmt"
	"strconv"
)

// sbi - C2S_SHOP_BUY_ITEM
func BuyIngredient(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Dont allow players to modify the packet and sending us MJM while in editor.
	if c.Location.Cafe().InEditorMode {
		return nil
	}

	ingredientID, err := strconv.Atoi(req.Args[2])
	if err != nil {
		fmt.Printf("Can't parse id to int: %v", err)
		return err
	}

	ingredientInfo, err := utils.GetIngredient(ingredientID)
	if err != nil {
		fmt.Printf("Invalid ingredient ID: %v", err)
		return err
	}

	ingredientAmount, err := strconv.Atoi(req.Args[3])
	if err != nil || ingredientAmount <= 0 {
		fmt.Printf("Invalid ingredient amount: %v", err)
		return err
	}

	if c.Location.Cafe().GetFridgeFreeSpace() < ingredientAmount {
		c.SendExtensionResponse("sbi", "-1", "40")
		return nil
	}

	if ingredientInfo.Cash != 0 {
		if c.Player.Cash < ingredientInfo.Cash {
			c.SendExtensionResponse("sbi", "-1", "4")
			return nil
		} else {
			c.Player.Cash -= ingredientInfo.Cash
		}
	}

	if ingredientInfo.Gold != 0 {
		if c.Player.Gold < ingredientInfo.Gold {
			c.SendExtensionResponse("sbi", "-1", "4")
			return nil
		} else {
			c.Player.Gold -= ingredientInfo.Gold
		}
	}

	if c.Location.Cafe().FridgeInventory[ingredientID] != 0 {
		c.Location.Cafe().FridgeInventory[ingredientID] += ingredientAmount
	} else {
		c.Location.Cafe().FridgeInventory[ingredientID] = ingredientAmount
	}

	c.SendExtensionResponse("sbi", "-1", "0", strconv.Itoa(ingredientID), strconv.Itoa(ingredientAmount))

	// TODO: Update it in the database

	return nil
}
