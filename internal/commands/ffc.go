package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
)

// vck - VersionCheck
func FastFoodCook(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	items, err := utils.MultiAtoi(req.Args[2:]...)
	if err != nil {
		return err
	}

	x, y, id := items[0], items[1], items[2]
	drink, err := utils.GetItem(id)
	if err != nil {
		return err
	}

	if drink.Gold > c.Player.GetGold() || drink.Cash > c.Player.GetCash() {
		c.SendExtensionResponse("ffc", "-1", "4")
		return nil
	}

	c.Player.AddGold(-drink.Gold)
	c.Player.AddCash(-drink.Cash)

	obj := c.Location.Cafe().GetObjectByPosXY(x, y)
	if obj == nil {
		return nil
	}

	obj.SetDishID(id)
	obj.SetDishAmount(drink.Servings)

	c.Location.Broadcast("ffc", "-1", "0",
		req.Args[2],
		req.Args[3],
		req.Args[4],
	)
	return nil
}
