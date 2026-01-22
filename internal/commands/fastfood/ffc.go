package fastfood

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
)

// min level = 6

func init() {
	commands.RegisterCommand(requests.C2S_FASTFOOD_COOK,
		commands.CommandConfig{
			Name:         "FastFoodRefill",
			Identifier:   responses.S2C_FASTFOOD_COOK,
			Description:  "Refill fastfood",
			Args:         "{posX} {posY} {DrinkID}",
			MinArgs:      5,
			MaxArgs:      5,
			FeatureLevel: 6,
		},
		FastFoodRefillValidator,
		FastFoodRefill,
		FastFoodDBSaver,
	)
}

// ffc - FastFoodRefill
func FastFoodRefill(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) error {
	items, _ := utils.MultiAtoi(req.Args[2:]...)
	x, y, id := items[0], items[1], items[2]
	drink, _ := utils.GetItem(id)

	c.Player.AddGold(-drink.Gold)
	c.Player.AddCash(-drink.Cash)

	obj := c.Location.Cafe().GetObjectByPosXY(x, y)

	obj.SetDishID(id)
	obj.SetDishAmount(drink.Servings)

	c.Location.Broadcast(cm.Identifier, "-1", "0",
		req.Args[2],
		req.Args[3],
		req.Args[4],
	)

	return nil
}

func FastFoodRefillValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	if !c.Location.IsRunning() {
		return "The location is not running", commands.NOT_DECLARED
	}

	if c.Location.Cafe().GetOwnerID() != c.Player.GetID() {
		return "Not the owner!", commands.NOT_DECLARED
	}

	if c.Player.GetLevel() < cm.FeatureLevel {
		return "Cant use fastfood refill, because the player not yet reached the feature.", commands.NOT_DECLARED
	}

	items, err := utils.MultiAtoi(req.Args[2:]...)
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	x, y, id := items[0], items[1], items[2]
	drink, err := utils.GetItem(id)
	if err != nil {
		return "Cant get fastfood info!", commands.NOT_DECLARED
	}

	obj := c.Location.Cafe().GetObjectByPosXY(x, y)
	if obj == nil {
		return "Invalid position!", commands.NOT_DECLARED
	}

	if drink.Gold > c.Player.GetGold() || drink.Cash > c.Player.GetCash() {
		return "Player not have enough money to refill the fastfood machine!", commands.NOT_ENOUGHT_MONEY
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func FastFoodDBSaver(c *client.Client) error {
	c.DB.UpdateCash(c.Player.GetID(), c.Player.GetCash())
	c.DB.UpdateGold(c.Player.GetID(), c.Player.GetGold())
	c.DB.UpdateObjects(c.Location.Cafe().GetID(), c.Location.Cafe().GetObjects().StringForDB())

	return nil
}
