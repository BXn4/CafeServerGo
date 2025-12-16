package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"strconv"
)

// min level = 7

func init() {
	RegisterCommand(requests.C2S_SHOP_BUY_ITEM,
		CommandConfig{
			Name:       "ShopBuyItem",
			Identifier: responses.S2C_SHOP_BUY_ITEM,
			MinArgs:    4,
			MaxArgs:    4,
		},
		BuyIngredientValidator,
		BuyIngredient,
	)
}

// sbi - C2S_SHOP_BUY_ITEM
func BuyIngredient(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	// SCP command hanles when buying from currior

	ingredientID, _ := strconv.Atoi(req.Args[2])

	ingredientInfo, _ := utils.GetIngredient(ingredientID)

	ingredientAmount, _ := strconv.Atoi(req.Args[3])

	c.Player.AddCash(-ingredientInfo.Cash)

	c.Player.AddGold(-ingredientInfo.Gold)

	c.Location.Cafe().AddToFridge(ingredientID, ingredientAmount)

	c.SendExtensionResponse("sbi", "-1", "0", strconv.Itoa(ingredientID), strconv.Itoa(ingredientAmount))

	c.Player.UpdateAchivementBoughtIngredients()
	if c.Player.IsTutorialCompleted {
		c.DB.UpdateAchievement(c.Player.ID, c.Player.GetAchivements().String())
		c.DB.UpdateCash(c.Player.ID, c.Player.GetCash())
		c.DB.UpdateGold(c.Player.ID, c.Player.GetGold())
		c.DB.UpdateFridgeInventory(c.Location.Cafe().ID, c.Location.Cafe().GetFridgeInventory().String())
	}

	return nil
}

func BuyIngredientValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
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

	ingredientInfo, err := utils.GetIngredient(ingredientID)
	if err != nil {
		return "Invalid ingredient ID", INVALID_VALUE
	}

	if ingredientInfo.Category == "fancy" && c.Player.GetLevel() < 7 {
		return "Cant buy that ingredient, because the player not yet reached the feature.", NOT_DECLARED
	}

	if c.Player.GetLevel() < ingredientInfo.Level {
		return "Cant buy that ingredient, because the player not yet reached the level.", NOT_DECLARED
	}

	ingredientAmount, err := strconv.Atoi(req.Args[3])
	if err != nil || ingredientAmount <= 0 {
		return "Invalid ingredient amount", INVALID_VALUE
	}

	if !c.Player.IsTutorialCompleted {
		if ingredientID != 1318 && ingredientAmount != 1 {
			return "Invalid tutorial values, needed: 1318 for ingredient, and 1 for amount!", INVALID_VALUE
		}
	}

	if c.Location.Cafe().GetFridgeFreeSpace() < ingredientAmount {
		return fmt.Sprintf("Can't buy %v ingredient amount, because the fridge just have %v free space",
			ingredientAmount, c.Location.Cafe().GetFridgeFreeSpace()), ERROR_FRIDGE_FULL
	}

	if ingredientInfo.Cash != 0 {
		if c.Player.GetCash() < ingredientInfo.Cash {
			return "Player not have enough cash to buy it", NOT_ENOUGHT_MONEY
		}
	}

	if ingredientInfo.Gold != 0 {
		if c.Player.GetGold() < ingredientInfo.Gold {
			return "Player not have enough gold to buy it", NOT_ENOUGHT_MONEY
		}
	}

	return "Command ran without any errors.", NO_ERROR
}
