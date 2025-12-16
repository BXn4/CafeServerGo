package commands

import (
	"cafego/internal/client"
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
	RegisterCommand(requests.C2S_SHOP_DELETE_ITEM,
		CommandConfig{
			Name:       "ShopSellIngredient",
			Identifier: responses.S2C_SHOP_DELETE_ITEM,
			MinArgs:    4,
			MaxArgs:    4,
		},
		SellIngredientValidator,
		SellIngredient,
	)
}

// sdi - C2S_SHOP_DELETE_ITEM
func SellIngredient(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	ingredientID, _ := strconv.Atoi(req.Args[2])
	sellAmount, _ := strconv.Atoi(req.Args[3])
	ingredientInfo, _ := utils.GetIngredient(ingredientID)

	// Calcualte money
	c.Player.AddCash(sellAmount * int(math.Round(float64(ingredientInfo.Cash)*(float64(balancing.BalancingConstants.SellFactorCash)/100)+float64(ingredientInfo.Gold)*(float64(balancing.BalancingConstants.SellFactorGold)))))

	c.Location.Cafe().RemoveFromFridge(ingredientID, sellAmount)

	c.SendExtensionResponse("sdi", "-1", "0", strconv.Itoa(ingredientID), strconv.Itoa(sellAmount))

	c.DB.UpdateCash(c.Player.ID, c.Player.GetCash())
	c.DB.UpdateFridgeInventory(c.Location.Cafe().ID, c.Location.Cafe().GetFridgeInventory().String())

	return nil
}

func SellIngredientValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	// Dont allow players to modify the packet and sending us MJM while in editor.
	if !c.Location.IsRunning() {
		return "The location is not running", LOCATION_NOT_RUNNING
	}

	if c.Location.Cafe().GetPlayerID() != c.Player.ID {
		return "Not the owner!", NOT_DECLARED
	}

	ingredientID, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	_, err = utils.GetIngredient(ingredientID)
	if err != nil {
		return "Invalid ingredient ID", INVALID_VALUE
	}

	sellAmount, err := strconv.Atoi(req.Args[3])
	if err != nil || sellAmount <= 0 {
		return "Invalid ingredient amount", INVALID_VALUE
	}

	if count, ok := c.Location.Cafe().GetFridgeInventory()[ingredientID]; !ok && count < sellAmount && sellAmount < 0 {
		return fmt.Sprintf("Invalid ingredient amount: %v, current amount: %v", sellAmount, count), NOT_DECLARED
	}

	return "Command ran without any errors.", NO_ERROR
}
