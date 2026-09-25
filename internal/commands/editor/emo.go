/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package editor

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/models/simple"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"strconv"
)

func init() {
	commands.RegisterCommand(requests.C2S_EDITOR_MOVE_OBJECT,
		commands.CommandConfig{
			Name:        "MoveObject",
			Identifier:  responses.S2C_EDITOR_MOVE_OBJECT,
			Description: "Moving an existing object in the Café",
			Args:        "{oldPosX} {oldPosY} {newPosX} {newPosY}",
			MinArgs:     6,
			MaxArgs:     6,
		},
		MoveObjectValidator,
		MoveObject,
		MoveObjectDBSaver,
	)
}

// emo - C2S_EDITOR_MOVE_OBJECT
func MoveObject(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {

	items, _ := utils.MultiAtoi(req.Args[2:]...)
	oldObjPos := simple.NewPosition(items[0], items[1])
	newObjPos := simple.NewPosition(items[2], items[3])
	obj := c.Location.Cafe().GetObjectByPos(oldObjPos)

	obj.SetPosXY(newObjPos.X, newObjPos.Y)

	c.SendExtensionResponse(cm.Identifier, "-1", "0", strconv.Itoa(oldObjPos.X), strconv.Itoa(oldObjPos.Y), strconv.Itoa(newObjPos.X), strconv.Itoa(newObjPos.X))
	return nil
}

func MoveObjectValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	// Dont allow players to modify the packet and sending us EMO while not in editor.
	if c.Location.IsRunning() {
		return "The location is running", commands.ERROR_EDITOR_ONLY_IN_EDITOR
	}

	if c.Location.Cafe().GetOwnerID() != c.Player.GetID() {
		return "Not the owner!", commands.NOT_DECLARED
	}

	items, err := utils.MultiAtoi(req.Args[2:]...)
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}
	oldObjPos := simple.NewPosition(items[0], items[1])
	newObjPos := simple.NewPosition(items[2], items[3])

	obj := c.Location.Cafe().GetObjectByPos(oldObjPos)
	if obj == nil {
		return "Target object not found!", commands.ERROR_EDITOR_POSITION_NOT_VALID
	}

	if obj.IsDoor() || obj.IsWall() || obj.IsWallObject() {
		size := c.Location.Cafe().GetSize()
		if newObjPos.X > size || newObjPos.Y > size || newObjPos.X < 0 || newObjPos.Y < 0 {
			return "Invalid position!", commands.ERROR_EDITOR_POSITION_NOT_VALID
		}
	} else {
		size := c.Location.Cafe().GetSize()
		if newObjPos.X > size || newObjPos.Y > size || newObjPos.X < 1 || newObjPos.Y < 1 {
			return "Invalid position!", commands.ERROR_EDITOR_POSITION_NOT_VALID
		}

		if newObjPos == c.Location.Cafe().GetPlayerStart() {
			return "Cant move the object to the playerstart!", commands.ERROR_EDITOR_POSITION_NOT_VALID
		}
	}
	// Theres a object in that space. The game removes the door render on door drag. We need to enable that to place back to the og. pos,
	if obj := c.Location.Cafe().GetObjectByPosXY(newObjPos.X, newObjPos.Y); obj != nil {
		// If not a door, give error.
		// I'm not sure if its really works, but if its works, we need to handle it.
		if !obj.IsDoor() {
			return "Theres an object in the pos!", commands.ERROR_EDITOR_POSITION_NOT_VALID
		}
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func MoveObjectDBSaver(c *client.Client) error {
	c.DB.UpdateObjects(c.Location.Cafe().GetID(), c.Location.Cafe().GetObjects().StringForDB())

	return nil
}
