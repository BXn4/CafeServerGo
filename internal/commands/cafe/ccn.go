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
)

func init() {
	commands.RegisterCommand(requests.C2S_CAFE_CLEAN,
		commands.CommandConfig{
			Name:        "Clean",
			Identifier:  responses.S2C_CAFE_CLEAN,
			Description: "Cleaning stoves / counters",
			Args:        "{0?} {objPosX} {objPosY} {isStove 0/1}",
			MinArgs:     4,
			MaxArgs:     4,
		},
		CleanValidator,
		Clean,
		CleanDBSaver,
	)
}

func Clean(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {

	objX, _ := strconv.Atoi(req.Args[2])
	objY, _ := strconv.Atoi(req.Args[3])
	obj := c.Location.Cafe().GetObjectByPosXY(objX, objY)

	if obj.IsStove() {
		c.Player.AddCash(-balancing.BalancingConstants.CleanCostCash)
		if obj.GetDishID() > 0 && obj.GetIsRotten() {
			c.Player.UpdateAchivementOvercookedFoods() // if the player cleans rotten food
		}
		obj.SetDishID(-1)
	}

	if obj.IsCounter() {
		obj.SetDishID(-1)
		obj.SetDishAmount(-1)
	}

	c.Location.Broadcast(
		cm.Identifier, "-1",
		"0",
		req.Args[2],
		req.Args[3],
		utils.If(obj.IsStove(), "1", "0"),
	)

	return nil
}

func CleanValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	// Dont allow players to modify the packet and sending us CCN while in editor.
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

	obj := c.Location.Cafe().GetObjectByPosXY(posX, posY)
	if obj == nil {
		return fmt.Sprintf("No object found at: %d:%d", posX, posY), commands.NOT_DECLARED
	}

	if obj.IsStove() {
		if c.Player.GetCash() < balancing.BalancingConstants.CleanCostCash {
			return "Player not have enough cash to clean the stove", commands.NOT_ENOUGHT_MONEY
		}
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func CleanDBSaver(c *client.Client) error {
	c.DB.UpdateObjects(c.Location.Cafe().GetID(), c.Location.Cafe().GetObjects().StringForDB())
	c.DB.UpdateCash(c.Player.GetID(), c.Player.GetCash())
	c.DB.UpdateAchievement(c.Player.GetID(), c.Player.GetAchivements().String())
	return nil
}
