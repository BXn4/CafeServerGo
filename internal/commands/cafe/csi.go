/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package cafe

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/models/object"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"
)

func init() {
	commands.RegisterCommand(requests.C2S_CAFE_STOVE_DELIVER_INFO,
		commands.CommandConfig{
			Name:        "StoveDeliverInfo",
			Identifier:  responses.S2C_CAFE_STOVE_DELIVER_INFO,
			Description: "Stove delivery sending counter location",
			Args:        "{counterX} {counterY}",
			MinArgs:     4,
			MaxArgs:     4,
		},
		StoveDeliverInfoValidator,
		StoveDeliverInfo,
		nil,
	)
}

func StoveDeliverInfo(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	stoveX, _ := strconv.Atoi(req.Args[2])
	stoveY, _ := strconv.Atoi(req.Args[3])

	stove := c.Location.Cafe().GetObjectByPosXY(stoveX, stoveY)

	// Choose counter that is empty or has the same food type
	var counter *object.Object
	for _, object := range c.Location.Cafe().GetObjects() {
		if !object.IsCounter() {
			continue
		}

		if stove.GetDishID() == object.GetDishID() {
			counter = object
			break
		} else if object.GetDishID() == -1 {
			counter = object
		}
	}

	// Set args
	counterX := counter.GetPos().X
	counterY := counter.GetPos().Y

	/*
		py	%xt%csi%-1%0%1%5%3%6%
		go	%xt%csi%-1%0%0%1%5%3%6%
	*/

	c.Location.Broadcast(
		cm.Identifier, "-1",
		"0",
		req.Args[2],
		req.Args[3],
		strconv.Itoa(counterX),
		strconv.Itoa(counterY),
	)

	return nil
}

func StoveDeliverInfoValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	// Dont allow players to modify the packet and sending us CSI while in editor.
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
		return "Cant use stove deliver info, because the stove not have any valid dish ID", commands.NOT_DECLARED
	}

	if stove.GetIsRotten() {
		return "Cant use stove deliver info, because the dish is rotten", commands.NOT_DECLARED
	}

	var counter *object.Object
	for _, object := range c.Location.Cafe().GetObjects() {
		if !object.IsCounter() {
			continue
		}

		if stove.GetDishID() == object.GetDishID() {
			counter = object
			break
		} else if object.GetDishID() == -1 {
			counter = object
		}
	}

	if counter == nil {
		return "No valid counters found", commands.NOT_DECLARED
	}

	return "Command ran without any errors.", commands.NO_ERROR
}
