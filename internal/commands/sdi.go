package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"fmt"
	"slices"
	"strconv"
)

// sdi - C2S_SHOP_DELETE_ITEM
func SellIngredient(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {

	ingredients := utils.Ingredients

	var ingredientIDs []int
	for _, ingredient := range ingredients {
		ingredientIDs = append(ingredientIDs, ingredient.ID)
	}

	ingredientID, err := strconv.Atoi(req.Args[2])
	if err != nil || !slices.Contains(ingredientIDs, ingredientID) {
		fmt.Printf("invalid ingredient ID: %v", err)
		return err
	}

	ingredientAmount, err := strconv.Atoi(req.Args[3])
	if err != nil || ingredientAmount <= 0 {
		fmt.Printf("invalid ingredient amount: %v", err)
		return err
	}

	ingredientInfo, err := utils.GetIngredientInfo(ingredientID)
	if err != nil {
		return nil
	}

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
