package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/balancing"
	"cafego/internal/models/shop"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"
)

func init() {
	RegisterCommand(requests.C2S_SHOP_CARRIER_PIGEON,
		CommandConfig{
			Name:       "CarrierPigeon",
			Identifier: responses.S2C_SHOP_CARRIER_PIGEON,
			MinArgs:    5,
			MaxArgs:    5,
		},
		BuyIngredientFromShopCarrierValidator,
		BuyIngredientFromShopCarrier,
	)
}

// scp - C2S_SHOP_CARRIER_PIGEON
func BuyIngredientFromShopCarrier(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	ingredientID, _ := strconv.Atoi(req.Args[2])

	ingredientAmount, _ := strconv.Atoi(req.Args[3])

	if shop.IsIngredientUnavailable(ingredientID) { // return true when the shop is un.
		c.Location.Cafe().AddToFridge(ingredientID, ingredientAmount)
		c.Player.UpdateAchivementCurierCount(ingredientAmount)
		c.Player.AddGold(-balancing.BalancingConstants.CourierPrice)

		c.DB.UpdateFridgeInventory(c.Location.Cafe().ID, c.Location.Cafe().GetFridgeInventory().String())
		c.DB.UpdateGold(c.Player.ID, c.Player.GetGold())

		c.SendExtensionResponse("scp", "-1", "0", strconv.Itoa(ingredientID), strconv.Itoa(ingredientAmount))
	} else {
		ShopAvailibility(req, c, gm)
		return nil
	}

	return nil
}

func BuyIngredientFromShopCarrierValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	// Dont allow players to modify the packet and sending us SCP while in editor.
	if !c.Location.IsRunning() {
		return "The location is not running", NOT_DECLARED
	}

	if c.Location.Cafe().GetPlayerID() != c.Player.ID {
		return "Not the owner!", NOT_DECLARED
	}

	ingredientID, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	ingredientAmount, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	if shop.IsIngredientUnavailable(ingredientID) { // return true when the shop is un.
		if c.Location.Cafe().GetFridgeFreeSpace() < ingredientAmount {
			return fmt.Sprintf("Can't buy %v ingredient amount, because the fridge just have %v free space",
				ingredientAmount, c.Location.Cafe().GetFridgeFreeSpace()), ERROR_FRIDGE_FULL
		}

		if c.Player.GetGold() < balancing.BalancingConstants.CourierPrice || ingredientAmount > balancing.BalancingConstants.MaxCourierSize {
			return "Player not have enough money", NOT_ENOUGHT_MONEY
		}

		if ingredientAmount > balancing.BalancingConstants.MaxCourierSize {
			return "Player tried to buy more than the limit!", ERROR_COURIER_SHIPPING_ERROR
		}

		if ingredientAmount < 1 {
			return "Player tried to buy less than 1!", ERROR_COURIER_SHIPPING_ERROR
		}

	} else {
		return "The ingredient is now available!", ERROR_COURIER_INGREDIENT_NOW_AVAILABLE
	}

	return "Command ran without any errors.", NO_ERROR
}
