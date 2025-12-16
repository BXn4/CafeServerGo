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
)

func init() {
	RegisterCommand(requests.C2S_CAFE_CLEAN,
		CommandConfig{
			Name:       "Clean",
			Identifier: responses.S2C_CAFE_CLEAN,
			MinArgs:    4,
			MaxArgs:    4,
		},
		CleanValidator,
		Clean,
	)
}

func Clean(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	objX, _ := strconv.Atoi(req.Args[2])
	objY, _ := strconv.Atoi(req.Args[3])
	obj := c.Location.Cafe().GetObjectByPosXY(objX, objY)

	if obj.IsStove() {
		c.Player.AddCash(-balancing.BalancingConstants.CleanCostCash)
		if obj.GetDishID() > 0 && obj.GetIsRotten() {
			c.Player.UpdateAchivementOvercookedFoods() // if the player cleans rotten food

			c.DB.UpdateAchievement(c.Player.ID, c.Player.GetAchivements().String())
		}
		obj.SetDishID(-1)
	}

	if obj.IsCounter() {
		obj.SetDishID(-1)
		obj.SetDishAmount(-1)
	}

	c.Location.Broadcast(
		"ccn", "-1",
		"0",
		req.Args[2],
		req.Args[3],
		utils.If(obj.IsStove(), "1", "0"),
	)

	c.DB.UpdateObjects(c.Location.Cafe().ID, c.Location.Cafe().Objects.StringForDB())
	c.DB.UpdateCash(c.Player.ID, c.Player.GetCash())

	return nil
}

func CleanValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	// Dont allow players to modify the packet and sending us CCN while in editor.
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

	obj := c.Location.Cafe().GetObjectByPosXY(posX, posY)
	if obj == nil {
		return fmt.Sprintf("No object found at: %d:%d", posX, posY), NOT_DECLARED
	}

	if obj.IsStove() {
		if c.Player.GetCash() < balancing.BalancingConstants.CleanCostCash {
			return "Player not have enough cash to clean the stove", NOT_ENOUGHT_MONEY
		}
	}

	return "Command ran without any errors.", NO_ERROR
}
