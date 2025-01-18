package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"strconv"
)

func StoveDeliver(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Dont allow players to modify the packet and sending us CSD while in editor.
	if c.Location.Cafe().InEditorMode {
		return nil
	}

	stoveX, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	stoveY, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return err
	}

	counterX, err := strconv.Atoi(req.Args[4])
	if err != nil {
		return err
	}

	counterY, err := strconv.Atoi(req.Args[5])
	if err != nil {
		return err
	}

	stove := c.Location.Cafe().GetObjectByPos(stoveX, stoveY)
	counter := c.Location.Cafe().GetObjectByPos(counterX, counterY)

	// Add to counter
	wod, err := utils.GetDish(stove.DishID)
	if err != nil {
		return err
	}

	stove.DishAmount = wod.Servings

	if counter.DishID == stove.DishID {
		counter.DishAmount += stove.DishAmount
	} else {
		counter.DishID = stove.DishID
		counter.DishAmount = stove.DishAmount
	}

	// Get dish
	dish, err := utils.GetDish(stove.DishID)
	if err != nil {
		return err
	}

	// Reset stove
	stove.DishID = -2 // Dirty
	stove.FancyIng = false
	stove.StartedAt = nil
	stove.FinishesAt = nil

	// Increase xp
	c.Player.XP += dish.XP

	// Increase mastery
	// TODO: Create masteryies

	// response = ExtensionResponse('csd', '-1', '0', stove_x, stove_y, counter_x, counter_y, str(player.id))
	c.SendExtensionResponse(
		"csd", "-1", "0",
		req.Args[2],
		req.Args[3],
		req.Args[4],
		req.Args[5],
		strconv.Itoa(c.Player.ID),
	)

	return nil
}
