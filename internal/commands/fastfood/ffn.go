package fastfood

import (
	"cafego/internal/agents"
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
	commands.RegisterCommand(requests.C2S_FASTFOOD_NPC,
		commands.CommandConfig{
			Name:         "FastFoodUse",
			Identifier:   responses.S2C_FASTFOOD_NPC,
			Description:  "The owner clicks on the drink icon",
			Args:         "{objX} {objY} {DrinkID} {CustomerID}",
			MinArgs:      5,
			MaxArgs:      5,
			FeatureLevel: 6,
		},
		FastFoodUseValidator,
		FastFoodUse,
		FastFoodUseDBSaver,
	)
}

// ffn
func FastFoodUse(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) error {

	items, err := utils.MultiAtoi(req.Args[2:]...)
	if err != nil {
		return fmt.Errorf("Cant convert args to integers")
	}

	_, x, y := items[0], items[1], items[2]
	obj := c.Location.Cafe().GetObjectByPosXY(x, y)
	// cs := c.Location.Cafe().GetCustomer(customerID)

	drink, err := utils.GetItem(obj.GetDishID())

	// Decrease dish amount
	obj.AddDishAmount(-1)

	// Add rewards
	c.Player.AddXP(drink.XP)
	c.Player.AddCash(drink.Cash)
	c.Location.Cafe().AddRating(drink.RatingBonus / 10)

	c.SendExtensionResponse(cm.Identifier, "-1", "0", req.Args[2], req.Args[3], req.Args[4])
	return nil
}

func FastFoodUseValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) (string, commands.ErrorCodes) {
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
		return "Cant use fastfood, because the player not yet reached the feature.", commands.NOT_DECLARED
	}

	items, err := utils.MultiAtoi(req.Args[2:]...)
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}
	customerID, x, y := items[0], items[1], items[2]

	obj := c.Location.Cafe().GetObjectByPosXY(x, y)
	if obj == nil {
		return "Invalid position!", commands.ERROR_FASTFOOD_OBJECT_NOT_FOUND
	}

	cs := c.Location.Cafe().GetCustomer(customerID)
	if cs == nil {
		return "Customer not found!", commands.NOT_DECLARED
	}

	start := agents.NewCafePoint(cs.GetPos(), c.Location.Cafe())
	end := agents.NewCafePoint(obj.GetPos(), c.Location.Cafe())
	_, _, reachable := agents.Path(start, end)

	if !reachable {
		return "Not reachable!", commands.ERROR_FASTFOOD_OUT_OF_REACH
	}

	_, err = utils.GetItem(obj.GetDishID())
	if err != nil {
		return "Cant get fastfood drink info!", commands.NOT_DECLARED
	}

	if obj.GetDishAmount() <= 0 {
		return "Fastfood empty!", commands.ERROR_FASTFOOD_EMPTY
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func FastFoodUseDBSaver(c *client.Client) error {
	c.DB.UpdateAchievement(c.Player.GetID(), c.Player.GetAchivements().String())

	c.DB.UpdateCash(c.Player.GetID(), c.Player.GetCash())
	c.DB.UpdateXP(c.Player.GetID(), c.Player.GetXP())
	c.DB.UpdateObjects(c.Location.Cafe().GetID(), c.Location.Cafe().GetObjects().StringForDB())
	c.DB.UpdateRating(c.Location.Cafe().GetID(), c.Location.Cafe().GetRating())

	return nil
}
