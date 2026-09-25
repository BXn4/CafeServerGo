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
	"cafego/internal/models/event"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"strconv"
)

// min level = 6 vending
// min level = 8 premium deco (gold?)

func init() {
	commands.RegisterCommand(requests.C2S_EDITOR_BUY_OBJECT,
		commands.CommandConfig{
			Name:        "BuyObject",
			Identifier:  responses.S2C_EDITOR_BUY_OBJECT,
			Description: "Buying objects",
			Args:        "{objX} {objY} {objID} {objRotation}",
			MinArgs:     6,
			MaxArgs:     6,
		},
		BuyObjectValidator,
		BuyObject,
		BuyObjectDBSaver,
	)
}

// ebu - C2S_EDITOR_BUY_OBJECT
// TODO: Need to check level.
func BuyObject(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	items, _ := utils.MultiAtoi(req.Args[2:]...)
	objX, objY, objID, objRotation := items[0], items[1], items[2], items[3]
	objectInfo, _ := utils.GetItem(objID)

	// If the player does not have the object in their inventory, dont remove cash, gold
	if c.Location.Cafe().GetFurnitureInventory()[objID] == 0 {
		if objectInfo.Cash > 0 {
			c.Player.AddCash(-objectInfo.Cash)
		}
		if objectInfo.Gold > 0 {
			c.Player.AddGold(-objectInfo.Gold)
		}
	} else {
		c.Location.Cafe().RemoveFurnitures(objID, 1)
	}

	// Need to add back the old wall in the inventory
	switch objectInfo.Group {
	case "Wall":
		oldWallID := c.Location.Cafe().GetTiles()[objX][objY]
		// If the old wall have luxury value, remove it from the Cafe
		c.Location.Cafe().AddLuxury(-(objectInfo.Cash / 4000) + (objectInfo.Gold * 2))
		c.Location.Cafe().AddFurnitures(oldWallID, 1)

		// Add the new wall
		c.Location.Cafe().SetTile(objX, objY, objID)

	case "Door":
		oldDoor := c.Location.Cafe().GetDoor()
		// If the old door have luxury value, remove it from the Cafe
		obj, _ := utils.GetDoor(int(oldDoor.GetKind()))
		c.Location.Cafe().AddLuxury(-(obj.Cash / 4000) + (obj.Gold * 2))
		// KIND = ID!!!
		c.Location.Cafe().GetFurnitureInventory()[int(oldDoor.GetKind())] = 1
		c.Location.Cafe().AddNewObject(objX, objY, objID, objRotation)
		c.Location.Cafe().RemoveObject(oldDoor.GetPos())

	case "Deco":
		c.Location.Cafe().AddNewObject(objX, objY, objID, objRotation)
		c.Player.UpdateAchivementBoughtDecoration()
		c.DB.UpdateAchievement(c.Player.GetID(), c.Player.GetAchivements().String())

	case "Fridge":
		c.Location.Cafe().AddNewObject(objX, objY, objID, objRotation)

	case "Stove":
		c.Location.Cafe().AddNewObject(objX, objY, objID, objRotation)

	case "Counter":
		c.Location.Cafe().AddNewObject(objX, objY, objID, objRotation)

	default:
		c.Location.Cafe().AddNewObject(objX, objY, objID, objRotation)
	}

	c.Location.Cafe().AddLuxury((objectInfo.Cash / 4000) + (objectInfo.Gold * 2))

	c.SendExtensionResponse(cm.Identifier, "-1", "0", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(objID), strconv.Itoa(objRotation))
	return nil
}

func BuyObjectValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	// Dont allow players to modify the packet and sending us EBU while not in editor.
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
	objX, objY, objID, objRotation := items[0], items[1], items[2], items[3]

	size := c.Location.Cafe().GetSize()
	if objX > size || objY > size || objX < 0 || objY < 0 {
		return "Invalid position!", commands.ERROR_EDITOR_POSITION_NOT_VALID
	}

	if objRotation < 0 || objRotation > 3 {
		return "Invalid rotation!", commands.ERROR_EDITOR_WATCH_OUT
	}

	objectInfo, err := utils.GetItem(objID)
	if err != nil {
		return "Object info not found!", commands.ERROR_EDITOR_WATCH_OUT
	}

	// Theres a object in that space. The game removes the door render on door drag. We need to enable that to place back to the og. pos,
	if obj := c.Location.Cafe().GetObjectByPosXY(objX, objY); obj != nil {
		// If not a door, give error.
		// I'm not sure if its really works, but if its works, we need to handle it.
		if !obj.IsDoor() {
			return "Theres an object in the pos!", commands.ERROR_EDITOR_POSITION_NOT_VALID
		}
	}

	// 2 0 0
	if (c.Location.Cafe().GetFurnitureInventory()[objID] == 0) && (event.GetEvent() < objectInfo.Events) {
		return "Cant buy that object, because theres no event!", commands.ERROR_EDITOR_WATCH_OUT
	}

	if objectInfo.Group == "Vendingmachine" && c.Player.GetLevel() < 6 {
		return "Cant buy that object, because the player not yet reached the feature.", commands.NOT_DECLARED
	}

	if objectInfo.Gold != 0 && c.Player.GetLevel() < 8 {
		return "Cant buy that object, because the player not yet reached the feature.", commands.NOT_DECLARED
	}

	// If the player does not have the object in their inventory, try to buy it.
	if c.Location.Cafe().GetFurnitureInventory()[objID] == 0 {
		if objectInfo.Cash > c.Player.GetCash() || objectInfo.Gold > c.Player.GetGold() {
			return "Player not have enough money to buy the object", commands.NOT_ENOUGHT_MONEY
		}
	}

	if objectInfo.Level > c.Player.GetLevel() {
		return "Player not yet unlocked the object", commands.ERROR_EDITOR_LEVEL_NOT_REACHED
	}

	if objectInfo.Category != "Door" {
		if objX == c.Location.Cafe().GetPlayerStart().X && objY == c.Location.Cafe().GetPlayerStart().Y {
			return "Cant move the object to the playerstart!", commands.ERROR_EDITOR_POSITION_NOT_VALID
		}
	}

	switch objectInfo.Group {
	case "Door":
		oldDoor := c.Location.Cafe().GetDoor()
		if oldDoor == nil {
			return "Cant get the old door!", commands.ERROR_EDITOR_WATCH_OUT
		}
		_, err := utils.GetDoor(int(oldDoor.GetKind()))
		if err != nil {
			return "Cant get the old door!", commands.ERROR_EDITOR_WATCH_OUT
		}

	case "Fridge":
		numberOfFridges := c.Location.Cafe().GetFridgeMaxCapacity() / 50
		if numberOfFridges > utils.GetLevelFridgesLimit(c.Player.GetLevel()) {
			return "Player have max level fridges in the cafe", commands.ERROR_EDITOR_WATCH_OUT
		}
	case "Stove":
		numberOfStoves := 0

		for _, obj := range c.Location.Cafe().GetObjects() {
			if obj.IsStove() {
				numberOfStoves++
			}
		}

		if numberOfStoves > utils.GetLevelStovesLimit(c.Player.GetLevel()) {
			return "Player have max level stoves in the cafe", commands.ERROR_EDITOR_WATCH_OUT
		}
	case "Counter":
		numberOfCounters := 0

		for _, obj := range c.Location.Cafe().GetObjects() {
			if obj.IsCounter() {
				numberOfCounters++
			}
		}

		if numberOfCounters > utils.GetLevelStovesLimit(c.Player.GetLevel()) {
			return "Player have max level counters in the cafe", commands.ERROR_EDITOR_WATCH_OUT
		}
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func BuyObjectDBSaver(c *client.Client) error {
	c.DB.UpdateCash(c.Player.GetID(), c.Player.GetCash())
	c.DB.UpdateGold(c.Player.GetID(), c.Player.GetGold())
	c.DB.UpdateObjects(c.Location.Cafe().GetID(), c.Location.Cafe().GetObjects().StringForDB())
	c.DB.UpdateLuxury(c.Location.Cafe().GetID(), c.Location.Cafe().GetLuxury())
	c.DB.UpdateFurnitureInventory(c.Location.Cafe().GetID(), c.Location.Cafe().GetFurnitureInventory().String())

	return nil
}
