package player

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/models/simple"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"
)

func init() {
	commands.RegisterCommand(requests.C2S_CAFE_WALK,
		commands.CommandConfig{
			Name:       "Walk",
			Identifier: responses.S2C_CAFE_WALK,
			MinArgs:    4,
			MaxArgs:    4,
		},
		CafeWalkValidator,
		CafeWalk,
		nil,
	)
}

// cwa - C2S_CAFE_WALK
func CafeWalk(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	posX, _ := strconv.Atoi(req.Args[2])
	posY, _ := strconv.Atoi(req.Args[3])

	c.Player.SetPos(simple.NewPosition(posX, posY))

	c.Location.Broadcast(cm.Identifier, "-1", "0", strconv.Itoa(c.Player.GetID()), "0", req.Args[2], req.Args[3])

	return nil
}

func CafeWalkValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	// Dont allow players to modify the packet and sending us CWA while in editor.
	if !c.Location.IsRunning() {
		return "The location is not running", commands.LOCATION_NOT_RUNNING
	}

	posX, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	posY, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	size := c.Location.Cafe().GetSize()
	if posX > size || posY > size || posX < 1 || posY < 1 {
		return "Cant walk there, because its an invalid pos!", commands.WALK_CANT_GO_THERE
	}

	if c.Location.Cafe().GetObjectByPosXY(posX, posY) != nil {
		return "Cant walk there, because theres an object!", commands.WALK_CANT_GO_THERE
	}

	return "Command ran without any errors.", commands.NO_ERROR
}
