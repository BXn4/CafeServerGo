package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"fmt"
	"strconv"
  "errors"
)

// sdi - C2S_SHOP_DELETE_ITEM
func SellIngredient(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {

  if c.Cafe.Owner() != c.Player.ID {
    return errors.New("You dont own this cafe!")
  }

  // Convert id got from cmd to int
	ingredientID, err := strconv.Atoi(req.Args[2])
	if err != nil {
		fmt.Printf("Can parse id to int: %v", err)
		return err
	}

  // Check if there is ingredient with that id
  ingredientInfo, err := utils.GetIngredient(ingredientID)
  if  err != nil {
		fmt.Printf("Invalid ingredient ID: %v", err)
  }

  // Convert amount got from cmd to int
	ingredientAmount, err := strconv.Atoi(req.Args[3])
	if err != nil {
		fmt.Printf("Can parse ingredient amount to int: %v", err)
		return err
	}

  // Check if amount is right
  if count, ok := c.Cafe.Fridge()[ingredientID]; !ok && count < ingredientAmount && ingredientAmount < 0{
    return fmt.Errorf("Invalid ingredient amount: %v, current amount: %v", ingredientAmount, count)
  }

  // Calcualte money
	if ingredientInfo.Cash != 0 {
		c.Player.Cash += int(float64(ingredientInfo.Cash) * 0.2 * float64(ingredientAmount))
	} else if ingredientInfo.Gold != 0 {
		c.Player.Cash += int(float64(ingredientInfo.Gold) * 0.2 * float64(ingredientAmount))
	}

	c.Cafe.Fridge()[ingredientID] -= ingredientAmount
	if c.Cafe.Fridge()[ingredientID] == 0 {
		// print("REMOVED!")
		delete(c.Cafe.Fridge(), ingredientID)
	}

	c.SendExtensionResponse("sdi", "-1", "0", strconv.Itoa(ingredientID), strconv.Itoa(ingredientAmount))

	return nil
}
