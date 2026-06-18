package shop

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/models/balancing"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"math"
	"strconv"
)

func init() {
	commands.RegisterCommand(requests.C2S_SHOP_DELETE_ITEM,
		commands.CommandConfig{
			Name:        "ShopSellIngredient",
			Identifier:  responses.S2C_SHOP_DELETE_ITEM,
			Description: "Sell an ingredient",
			Args:        "{ingredientID} {ingredientAmount}",
			MinArgs:     4,
			MaxArgs:     4,
		},
		SellIngredientValidator,
		SellIngredient,
		SellIngredientDBSaver,
	)
}

// sdi - C2S_SHOP_DELETE_ITEM
func SellIngredient(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	ingredientID, _ := strconv.Atoi(req.Args[2])
	sellAmount, _ := strconv.Atoi(req.Args[3])
	ingredientInfo, _ := utils.GetIngredient(ingredientID)

	// Calcualte money
	c.Player.AddCash(sellAmount * int(math.Round(float64(ingredientInfo.Cash)*(float64(balancing.BalancingConstants.SellFactorCash)/100)+float64(ingredientInfo.Gold)*(float64(balancing.BalancingConstants.SellFactorGold)))))

	c.Location.Cafe().RemoveFromFridge(ingredientID, sellAmount)

	c.SendExtensionResponse(cm.Identifier, "-1", "0", strconv.Itoa(ingredientID), strconv.Itoa(sellAmount))

	return nil
}

func SellIngredientValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	// Dont allow players to modify the packet and sending us MJM while in editor.
	if !c.Location.IsRunning() {
		return "The location is not running", commands.LOCATION_NOT_RUNNING
	}

	ingredientID, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	_, err = utils.GetIngredient(ingredientID)
	if err != nil {
		return "Invalid ingredient ID", commands.INVALID_VALUE
	}

	sellAmount, err := strconv.Atoi(req.Args[3])
	if err != nil || sellAmount <= 0 {
		return "Invalid ingredient amount", commands.INVALID_VALUE
	}

	if count, ok := c.Location.Cafe().GetFridgeInventory()[ingredientID]; !ok && count < sellAmount && sellAmount < 0 {
		return fmt.Sprintf("Invalid ingredient amount: %v, current amount: %v", sellAmount, count), commands.NOT_DECLARED
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func SellIngredientDBSaver(c *client.Client) error {
	c.DB.UpdateCash(c.Player.GetID(), c.Player.GetCash())
	c.DB.UpdateFridgeInventory(c.Location.Cafe().GetID(), c.Location.Cafe().GetFridgeInventory().String())

	return nil
}
