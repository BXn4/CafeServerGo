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
	if !c.Location.IsRunning() {
		return nil
	}

	ingredientID, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return fmt.Errorf("Can't parse id to int: %v", err)
	}

	ingredientInfo, err := utils.GetIngredient(ingredientID)
	if err != nil {
		return fmt.Errorf("Invalid ingredient ID: %v", err)
	}

	ingredientAmount, err := strconv.Atoi(req.Args[3])
	if err != nil || ingredientAmount <= 0 {
		return fmt.Errorf("Invalid ingredient amount: %v", err)
	}

	if c.Location.Cafe().GetFridgeFreeSpace() < ingredientAmount {
		c.SendExtensionResponse("sbi", "-1", "40")
		return nil
	}

	if ingredientInfo.Cash != 0 {
		if c.Player.GetCash() < ingredientInfo.Cash {
			c.SendExtensionResponse("sbi", "-1", "4")
			return nil
		} else {
			c.Player.AddCash(-ingredientInfo.Cash)

			c.Player.UpdateAchivementBoughtIngredients()

			c.DB.UpdateAchievement(c.Player.ID, c.Player.GetAchivements().String())
		}
	}

	if ingredientInfo.Gold != 0 {
		if c.Player.GetGold() < ingredientInfo.Gold {
			c.SendExtensionResponse("sbi", "-1", "4")
			return nil
		} else {
			c.Player.AddGold(ingredientInfo.Gold)
		}
	}

	c.Location.Cafe().AddToFridge(ingredientID, ingredientAmount)

	c.SendExtensionResponse("sbi", "-1", "0", strconv.Itoa(ingredientID), strconv.Itoa(ingredientAmount))

	c.DB.UpdateCash(c.Player.ID, c.Player.GetCash())
	c.DB.UpdateGold(c.Player.ID, c.Player.GetGold())
	c.DB.UpdateFridgeInventory(c.Location.Cafe().ID, c.Location.Cafe().GetFridgeInventory().String())

	return nil
}
