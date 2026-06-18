package shop

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"strconv"
)

// min level = 7

func init() {
	commands.RegisterCommand(requests.C2S_SHOP_BUY_ITEM,
		commands.CommandConfig{
			Name:        "ShopBuyItem",
			Identifier:  responses.S2C_SHOP_BUY_ITEM,
			Description: "Buy an ingredient from the shop",
			Args:        "{id} {amount}",
			MinArgs:     4,
			MaxArgs:     4,
		},
		BuyIngredientValidator,
		BuyIngredient,
		BuyIngredientDBSaver,
	)
}

// sbi - C2S_SHOP_BUY_ITEM
func BuyIngredient(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	// SCP command hanles when buying from currior

	ingredientID, _ := strconv.Atoi(req.Args[2])

	ingredientInfo, _ := utils.GetIngredient(ingredientID)

	ingredientAmount, _ := strconv.Atoi(req.Args[3])

	c.Player.AddCash(-ingredientInfo.Cash)

	c.Player.AddGold(-ingredientInfo.Gold)

	c.Location.Cafe().AddToFridge(ingredientID, ingredientAmount)

	c.SendExtensionResponse(cm.Identifier, "-1", "0", strconv.Itoa(ingredientID), strconv.Itoa(ingredientAmount))

	c.Player.UpdateAchivementBoughtIngredients()

	return nil
}

func BuyIngredientValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
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

	if c.Location.Cafe().GetOwnerID() != c.Player.GetID() {
		return "Not the owner!", commands.NOT_DECLARED
	}

	ingredientID, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	ingredientInfo, err := utils.GetIngredient(ingredientID)
	if err != nil {
		return "Invalid ingredient ID", commands.INVALID_VALUE
	}

	if ingredientInfo.Category == "fancy" && c.Player.GetLevel() < 7 {
		return "Cant buy that ingredient, because the player not yet reached the feature.", commands.NOT_DECLARED
	}

	if c.Player.GetLevel() < ingredientInfo.Level {
		return "Cant buy that ingredient, because the player not yet reached the level.", commands.NOT_DECLARED
	}

	ingredientAmount, err := strconv.Atoi(req.Args[3])
	if err != nil || ingredientAmount <= 0 {
		return "Invalid ingredient amount", commands.INVALID_VALUE
	}

	if !c.Player.GetIsTutorialCompleted() {
		if ingredientID != 1318 && ingredientAmount != 1 {
			return "Invalid tutorial values, needed: 1318 for ingredient, and 1 for amount!", commands.INVALID_VALUE
		}
	}

	if c.Location.Cafe().GetFridgeFreeSpace() < ingredientAmount {
		return fmt.Sprintf("Can't buy %v ingredient amount, because the fridge just have %v free space",
			ingredientAmount, c.Location.Cafe().GetFridgeFreeSpace()), commands.ERROR_FRIDGE_FULL
	}

	if ingredientInfo.Cash != 0 {
		if c.Player.GetCash() < ingredientInfo.Cash {
			return "Player not have enough cash to buy it", commands.NOT_ENOUGHT_MONEY
		}
	}

	if ingredientInfo.Gold != 0 {
		if c.Player.GetGold() < ingredientInfo.Gold {
			return "Player not have enough gold to buy it", commands.NOT_ENOUGHT_MONEY
		}
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func BuyIngredientDBSaver(c *client.Client) error {
	if c.Player.GetIsTutorialCompleted() {
		c.DB.UpdateAchievement(c.Player.GetID(), c.Player.GetAchivements().String())
		c.DB.UpdateCash(c.Player.GetID(), c.Player.GetCash())
		c.DB.UpdateGold(c.Player.GetID(), c.Player.GetGold())
		c.DB.UpdateFridgeInventory(c.Location.Cafe().GetID(), c.Location.Cafe().GetFridgeInventory().String())
	}

	return nil
}
