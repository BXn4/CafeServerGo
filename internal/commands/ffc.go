package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
)

// min level = 6

func init() {
	RegisterCommand(requests.C2S_FASTFOOD_COOK,
		CommandConfig{
			Name:       "FastFoodRefill",
			Identifier: responses.S2C_FASTFOOD_COOK,
			MinArgs:    5,
			MaxArgs:    5,
		},
		FastFoodRefillValidator,
		FastFoodRefill,
	)
}

// ffc - FastFoodRefill
func FastFoodRefill(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	items, _ := utils.MultiAtoi(req.Args[2:]...)
	x, y, id := items[0], items[1], items[2]
	drink, _ := utils.GetItem(id)

	c.Player.AddGold(-drink.Gold)
	c.Player.AddCash(-drink.Cash)

	obj := c.Location.Cafe().GetObjectByPosXY(x, y)

	obj.SetDishID(id)
	obj.SetDishAmount(drink.Servings)

	c.Location.Broadcast("ffc", "-1", "0",
		req.Args[2],
		req.Args[3],
		req.Args[4],
	)

	c.DB.UpdateCash(c.Player.ID, c.Player.GetCash())
	c.DB.UpdateGold(c.Player.ID, c.Player.GetGold())
	c.DB.UpdateObjects(c.Location.Cafe().ID, c.Location.Cafe().Objects.StringForDB())

	return nil
}

func FastFoodRefillValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	if !c.Location.IsRunning() {
		return "The location is not running", NOT_DECLARED
	}

	if c.Location.Cafe().GetPlayerID() != c.Player.ID {
		return "Not the owner!", NOT_DECLARED
	}

	items, err := utils.MultiAtoi(req.Args[2:]...)
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	x, y, id := items[0], items[1], items[2]
	drink, err := utils.GetItem(id)
	if err != nil {
		return "Cant get fastfood info!", NOT_DECLARED
	}

	obj := c.Location.Cafe().GetObjectByPosXY(x, y)
	if obj == nil {
		return "Invalid position!", NOT_DECLARED
	}

	if drink.Gold > c.Player.GetGold() || drink.Cash > c.Player.GetCash() {
		return "Player not have enough money to refill the fastfood machine!", NOT_ENOUGHT_MONEY
	}

	return "Command ran without any errors.", NO_ERROR
}
