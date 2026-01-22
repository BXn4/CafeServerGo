package cafe

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/commands/coops"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

func init() {
	commands.RegisterCommand(requests.C2S_CAFE_STOVE_DELIVER,
		commands.CommandConfig{
			Name:        "StoveDeliver",
			Identifier:  responses.S2C_CAFE_STOVE_DELIVER,
			Description: "Deliver dishes from the stove to counters",
			Args:        "{stoveX} {stoveY} {counterX} {counterY} {playerID}",
			MinArgs:     6,
			MaxArgs:     6,
		},
		StoveDeliverValidator,
		StoveDeliver,
		StoveDeliverDBSaver,
	)
}

func StoveDeliver(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) error {
	stoveX, _ := strconv.Atoi(req.Args[2])
	stoveY, _ := strconv.Atoi(req.Args[3])
	counterX, _ := strconv.Atoi(req.Args[4])
	counterY, _ := strconv.Atoi(req.Args[5])

	stove := c.Location.Cafe().GetObjectByPosXY(stoveX, stoveY)
	counter := c.Location.Cafe().GetObjectByPosXY(counterX, counterY)
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

	masteryXP := c.Player.GetDishMasteryXP(dish.ID)

	println(masteryXP)

	// Increase xp
	c.Player.AddXP(masteryXP)

	c.Player.UpdateMastery(dish.ID)
	// TODO: Create masteryies

	c.Player.UpdateAchivementCookingCount()

	if c.Player.GetIsInCoop() {
		coop, _ := c.DB.GetCoop(c.Player.GetCoopID())

		coopInfo, _ := utils.GetCoop(coop.ActiveCoop)

		coopDishes := strings.Split(coopInfo.Dishes, "#")

		var dishes []int

		for _, dishRequirements := range coopDishes {
			dishRequirement := strings.Split(dishRequirements, "+")
			dishID, _ := strconv.Atoi(dishRequirement[0])
			dishes = append(dishes, dishID)
		}

		if slices.Contains(dishes, dish.ID) {
			coop.AddDish(dish.ID)

			if coop.IsDone() {
				coop.FinishLevel = coop.CalculateFinishLevel()
				coops.CoopFinish(&coop, gm)
			}
			c.DB.SaveCoop(&coop)
		}

	}

	// response = ExtensionResponse('csd', '-1', '0', stove_x, stove_y, counter_x, counter_y, str(player.id))
	c.Location.Broadcast(
		cm.Identifier, "-1", "0",
		req.Args[2],
		req.Args[3],
		req.Args[4],
		req.Args[5],
		strconv.Itoa(c.Player.GetID()),
	)

	return nil
}

func StoveDeliverValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	// Dont allow players to modify the packet and sending us CSD while in editor.
	if !c.Location.IsRunning() {
		return "The location is not running", commands.LOCATION_NOT_RUNNING
	}

	if c.Location.Cafe().GetOwnerID() != c.Player.GetID() {
		return "Not the owner!", commands.NOT_DECLARED
	}

	posX, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	posY, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	counterX, err := strconv.Atoi(req.Args[4])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	counterY, err := strconv.Atoi(req.Args[5])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	stove := c.Location.Cafe().GetObjectByPosXY(posX, posY)

	if stove == nil {
		return fmt.Sprintf("No stove found at: %d:%d", posX, posY), commands.NOT_DECLARED
	}

	counter := c.Location.Cafe().GetObjectByPosXY(counterX, counterY)
	if counter == nil {
		return fmt.Sprintf("No counter found at: %d:%d", posX, posY), commands.NOT_DECLARED
	}

	dishID := stove.GetDishID()
	if dishID == -1 {
		return "Cant use stove deliver info, because the stove not have any valid dish ID", commands.NOT_DECLARED
	}

	if stove.GetIsRotten() {
		return "Cant use stove deliver info, because the dish is rotten", commands.NOT_DECLARED
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func StoveDeliverDBSaver(c *client.Client) error {
	if !c.Player.GetIsTutorialCompleted() {
		return nil
	}

	c.DB.UpdateXP(c.Player.GetID(), c.Player.GetXP())
	c.DB.UpdateAchievement(c.Player.GetID(), c.Player.GetAchivements().String())
	c.DB.UpdateObjects(c.Location.Cafe().GetID(), c.Location.Cafe().GetObjects().StringForDB())

	return nil
}
