package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"math"
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

	// Get dish
	dish, err := utils.GetDish(stove.GetDishID())
	if err != nil {
		return err
	}

	dishAmount := c.Player.GetDishMasteryServings(dish.ID)
	dishXP := c.Player.GetDishMasteryXP(dish.ID)

	if stove.GetFancyIng() {
		/*
			var _loc1_:int = Math.max(this.MIN_FANCY_SUBTRAHEND,Math.min(this.MAX_FANCY_SUBTRAHEND,int(this.getDuration() / 60 / 3)));
			return this.getMasterdServings() * (CafeConstants.fancyFactorServings - _loc1_) / 100;
		*/

		loc1 := int(dish.Duration / 60 / 3)
		loc1 = int(math.Max(0, math.Min(19, float64(loc1))))

		fancyAmount := dishAmount * (20 - loc1) / 100

		dishAmount += fancyAmount

		/*
			var _loc1_:int = Math.max(this.MIN_FANCY_SUBTRAHEND,Math.min(this.MAX_FANCY_SUBTRAHEND,int(this.getDuration() / 60 / 3)));
			return this.getMasterdXp() * (CafeConstants.fancyFactorXp - _loc1_) / 100;
		*/

		fancyXP := dishXP * (20 - loc1) / 100

		dishXP += fancyXP

		print(dishAmount)
	}

	if counter.GetDishID() == stove.GetDishID() {
		counter.AddDishAmount(dishAmount)
	} else {
		counter.SetDishID(stove.GetDishID())
		counter.AddDishAmount(dishAmount)
	}

	// Reset stove
	stove.SetDishID(-2) // Dirty
	stove.SetFancyIng(false)
	stove.SetStartedAt(nil)
	stove.SetFinishesAt(nil)

	// Increase xp
	c.Player.XP += c.Player.GetDishMasteryXP(dish.ID)

	c.Player.UpdateMastery(dish.ID)
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
