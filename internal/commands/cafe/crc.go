package cafe

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/models/balancing"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"
	"time"
)

func init() {
	commands.RegisterCommand(requests.C2S_CAFE_RECOOK,
		commands.CommandConfig{
			Name:         "Recook",
			Identifier:   responses.S2C_CAFE_RECOOK,
			Description:  "Save the rotten dish",
			Args:         "{objX} {objY}",
			MinArgs:      4,
			MaxArgs:      4,
			FeatureLevel: 5,
		},
		RecookValidator,
		Recook,
		RecookDBSaver,
	)
}

// crc - Recook
func Recook(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {

	objX, _ := strconv.Atoi(req.Args[2])
	objY, _ := strconv.Atoi(req.Args[3])
	stove := c.Location.Cafe().GetObjectByPosXY(objX, objY)

	c.Player.AddGold(-balancing.BalancingConstants.RefreshFoodCost)
	// The recook works different!!
	// this.rottenCookTime = getTimer() + currentDish.baseDuration * 60000 / CafeConstants.timeFactor;
	// https://www.youtube.com/watch?v=GH_Fw6yAjJo
	cookingTime := c.Player.GetDishMasteryDuration(stove.GetDishID()) // returns in seconds
	currentTime := time.Now().UTC()

	if cookingTime < 60*60 {
		startedAt := currentTime.Add(-time.Duration(cookingTime)*time.Second - time.Hour) // from the current hour remove X seconds, then 1 hour
		finishedAt := currentTime.Add(time.Duration(cookingTime)*time.Second - time.Hour)

		stove.SetStartedAt(&startedAt)
		stove.SetFinishesAt(&finishedAt)
	} else {
		startedAt := currentTime
		finishedAt := currentTime

		stove.SetStartedAt(&startedAt)
		stove.SetFinishesAt(&finishedAt)
	}

	c.Player.UpdateAchivementOvercookedFoods()

	c.Location.Broadcast(cm.Identifier, "-1", "0", strconv.Itoa(objX), strconv.Itoa(objY))

	return nil
}

func RecookValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
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

	stove := c.Location.Cafe().GetObjectByPosXY(posX, posY)

	if stove == nil {
		return fmt.Sprintf("No stove found at: %d:%d", posX, posY), commands.NOT_DECLARED
	}

	dishID := stove.GetDishID()
	if dishID == -1 {
		return "Cant use instant cook, because the stove not have any valid dish ID", commands.NOT_DECLARED
	}

	if !stove.GetIsRotten() {
		return "Cant use recook, because the dish is not rotten!", commands.NOT_DECLARED
	}

	if c.Player.GetGold() < balancing.BalancingConstants.RefreshFoodCost {
		return "Player not have enough gold to use recook", commands.NOT_ENOUGHT_MONEY
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func RecookDBSaver(c *client.Client) error {
	c.DB.UpdateAchievement(c.Player.GetID(), c.Player.GetAchivements().String())

	c.DB.UpdateGold(c.Player.GetID(), c.Player.GetGold())
	c.DB.UpdateObjects(c.Location.Cafe().GetID(), c.Location.Cafe().GetObjects().StringForDB())

	return nil
}
