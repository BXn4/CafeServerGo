package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/balancing"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"strconv"
	"time"
)

// min level = 5

func init() {
	RegisterCommand(requests.C2S_CAFE_INSTANTCOOK,
		CommandConfig{
			Name:       "InstantCook",
			Identifier: responses.S2C_CAFE_INSTANTCOOK,
			MinArgs:    4,
			MaxArgs:    4,
		},
		InstantCookValidator,
		InstantCook,
	)
}

// cic - S2C_CAFE_INSTANTCOOK
// TODO: level ratio
func InstantCook(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	objX, _ := strconv.Atoi(req.Args[2])
	objY, _ := strconv.Atoi(req.Args[3])

	stove := c.Location.Cafe().GetObjectByPosXY(objX, objY)
	dishID := stove.GetDishID()
	dishInfo, _ := utils.GetDish(dishID)

	dishDuration := dishInfo.Duration
	dishCostGoldPerHours := balancing.BalancingConstants.InstantCookHourPerGold
	if dishDuration > 60 {
		dishCostGoldPerHours = dishDuration / 60
	}

	if c.Player.IsTutorialCompleted {
		c.Player.AddGold(-dishCostGoldPerHours)
		c.Player.RemoveInstantCooking()
	}

	currentTime := time.Now().UTC()
	stove.SetStartedAt(&currentTime)
	stove.SetFinishesAt(&currentTime)

	c.Player.UpdateAchivementInstantCount()

	if c.Player.IsTutorialCompleted {
		c.DB.UpdateAchievement(c.Player.ID, c.Player.GetAchivements().String())
		c.DB.UpdateGold(c.Player.ID, c.Player.GetGold())
		c.DB.UpdateObjects(c.Location.Cafe().ID, c.Location.Cafe().Objects.StringForDB())
		c.DB.UpdateInstantCookings(c.Player.ID, c.Player.GetInstantCookings())
	}

	c.Location.Broadcast("cic", "-1", "0", strconv.Itoa(objX), strconv.Itoa(objY))

	return nil
}

func InstantCookValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	// Dont allow players to modify the packet and sending us CIC while in editor.
	if !c.Location.IsRunning() {
		return "The location is not running", LOCATION_NOT_RUNNING
	}

	if c.Location.Cafe().GetPlayerID() != c.Player.ID {
		return "Not the owner!", NOT_DECLARED
	}

	posX, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	posY, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	if c.Player.IsTutorialCompleted {
		if c.Player.GetLevel() < 5 {
			return "Cant use instant cook, because the player not yet reached the feature.", NOT_DECLARED
		}
	}

	stove := c.Location.Cafe().GetObjectByPosXY(posX, posY)

	if stove == nil {
		return fmt.Sprintf("No stove found at: %d:%d", posX, posY), NOT_DECLARED
	}

	dishID := stove.GetDishID()
	if dishID == -1 {
		return "Cant use instant cook, because the stove not have any valid dish ID", NOT_DECLARED
	}

	dishInfo, err := utils.GetDish(dishID)
	if err != nil {
		return "Invalid dish ID", INVALID_VALUE
	}

	if stove.GetIsRotten() {
		return "Cant use instant cook, because the dish is rotten!", NOT_DECLARED
	}

	if stove.GetRemaingTime() <= 0 {
		return "Cant use instant cook, because the dish is already done!", NOT_DECLARED
	}

	dishCostGoldPerHours := balancing.BalancingConstants.InstantCookHourPerGold
	dishDuration := dishInfo.Duration
	if dishDuration > 60 {
		dishCostGoldPerHours = dishDuration / 60
	}

	if c.Player.GetGold() < dishCostGoldPerHours || c.Player.GetInstantCookings() > c.Player.GetMaxInstantCookings() {
		return "Player not have enough gold / instant cook lefts to use instant cook", NOT_ENOUGHT_MONEY
	}

	return "Command ran without any errors.", NO_ERROR
}
