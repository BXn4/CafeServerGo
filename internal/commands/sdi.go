package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"errors"
	"fmt"
	"math"
	"strconv"
)

// sdi - C2S_SHOP_DELETE_ITEM
func SellIngredient(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Dont allow players to modify the packet and sending us MJM while in editor.
	if !c.Location.IsRunning() {
		return nil
	}

	if c.Location.Cafe().GetPlayerID() != c.Player.ID {
		return errors.New("You dont own this cafe!")
	}

	// Convert id got from cmd to int
	ingredientID, err := strconv.Atoi(req.Args[2])
	if err != nil {

		return fmt.Errorf("Can parse id to int: %v", err)
	}

	// Check if there is ingredient with that id
	ingredientInfo, err := utils.GetIngredient(ingredientID)
	if err != nil {
		return fmt.Errorf("Invalid ingredient ID: %v", err)
	}

	// Convert amount got from cmd to int
	sellAmount, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return fmt.Errorf("Can parse ingredient amount to int: %v", err)
	}

	// Check if amount is right
	if count, ok := c.Location.Cafe().GetFridgeInventory()[ingredientID]; !ok && count < sellAmount && sellAmount < 0 {
		return fmt.Errorf("Invalid ingredient amount: %v, current amount: %v", sellAmount, count)
	}

	// Calcualte money
	c.Player.AddCash(sellAmount * int(math.Round(float64(ingredientInfo.Cash)*0.2+float64(ingredientInfo.Gold)*0.2)))

	c.Location.Cafe().RemoveFromFridge(ingredientID, sellAmount)

	c.SendExtensionResponse("sdi", "-1", "0", strconv.Itoa(ingredientID), strconv.Itoa(sellAmount))

	c.DB.UpdateCash(c.Player.ID, c.Player.GetCash())
	c.DB.UpdateFridgeInventory(c.Location.Cafe().ID, c.Location.Cafe().GetFridgeInventory().String())

	return nil
}
