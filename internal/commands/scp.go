package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"fmt"
	"strconv"
)

// scp - C2S_SHOP_CARRIER_PIGEON
func BuyIngredientFromShopCarrier(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Dont allow players to modify the packet and sending us SCP while in editor.
	if !c.Location.IsRunning() {
		return nil
	}

	ingredientID, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return fmt.Errorf("Can't parse id to int: %v", err)
	}

	ingredientAmount, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return fmt.Errorf("Can't parse amount to int: %v", err)
	}

	if gm.IsIngredientUnavailable(ingredientID) { // return true when the shop is un.
		if c.Location.Cafe().GetFridgeFreeSpace() < ingredientAmount {
			c.SendExtensionResponse("scp", "-1", "40")
			return nil
		}

		if c.Player.GetGold() <= 0 {
			c.SendExtensionResponse("scp", "-1", "4")
			return nil
		}
	} else {
		ShopAvailibility(req, c, gm)
		c.SendExtensionResponse("scp", "-1", "42")
		return nil
	}

	c.Location.Cafe().AddToFridge(ingredientID, ingredientAmount)
	c.Player.UpdateAchivementCurierCount(ingredientAmount)
	c.Player.AddGold(-1)

	c.DB.UpdateFridgeInventory(c.Location.Cafe().ID, c.Location.Cafe().GetFridgeInventory().String())
	c.DB.UpdateGold(c.Player.ID, c.Player.GetGold())

	c.SendExtensionResponse("scp", "-1", "0", strconv.Itoa(ingredientID), strconv.Itoa(ingredientAmount))

	return nil
}
