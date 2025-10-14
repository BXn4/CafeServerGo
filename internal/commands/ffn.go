package commands

import (
	"cafego/internal/agents"
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"fmt"
)

// ffn
func FastFoodCustomer(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	items, err := utils.MultiAtoi(req.Args[2:]...)
	if err != nil {
		return fmt.Errorf("Cant convert args to integers")
	}

	customerID, x, y := items[0], items[1], items[2]
	obj := c.Location.Cafe().GetObjectByPosXY(x, y)
	// If object not found
	if obj == nil {
		c.SendExtensionResponse("ffn", "-1", "999")
		return nil
	}

	cs := c.Location.Cafe().GetCustomer(customerID)
	if cs == nil {
		return fmt.Errorf("Cant find customer with id: %v", customerID)
	}

	// Check if object is reachable
	start := agents.NewCafePoint(cs.GetPos(), c.Location.Cafe())
	end := agents.NewCafePoint(obj.GetPos(), c.Location.Cafe())
	_, _, reachable := agents.Path(start, end)

	if !reachable {
		c.SendExtensionResponse("ffn", "-1", "36")
		return nil
	}

	drink, err := utils.GetItem(obj.GetDishID())
	if err != nil {
		return fmt.Errorf("Cant find dish with id: %v", obj.GetDishID())
	}

	// Decrease dish amount if cant send error
	if obj.AddDishAmount(-1) {
		c.SendExtensionResponse("ffn", "-1", "103")
		return nil
	}

	// Add rewards
	c.Player.AddXP(drink.XP)
	c.Player.AddCash(drink.Cash)
	c.Location.Cafe().AddRating(drink.RatingBonus)

	c.DB.UpdateAchievement(c.Player.ID, c.Player.GetAchivements().String())

	c.SendExtensionResponse("vck", "-1", "0", req.Args[2], req.Args[3], req.Args[4])
	return nil
}
