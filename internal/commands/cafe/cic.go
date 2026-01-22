package cafe

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/models/balancing"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"strconv"
	"time"
)

func init() {
	commands.RegisterCommand(requests.C2S_CAFE_INSTANTCOOK,
		commands.CommandConfig{
			Name:         "InstantCook",
			Identifier:   responses.S2C_CAFE_INSTANTCOOK,
			Description:  "Skip the dish cooking time, and instant cook it",
			Args:         "{objX} {objY}",
			MinArgs:      4,
			MaxArgs:      4,
			FeatureLevel: 5,
		},
		InstantCookValidator,
		InstantCook,
		InstantCookDBSaver,
	)
}

// cic - S2C_CAFE_INSTANTCOOK
// TODO: level ratio
func InstantCook(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) error {
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

	if c.Player.GetIsTutorialCompleted() {
		c.Player.AddGold(-dishCostGoldPerHours)
		c.Player.RemoveInstantCooking()
	}

	currentTime := time.Now().UTC()
	stove.SetStartedAt(&currentTime)
	stove.SetFinishesAt(&currentTime)

	c.Player.UpdateAchivementInstantCount()

	c.Location.Broadcast(cm.Identifier, "-1", "0", strconv.Itoa(objX), strconv.Itoa(objY))

	return nil
}

func InstantCookValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	// Dont allow players to modify the packet and sending us CIC while in editor.
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

	if c.Player.GetIsTutorialCompleted() {
		if c.Player.GetLevel() < cm.FeatureLevel {
			return "Cant use instant cook, because the player not yet reached the feature.", commands.NOT_DECLARED
		}
	}

	stove := c.Location.Cafe().GetObjectByPosXY(posX, posY)

	if stove == nil {
		return fmt.Sprintf("No stove found at: %d:%d", posX, posY), commands.NOT_DECLARED
	}

	dishID := stove.GetDishID()
	if dishID == -1 {
		return "Cant use instant cook, because the stove not have any valid dish ID", commands.NOT_DECLARED
	}

	dishInfo, err := utils.GetDish(dishID)
	if err != nil {
		return "Invalid dish ID", commands.INVALID_VALUE
	}

	if stove.GetIsRotten() {
		return "Cant use instant cook, because the dish is rotten!", commands.NOT_DECLARED
	}

	if stove.GetRemaingTime() <= 0 {
		return "Cant use instant cook, because the dish is already done!", commands.NOT_DECLARED
	}

	dishCostGoldPerHours := balancing.BalancingConstants.InstantCookHourPerGold
	dishDuration := dishInfo.Duration
	if dishDuration > 60 {
		dishCostGoldPerHours = dishDuration / 60
	}

	if c.Player.GetGold() < dishCostGoldPerHours || c.Player.GetInstantCookingsUsed() > c.Player.GetMaxInstants() {
		return "Player not have enough gold / instant cook lefts to use instant cook", commands.NOT_ENOUGHT_MONEY
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func InstantCookDBSaver(c *client.Client) error {
	if !c.Player.GetIsTutorialCompleted() {
		return nil
	}

	c.DB.UpdateAchievement(c.Player.GetID(), c.Player.GetAchivements().String())
	c.DB.UpdateGold(c.Player.GetID(), c.Player.GetGold())
	c.DB.UpdateObjects(c.Location.Cafe().GetID(), c.Location.Cafe().GetObjects().StringForDB())
	c.DB.UpdateInstantCookings(c.Player.GetID(), c.Player.GetInstantCookingsUsed())
	return nil
}
