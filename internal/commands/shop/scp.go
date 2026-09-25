/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package shop

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/models/balancing"
	"cafego/internal/models/shop"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"
)

func init() {
	commands.RegisterCommand(requests.C2S_SHOP_CARRIER_PIGEON,
		commands.CommandConfig{
			Name:        "CarrierPigeon",
			Identifier:  responses.S2C_SHOP_CARRIER_PIGEON,
			Description: "Buy out of stock ingredients",
			Args:        "{ingredientID} {ingredientAmount}",
			MinArgs:     5,
			MaxArgs:     5,
		},
		BuyIngredientFromShopCarrierValidator,
		BuyIngredientFromShopCarrier,
		BuyIngredientFromShopCarrierDBSaver,
	)
}

// scp - C2S_SHOP_CARRIER_PIGEON
func BuyIngredientFromShopCarrier(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {

	ingredientID, _ := strconv.Atoi(req.Args[2])

	ingredientAmount, _ := strconv.Atoi(req.Args[3])

	if shop.IsIngredientUnavailable(ingredientID) { // return true when the shop is un.
		c.Location.Cafe().AddToFridge(ingredientID, ingredientAmount)
		c.Player.UpdateAchivementCurierCount(ingredientAmount)
		c.Player.AddGold(-balancing.BalancingConstants.CourierPrice)

		c.SendExtensionResponse(cm.Identifier, "-1", "0", strconv.Itoa(ingredientID), strconv.Itoa(ingredientAmount))
	} else {
		ShopAvailibility(req, c, gm, nil) // -- cm is not used here
		return nil
	}

	return nil
}

func BuyIngredientFromShopCarrierValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	// Dont allow players to modify the packet and sending us SCP while in editor.
	if !c.Location.IsRunning() {
		return "The location is not running", commands.NOT_DECLARED
	}

	ingredientID, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	ingredientAmount, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	if shop.IsIngredientUnavailable(ingredientID) { // return true when the shop is un.
		if c.Location.Cafe().GetFridgeFreeSpace() < ingredientAmount {
			return fmt.Sprintf("Can't buy %v ingredient amount, because the fridge just have %v free space",
				ingredientAmount, c.Location.Cafe().GetFridgeFreeSpace()), commands.ERROR_FRIDGE_FULL
		}

		if c.Player.GetGold() < balancing.BalancingConstants.CourierPrice || ingredientAmount > balancing.BalancingConstants.MaxCourierSize {
			return "Player not have enough money", commands.NOT_ENOUGHT_MONEY
		}

		if ingredientAmount > balancing.BalancingConstants.MaxCourierSize {
			return "Player tried to buy more than the limit!", commands.ERROR_COURIER_SHIPPING_ERROR
		}

		if ingredientAmount < 1 {
			return "Player tried to buy less than 1!", commands.ERROR_COURIER_SHIPPING_ERROR
		}

	} else {
		return "The ingredient is now available!", commands.ERROR_COURIER_INGREDIENT_NOW_AVAILABLE
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func BuyIngredientFromShopCarrierDBSaver(c *client.Client) error {
	c.DB.UpdateFridgeInventory(c.Location.Cafe().GetID(), c.Location.Cafe().GetFridgeInventory().String())
	c.DB.UpdateGold(c.Player.GetID(), c.Player.GetGold())

	return nil
}
